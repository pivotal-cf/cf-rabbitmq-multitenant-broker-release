package broker

import (
	"context"
	"errors"

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

func (b *RabbitMQServiceBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	return brokerapi.GetInstanceDetailsSpec{}, errors.New("Not implemented")
}

func (b *RabbitMQServiceBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, errors.New("Not implemented")
}

func (b *RabbitMQServiceBroker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, errors.New("Not implemented")
}

func (b *RabbitMQServiceBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	return brokerapi.GetBindingSpec{}, errors.New("Not implemented")
}

func (b *RabbitMQServiceBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, errors.New("Not implemented")
}
