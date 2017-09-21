require 'spec_helper'
require 'prof/marketplace_service'
require 'httparty'

require File.expand_path('../../../assets/rabbit-labrat/lib/lab_rat/aggregate_health_checker.rb', __FILE__)

RSpec.describe 'Logging CF service broker is well wired up' do
  SYSLOG_TEST_SERVER_URL = 'http://10.244.16.126:8080'

  let(:service_name) { environment.bosh_manifest.property('rabbitmq-broker.service.name') }
  let(:service) { Prof::MarketplaceService.new(name: service_name, plan: 'standard') }

  describe 'provisions a service' do
    it 'and writes the operation into the stdout logs', :creates_service_key do
      cf.provision_and_create_service_key(service) do |_, _, service_key_data|
        service_instance_id = service_key_data['vhost']
        expect(logs.find { |message| message.include? "Asked to provision a service: #{service_instance_id}"}).not_to be_nil
      end
    end
  end

  def logs
    logs = []

    loop do
      response = get(SYSLOG_TEST_SERVER_URL)

      break if has_no_logs? response

      log = JSON.parse(response.body)
      logs.push(log['message'])
    end

    logs
  end

  def get(url)
    HTTParty.get(url)
  end

  def has_no_logs?(response)
    !response.success? || response.body.start_with?('no logs available')
  end
end
