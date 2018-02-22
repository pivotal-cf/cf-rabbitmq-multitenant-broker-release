require 'spec_helper'

require 'json'
require 'ostruct'
require 'tempfile'
require 'net/https'
require 'uri'

require 'rabbitmq/http/client'

require 'mqtt'
require 'stomp'
require 'net/https'
require 'httparty'

require File.expand_path('../../../assets/rabbit-labrat/lib/lab_rat/aggregate_health_checker.rb', __FILE__)

RSpec.describe 'Using a Cloud Foundry service broker' do
  def service_offering
    'standard'
  end

  def service_name
    'p-rabbitmq'
  end

  def app_name
    'testapp'
  end

  def service_instance_name
    'testservicename'
  end

  context 'when bound to a service instance' do
    before :all do
      @app = cf.push_app(test_app_path, app_name)

      cf.create_service_instance(service_name, service_offering, service_instance_name)
      cf.bind_app_to_service(app_name, service_instance_name)
      cf.restage_app(app_name)
    end

    after :all do
      cf.tear_down_org_and_space
    end

    it 'provides default connectivity', :pushes_cf_app do
      provides_amqp_connectivity(@app)
      provides_mqtt_connectivity(@app)
      provides_stomp_connectivity(@app)
    end

    it 'should provision a service key with defaults' do
      service_key_name = 'testservicekey'

      cf.service_key(service_instance_name, service_key_name) do |service_key_data|
        check_that_amqp_connection_data_is_present_in(service_key_data)
        check_that_stomp_connection_data_is_present_in(service_key_data)
      end
    end
  end

  context 'when unbinding a service instance' do
    after :each do
      cf.tear_down_org_and_space
    end

    it 'fails to connect if bindings are deleted', :pushes_cf_app do
      app = cf.push_app(test_app_path, app_name)
      cf.create_service_instance(service_name, service_offering, service_instance_name)
      cf.bind_app_to_service(app_name, service_instance_name)
      cf.restage_app(app_name)

      cf.unbind_app_from_service(app_name, service_instance_name)

      provides_no_amqp_connectivity(app)
      provides_no_mqtt_connectivity(app)
    end
  end

  describe 'redeploying cf-rabbitmq-multitenant-broker' do
    before :all do
      bosh.redeploy do |manifest|
        # Remove STOMP plugin
        rabbitmq_server_instance_group = manifest['instance_groups'].select { |instance_group| instance_group['name'] == 'rmq' }.first
        rabbitmq_server_job = rabbitmq_server_instance_group['jobs'].select { |job| job['name'] == 'rabbitmq-server' }.first
        rabbitmq_server_job['properties']['rabbitmq-server']['plugins'] = %w[rabbitmq_management rabbitmq_mqtt]

        # Broker with HA policy
        rabbitmq_broker_instance_group = manifest['instance_groups'].select { |instance_group| instance_group['name'] == 'rmq-broker' }.first
        rabbitmq_broker_job = rabbitmq_broker_instance_group['jobs'].select { |job| job['name'] == 'rabbitmq-broker' }.first
        rabbitmq_broker_job['properties']['rabbitmq-broker']['rabbitmq']['operator_set_policy'] = {
          'enabled' => true,
          'policy_name' => 'operator_set_policy',
          'policy_definition' => '{"ha-mode":"exactly","ha-params":2,"ha-sync-mode":"automatic"}',
          'policy_priority' => 50
        }

        # Configuring broker service metadata
        rabbitmq_broker_service = rabbitmq_broker_job['properties']['rabbitmq-broker']['service']
        rabbitmq_broker_service['name'] = 'service-name'
        rabbitmq_broker_service['display_name'] = 'apps-manager-test-name'
        rabbitmq_broker_service['offering_description'] = 'Some description of our service'
        rabbitmq_broker_service['long_description'] = 'Some long description of our service'
        rabbitmq_broker_service['icon_image'] = 'image-uri'
        rabbitmq_broker_service['provider_display_name'] = 'CompanyName'
        rabbitmq_broker_service['documentation_url'] = 'https://documentation.url'
        rabbitmq_broker_service['support_url'] = 'https://support.url'
      end

      # Setup Cloud Foundry
      @app = cf.push_app(test_app_path, app_name)
      cf.create_service_instance(service_name, service_offering, service_instance_name)
      cf.bind_app_to_service(app_name, service_instance_name)
      cf.restage_app(app_name)
    end

    after :all do
      bosh.deploy(test_manifest)
      cf.tear_down_org_and_space
    end

    context 'when stomp plugin is disabled' do
      it 'provides only amqp and mqtt connectivity', :pushes_cf_app do
        provides_amqp_connectivity(@app)
        provides_mqtt_connectivity(@app)
        provides_no_stomp_connectivity(@app)
      end
    end

    context 'when broker is configured with HA policy' do
      it 'sets queue policy to each created service instance', :pushes_cf_app do
        provides_mirrored_queue_policy_as_a_default(@app, service_name)
      end
    end

    describe 'the catalog' do
      let(:rmq_broker_username) do
        rabbitmq_broker_registrar_instance_group = test_manifest['instance_groups'].select { |instance_group| instance_group['name'] == 'rmq-broker' }.first
        rabbitmq_broker_registrar_job = rabbitmq_broker_registrar_instance_group['jobs'].select { |job| job['name'] == 'broker-registrar' }.first
        rabbitmq_broker_registrar_properties = rabbitmq_broker_registrar_job['properties']['broker']
        rabbitmq_broker_registrar_properties['username']
      end
      let(:rmq_broker_password) do
        rabbitmq_broker_registrar_instance_group = test_manifest['instance_groups'].select { |instance_group| instance_group['name'] == 'rmq-broker' }.first
        rabbitmq_broker_registrar_job = rabbitmq_broker_registrar_instance_group['jobs'].select { |job| job['name'] == 'broker-registrar' }.first
        rabbitmq_broker_registrar_properties = rabbitmq_broker_registrar_job['properties']['broker']
        rabbitmq_broker_registrar_properties['password']
      end
      let(:rmq_broker_host) do
        rabbitmq_broker_registrar_instance_group = test_manifest['instance_groups'].select { |instance_group| instance_group['name'] == 'rmq-broker' }.first
        rabbitmq_broker_registrar_job = rabbitmq_broker_registrar_instance_group['jobs'].select { |job| job['name'] == 'broker-registrar' }.first
        rabbitmq_broker_registrar_properties = rabbitmq_broker_registrar_job['properties']['broker']
        protocol = rabbitmq_broker_registrar_properties['protocol']
        host = rabbitmq_broker_registrar_properties['host']
        URI.parse("#{protocol}://#{host}")
      end
      let(:broker_catalog) do
        catalog_uri = URI.join(rmq_broker_host, '/v2/catalog')
        req = Net::HTTP::Get.new(catalog_uri)
        req.basic_auth(rmq_broker_username, rmq_broker_password)
        response = Net::HTTP.start(rmq_broker_host.hostname, rmq_broker_host.port, use_ssl: rmq_broker_host.scheme == 'https', verify_mode: OpenSSL::SSL::VERIFY_NONE) do |http|
          http.request(req)
        end
        JSON.parse(response.body)
      end
      let(:service_info) { broker_catalog['services'].first }
      let(:broker_catalog_metadata) { service_info['metadata'] }

      it 'has the correct name' do
        expect(service_info['name']).to eq('service-name')
      end

      it 'has the correct description' do
        expect(service_info['description']).to eq('Some description of our service')
      end

      it 'has the correct display name' do
        expect(broker_catalog_metadata['displayName']).to eq('apps-manager-test-name')
      end

      it 'has the correct long description' do
        expect(broker_catalog_metadata['longDescription']).to eq('Some long description of our service')
      end

      it 'has the correct image icon' do
        expect(broker_catalog_metadata['imageUrl']).to eq('data:image/png;base64,image-uri')
      end

      it 'has the correct provider display name' do
        expect(broker_catalog_metadata['providerDisplayName']).to eq('CompanyName')
      end

      it 'has the correct documentation url' do
        expect(broker_catalog_metadata['documentationUrl']).to eq('https://documentation.url')
      end

      it 'has the correct support url' do
        expect(broker_catalog_metadata['supportUrl']).to eq('https://support.url')
      end
    end
  end
