require 'spec_helper'
require 'ostruct'
require 'papertrail'
require 'httparty'
require 'uri'

REMOTE_LOG_DESTINATION = Papertrail::Connection.new(token: ENV.fetch("PAPERTRAIL_TOKEN"))
PAPERTRAIL_GROUP_ID = ENV.fetch("PAPERTRAIL_GROUP_ID")

def get_instances(bosh_director_url, deployment_name)
  bosh_director_uri = URI(bosh_director_url)
  bosh_director_username, bosh_director_password = URI.unescape(bosh_director_uri.userinfo).split(":")

  JSON.parse(
    HTTParty.get(
      "#{bosh_director_uri.scheme}://#{bosh_director_uri.host}:#{bosh_director_uri.port}/deployments/#{deployment_name}/instances",
      basic_auth: {username: bosh_director_username, password: bosh_director_password},
      verify: false
    )
  ).map { |instance| OpenStruct.new(instance) }
end

DEPLOYMENT_NAME = ENV.fetch("DEPLOYMENT_NAME")
BOSH_DIRECTOR_URL = ENV.fetch("BOSH_DIRECTOR_URL")
DEPLOYMENT_INSTANCES = get_instances(BOSH_DIRECTOR_URL, DEPLOYMENT_NAME)

def host_search_string(host)
  "[job=#{host.job} index=#{host.index} id=#{host.id}]"
end

def five_minutes_ago
  Time.now - (5 * 60)
end

RSpec.describe "Syslog forwarding" do
  def has_event_for?(log_entry)
    options = { :group_id => PAPERTRAIL_GROUP_ID, :min_time => five_minutes_ago }
    query = REMOTE_LOG_DESTINATION.query(log_entry, options)
    !query.search.events.empty?
  end

  describe "rmq_broker host" do
    rmq_broker_hosts = DEPLOYMENT_INSTANCES.select { |i| i.job == "rmq-broker" }

    rmq_broker_hosts.each do |rmq_broker_host|
      job_host_log_entry = host_search_string(rmq_broker_host)
      it "forwards rmq_broker hosts logs (#{job_host_log_entry})" do
        expect(has_event_for?(job_host_log_entry)).to be_truthy
      end
    end
  end
end
