require 'prof/external_spec/spec_helper'
require 'prof/environment/cloud_foundry'
require 'prof/environment_manager'
require 'yaml'
require 'pry'
require 'rspec/retry'

Dir[File.expand_path('support/**/*.rb', __dir__)].each do |file|
  require file
end

def environment
  @environment ||= begin
                     options = {
                       bosh_manifest_path: ENV.fetch('BOSH_MANIFEST') { File.expand_path('../../manifests/cf-rabbitmq-broker.yml', __FILE__) },
                       bosh_service_broker_job_name: 'cf-rabbitmq-multitenant-broker'
                     }

                     options[:cloud_foundry_domain]   = ENV.fetch('CF_DOMAIN', 'bosh-lite.com')
                     options[:cloud_foundry_username] = ENV.fetch('CF_USERNAME', 'admin')
                     options[:cloud_foundry_password] = ENV.fetch('CF_PASSWORD', 'admin')
                     options[:cloud_foundry_api_url]  = ENV.fetch('CF_API', 'api.bosh-lite.com')

                     options[:bosh_target]          = ENV.fetch('BOSH_TARGET', 'https://192.168.50.6:25555')
                     options[:bosh_username]        = ENV['BOSH_USERNAME']
                     options[:bosh_password]        = ENV['BOSH_PASSWORD']
                     options[:ssh_gateway_host]     = URI.parse(options[:bosh_target]).host

                     options[:ssh_gateway_username] = ENV.fetch('BOSH_SSH_USERNAME', 'vcap') if ENV.key?('BOSH_TARGET')

                     options.keep_if do |key, value|
                       not value.nil?
                     end

                     Prof::Environment::CloudFoundry.new(options)
                   end
end

BOSH_CLI = ENV.fetch("BOSH_CLI", 'bosh')

class Bosh2
  def initialize(bosh_cli = 'bosh', deployment = 'cf-rabbitmq-multitenant-broker')
    @bosh_cli = "#{bosh_cli} -n -d #{deployment}"

    version = `#{@bosh_cli} --version`
    raise 'BOSH CLI >= v2 required' unless version.start_with?('version 2.')
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
      output = ""
      exit_code = ::Open3.popen3("#{@bosh_cli} deploy #{manifest_file.path}") do |stdin, stdout, stderr, wait_thr|
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
    YAML.load(manifest)
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

def test_manifest
  YAML.load_file(ENV.fetch('BOSH_MANIFEST'))
end

def environment_manager
  cf_environment = OpenStruct.new(:cloud_foundry => cf)
  Prof::EnvironmentManager.new(cf_environment)
end

def cf
  @cf ||= environment.cloud_foundry
end

def test_app
  @test_app ||= Prof::TestApp.new(path: File.expand_path('../../assets/rabbit-labrat', __FILE__))
end

module ExcludeHelper
  def self.manifest
    @bosh_manifest ||= YAML.load(File.read(ENV['BOSH_MANIFEST']))
  end

  def self.metrics_available?
    return unless ENV['BOSH_MANIFEST']
    0 != manifest.fetch('releases').select{|i| i["name"] == "service-metrics" }.length
  end

  def self.warnings
    message = "\n"
    if !metrics_available?
      message += "WARNING: Skipping metrics tests, metrics are not available in this manifest\n"
    end

    message + "\n"
  end
end

puts ExcludeHelper::warnings

RSpec.configure do |config|
  config.include Matchers
  config.include TemplateHelpers, template: true

  Matchers::prints_logs_on_failure = true

  config.filter_run :focus
  config.run_all_when_everything_filtered = true
  config.filter_run_excluding :metrics => !ExcludeHelper::metrics_available?
  config.filter_run_excluding :test_with_errands => ENV.has_key?('SKIP_ERRANDS')
  config.filter_run_excluding :run_compliance_tests => (!ENV.has_key?('RUN_COMPLIANCE') && /darwin|mac os/ === RbConfig::CONFIG['host_os'] )

  config.around do |example|
    if example.metadata[:pushes_cf_app] || example.metadata[:creates_service_key]

      environment_manager.isolate_cloud_foundry do
        example.run
      end
    else
      example.run
    end
  end

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
