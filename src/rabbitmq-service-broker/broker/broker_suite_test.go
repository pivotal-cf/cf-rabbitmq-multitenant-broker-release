package broker_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"rabbitmq-service-broker/broker"
	"rabbitmq-service-broker/config"

	"github.com/pivotal-cf/brokerapi"

	"rabbitmq-service-broker/rabbithutch/fakes"

	"code.cloudfoundry.org/lager/lagertest"
)

func TestBroker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Broker Suite")
}

func defaultConfig() config.Config {
	return config.Config{
		RabbitMQ: config.RabbitMQ{
			ManagementDomain: "foo.bar.com",
			Administrator: config.AdminCredentials{
				Username: "default-admin-username",
			},
			Hosts: []string{"fake-hostname-1", "fake-hostname-2"},
		},
	}
}

func defaultConfigWithUserTags() config.Config {
	cfg := defaultConfig()
	cfg.RabbitMQ.RegularUserTags = "administrator"
	return cfg
}

func defaultConfigWithExternalLoadBalancer() config.Config {
	cfg := defaultConfig()
	cfg.RabbitMQ.DNSHost = "my-dns-host.com"
	cfg.RabbitMQ.Hosts = []string{}
	return cfg
}

func defaultServiceBroker(conf config.Config, rabbithutch *fakes.FakeRabbitHutch) brokerapi.ServiceBroker {
	return broker.New(conf, rabbithutch, lagertest.NewTestLogger("test"))
}