end

def get(url)
  HTTParty.get(url, verify: false, timeout: 2)
end

def provides_amqp_connectivity(app)
  response = get("#{app.url}/services/rabbitmq/protocols/amqp091")
  expect(response.code).to eql(200)
  expect(response.body).to include('amq.gen')
end

def provides_mqtt_connectivity(app)
  response = get("#{app.url}/services/rabbitmq/protocols/mqtt")

  expect(response.code).to eql(200)
  expect(response.body).to include('mqtt://')
  expect(response.body).to include('Payload published')
end

def provides_stomp_connectivity(app)
  response = get("#{app.url}/services/rabbitmq/protocols/stomp")

  expect(response.code).to eql(200)
  expect(response.body).to include('Payload published')
end

def provides_no_amqp_connectivity(app)
  provides_no_connectivity_for(app, 'amqp091')
end

def provides_no_mqtt_connectivity(app)
  provides_no_connectivity_for(app, 'mqtt')
end

def provides_no_stomp_connectivity(app)
  provides_no_connectivity_for(app, 'stomp')
end

def provides_no_connectivity_for(app, protocol)
  # This is a work-around for #144893311

  response = get("#{app.url}/services/rabbitmq/protocols/#{protocol}")
  expect(response.code).to eql(500)
rescue Net::ReadTimeout => e
  puts "Caught exception #{e}!"
