package broker

import (
	"rabbitmq-service-broker/config"
	"rabbitmq-service-broker/rabbithutch"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

type RabbitMQServiceBroker struct {
	cfg         config.Config
	logger      lager.Logger
	rabbithutch rabbithutch.RabbitHutch
}

func New(cfg config.Config, rabbithutch rabbithutch.RabbitHutch, logger lager.Logger) brokerapi.ServiceBroker {
	return &RabbitMQServiceBroker{
		cfg:         cfg,
		logger:      logger,
		rabbithutch: rabbithutch,
	}
}

func (b *RabbitMQServiceBroker) ensureServiceInstanceExists(logger lager.Logger, serviceInstanceID string) error {
	ok, err := b.rabbithutch.VHostExists(serviceInstanceID)
	if err != nil {
		logger.Error("get-vhost-failed", err)
		return err
	}

	if !ok {
		logger.Info("vhost-not-found")
		return brokerapi.ErrInstanceDoesNotExist
	}

	return nil
}
