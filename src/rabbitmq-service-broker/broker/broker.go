package broker

import (
	"rabbitmq-service-broker/config"
	"rabbitmq-service-broker/rabbithutch"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

type RabbitMQServiceBroker struct {
	cfg         config.Config
	client      rabbithutch.APIClient
	logger      lager.Logger
	rabbithutch rabbithutch.RabbitHutch
}

func New(cfg config.Config, client rabbithutch.APIClient, rabbithutch rabbithutch.RabbitHutch, logger lager.Logger) brokerapi.ServiceBroker {
	return &RabbitMQServiceBroker{
		cfg:         cfg,
		client:      client,
		logger:      logger,
		rabbithutch: rabbithutch,
	}
}
