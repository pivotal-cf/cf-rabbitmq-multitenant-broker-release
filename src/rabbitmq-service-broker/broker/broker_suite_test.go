package broker_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"rabbitmq-service-broker/broker"
	"rabbitmq-service-broker/config"

	"github.com/pivotal-cf/brokerapi"

	"rabbitmq-service-broker/broker/fakes"

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
		},
	}
}

func defaultServiceBroker(conf config.Config, client *fakes.FakeAPIClient) brokerapi.ServiceBroker {
	return broker.New(conf, client, lagertest.NewTestLogger("test"))
}