end

def check_that_amqp_connection_data_is_present_in(service_key_data)
  check_connection_data(service_key_data, 'amqp', 5672)
end

def check_that_stomp_connection_data_is_present_in(service_key_data)
  check_connection_data(service_key_data, 'stomp', 61_613)
end

def check_connection_data(service_key_data, protocol, port)
  expect(service_key_data).to have_key('protocols')
  expect(service_key_data['protocols']).to have_key(protocol)
  expect(service_key_data['protocols'][protocol]).to have_key('uri')
  expect(service_key_data['protocols'][protocol]['uri']).to start_with("#{protocol}://")
  expect(service_key_data['protocols'][protocol]).to have_key('host')
  expect(service_key_data['protocols'][protocol]['host']).not_to be_empty
  expect(service_key_data['protocols'][protocol]).to have_key('port')
  expect(service_key_data['protocols'][protocol]['port']).to eq(port)
  expect(service_key_data['protocols'][protocol]).to have_key('username')
  expect(service_key_data['protocols'][protocol]['username']).not_to be_empty
  expect(service_key_data['protocols'][protocol]).to have_key('password')
  expect(service_key_data['protocols'][protocol]['password']).not_to be_empty
  expect(service_key_data['protocols'][protocol]).to have_key('vhost')
  expect(service_key_data['protocols'][protocol]['vhost']).not_to be_empty
end

def provides_mirrored_queue_policy_as_a_default(app, service_name)
  vcap_services = cf.app_vcap_services(app.name)

  credentials = vcap_services[service_name].first['credentials']
  http_api_uris = credentials['http_api_uris']
  vhost = credentials['vhost']

  client = RabbitMQ::HTTP::Client.new(http_api_uris.first, ssl: { verify: false })
  policy = client.list_policies(vhost).find do |policy|
    policy['name'] == 'operator_set_policy'
  end

  expect(policy).to_not be_nil
  expect(policy['pattern']).to eq('.*')
  expect(policy['apply-to']).to eq('all')
  expect(policy['definition']).to eq('ha-mode' => 'exactly', 'ha-params' => 2, 'ha-sync-mode' => 'automatic')
  expect(policy['priority']).to eq(50)
end
