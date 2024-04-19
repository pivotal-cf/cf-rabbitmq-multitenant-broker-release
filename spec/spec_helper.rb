require 'yaml'
require 'pry'
require 'rspec/retry'
require 'open3'
require 'shellwords'

Dir[File.expand_path('support/**/*.rb', __dir__)].each do |file|
  require file
end

BOSH_CLI = ENV.fetch('BOSH_CLI', 'bosh')

class Bosh2
  def initialize(bosh_cli = 'bosh')
    deployment = ENV['DEPLOYMENT_NAME'] || 'cf-rabbitmq-multitenant-broker'
    @bosh_cli = "#{bosh_cli} -n -d #{deployment}"

    version = `#{@bosh_cli} --version`
    begin
      major_version = version.match('version (\d+)\..*')[1].to_i
    rescue NoMethodError
      $stderr.print 'Could not parse bosh-cli version ' + version + "\n"
      raise
    end

    raise 'BOSH CLI >= v2 required' unless major_version > 1

  end

  def ssh(instance, command)
    command_escaped = Shellwords.escape(command)
    output = `#{@bosh_cli} ssh --gw-private-key #{key_path} #{instance} -r --json -c #{command_escaped}`
    JSON.parse(output)
  end

  def indexed_instance(instance, index)
    output = `#{@bosh_cli} instances | grep #{instance} | cut -f1`
    output.split(' ')[index]
  end

  def deploy(manifest)
    Tempfile.open('manifest.yml') do |manifest_file|
      manifest_file.write(manifest.to_yaml)
      manifest_file.flush
      output = ''
      exit_code = ::Open3.popen3("#{@bosh_cli} deploy #{manifest_file.path}") do |_stdin, stdout, _stderr, wait_thr|
        output << stdout.read
        wait_thr.value
      end
      abort "Deployment failed\n#{output}" unless exit_code == 0
    end
  end

  def key_path
    ssh_key_file = ENV.fetch('INSTANCE_SSH_KEY_FILE')
    raise 'please set environment variable INSTANCE_SSH_KEY_FILE' if ssh_key_file.nil? || ssh_key_file.empty?
    ssh_key_file
  end

  def redeploy
    deployed_manifest = manifest
    yield deployed_manifest
    deploy deployed_manifest
  end

  def manifest
    manifest = `#{@bosh_cli} manifest`
    YAML.safe_load(manifest)
  end

  def start(instance)
    `#{@bosh_cli} start #{instance}`
  end

  def stop(instance)
    `#{@bosh_cli} stop #{instance}`
  end
end

def bosh
  @bosh ||= Bosh2.new(BOSH_CLI)
end

