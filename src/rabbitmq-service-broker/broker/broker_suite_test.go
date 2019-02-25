package broker_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"rabbitmq-service-broker/broker"

	"github.com/pivotal-cf/brokerapi"

	"rabbitmq-service-broker/broker/fakes"

	"code.cloudfoundry.org/lager/lagertest"
)

func TestBroker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Broker Suite")
}

func defaultConfig() broker.Config {
	return broker.Config{
		RabbitMQConfig: broker.RabbitMQConfig{
			ManagementDomain: "foo.bar.com",
			Administrator: broker.RabbitMQCredentials{
				Username: "default-admin-username",
			},
		},
	}
}

func defaultServiceBroker(config broker.Config, client *fakes.FakeAPIClient) brokerapi.ServiceBroker {
	return broker.New(config, client, lagertest.NewTestLogger("test"))
}
