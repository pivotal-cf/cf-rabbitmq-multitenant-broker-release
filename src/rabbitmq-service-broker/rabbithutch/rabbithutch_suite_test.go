package rabbithutch_test

import (
	"rabbitmq-service-broker/broker"
	"rabbitmq-service-broker/config"
	"rabbitmq-service-broker/rabbithutch/fakes"
	"testing"

	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
)

func TestRabbithutch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rabbithutch Suite")
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

func defaultServiceBroker(conf config.Config, client *fakes.FakeAPIClient, rabbithutch *fakes.FakeRabbitHutch) brokerapi.ServiceBroker {
	return broker.New(conf, client, rabbithutch, lagertest.NewTestLogger("test"))
}