class CloudFoundry
  attr_reader :domain, :api_url

  def initialize(domain, username, password, api_url)
    @domain = domain

    raise 'CF CLI is required' unless version.start_with?('cf version')

    target(api_url)
    login(username, password)
  end

  def target(api_url)
    cf("api #{api_url} --skip-ssl-validation")
  end

  def login(username, password)
    cf("auth #{username} #{password}")
  end

  def version
    cf('--version')
  end

  def create_service_instance(service, plan, service_instance_name)
    cf("create-service #{service} #{plan} #{service_instance_name}")
  end

  def bind_app_to_service(app_name, service_instance_name)
    cf("bind-service #{app_name} #{service_instance_name}")
  end

  def unbind_app_from_service(app_name, service_instance_name)
    cf("unbind-service #{app_name} #{service_instance_name}")
  end

  def service_key(service_instance_name, service_key_name)
    guid = service_instance_guid(service_instance_name)
    data = {
      'service_instance_guid' => guid,
      'name' => service_key_name
    }
    response = cf("curl -X POST /v2/service_keys -d '#{data.to_json}'")
    key = JSON.parse(response)
    yield key['entity']['credentials']
    cf("curl -X DELETE /v2/service_keys/#{key['metadata']['guid']}")
  end

  def service_instance_guid(service_instance_name)
    service_instance = cf("curl /v2/service_instances?q=name:#{service_instance_name}")
    JSON.parse(service_instance)['resources'].first['metadata']['guid']
  end

  def push_app(app_path, name)
    cf("push #{name} -p #{app_path}")

    App.new(name, "https://#{name}.#{domain}")
  end

  def start_app(name)
    cf("start #{name}")
  end

  def restage_app(name)
    cf("restage #{name}")
  end

  def url_for_app(app_name)
    "https://#{app_name}.#{domain}"
  end

  def app_vcap_services(app_name)
    all_apps = JSON.parse(cf('curl /v2/apps'))
    testapp = all_apps['resources'].find do |app|
      app['entity']['name'] == app_name
    end

    app_env = JSON.parse(cf("curl #{testapp['metadata']['url']}/env"))
    app_env['system_env_json']['VCAP_SERVICES']
  end

  def create_org_and_space(org_name, space_name)
    cf("create-org #{org_name}")
    cf("target -o #{org_name}")
    cf("create-space #{space_name}")
    cf("target -s #{space_name}")
    @org_name = org_name
    @space_name = space_name
  end

  def tear_down_org_and_space
    cf("delete-org #{org_name} -f")
  end

  def create_and_bind_security_group(security_group_name)
    cf("create-security-group #{security_group_name} #{security_group_path}")
    cf("bind-security-group #{security_group_name} #{org_name} --space #{space_name}")
  end

  private

  attr_reader :org_name, :space_name

  def cf(command)
    p "cf #{command}" unless command.include?('auth')
    stdout, stderr, status = Open3.capture3("cf #{command}")
    raise "error executing command\n#{stderr}" if status.exitstatus == 1

    stdout
  end
end

class App < Struct.new(:name, :url)
end

def cf
  return @cf unless @cf.nil?

  domain = ENV.fetch('CF_DOMAIN', 'bosh-lite.com')
  username = ENV.fetch('CF_USERNAME', 'admin')
  password = ENV.fetch('CF_PASSWORD', 'admin')
  api_url = ENV.fetch('CF_API', 'api.bosh-lite.com')

  @cf = CloudFoundry.new(domain, username, password, api_url)
  @cf.create_org_and_space("cf-org-#{random_string}", "cf-space-#{random_string}")
  @cf.create_and_bind_security_group("cf-sg-#{random_string}")
  @cf
end

def random_string
  [*('A'..'Z')].sample(8).join
end

def test_manifest
  YAML.load_file(ENV.fetch('BOSH_MANIFEST'))
end

def test_app_path
  File.expand_path('../../assets/rabbit-labrat', __FILE__)
end

def security_group_path
  File.expand_path('../assets/security-group.json', __FILE__)
end

module ExcludeHelper
  def self.manifest
    @bosh_manifest ||= YAML.safe_load(File.read(ENV['BOSH_MANIFEST']))
  end

  def self.metrics_available?
    return unless ENV['BOSH_MANIFEST']
    !manifest.fetch('releases').select { |i| i['name'] == 'service-metrics' }.empty?
  end

  def self.warnings
    message = "\n"
    unless metrics_available?
      message += "WARNING: Skipping metrics tests, metrics are not available in this manifest\n"
    end

    message + "\n"
  end
end

puts ExcludeHelper.warnings

RSpec.configure do |config|
  config.include Matchers
  config.include TemplateHelpers, template: true

  Matchers.prints_logs_on_failure = true

  config.filter_run :focus
  config.run_all_when_everything_filtered = true
  config.filter_run_excluding metrics: !ExcludeHelper.metrics_available?
  config.filter_run_excluding test_with_errands: ENV.key?('SKIP_ERRANDS')
  config.filter_run_excluding run_compliance_tests: (!ENV.key?('RUN_COMPLIANCE') && RbConfig::CONFIG['host_os'] === /darwin|mac os/)

  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end

  config.disable_monkey_patching!

  # show retry status in spec process
  config.verbose_retry = true
  # show exception that triggers a retry if verbose_retry is set to true
  config.display_try_failure_messages = true

  config.around :each, :retryable do |ex|
    ex.run_with_retry retry: 60, retry_wait: 10
  end

  Kernel.srand config.seed
end
