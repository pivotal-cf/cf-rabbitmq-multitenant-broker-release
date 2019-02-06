package broker

import (
	"context"
	"errors"

	"github.com/pivotal-cf/brokerapi"
)

type RabbitMQServiceBroker struct{}

func (rabbitmqServiceBroker RabbitMQServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	err := brokerapi.NewFailureResponse(errors.New("Services Not Implemented"), 404, "")
	return []brokerapi.Service{}, err
}

func (rabbitmqServiceBroker RabbitMQServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("Provision Not Implemented"), 404, "")
	return brokerapi.ProvisionedServiceSpec{}, err
}

func (rabbitmqServiceBroker RabbitMQServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("Deprovision Not Implemented"), 404, "")
	return brokerapi.DeprovisionServiceSpec{}, err
}

func (rabbitmqServiceBroker RabbitMQServiceBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("GetInstance Not Implemented"), 404, "")
	return brokerapi.GetInstanceDetailsSpec{}, err
}

func (rabbitmqServiceBroker RabbitMQServiceBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("Update Not Implemented"), 404, "")
	return brokerapi.UpdateServiceSpec{}, err
}

func (rabbitmqServiceBroker RabbitMQServiceBroker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	err := brokerapi.NewFailureResponse(errors.New("LastOperation Not Implemented"), 404, "")
	return brokerapi.LastOperation{}, err
}

func (rabbitmqServiceBroker RabbitMQServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	err := brokerapi.NewFailureResponse(errors.New("Bind Not Implemented"), 404, "")
	return brokerapi.Binding{}, err
}

func (rabbitmqServiceBroker RabbitMQServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("Unbind Not Implemented"), 404, "")
	return brokerapi.UnbindSpec{}, err
}

func (rabbitmqServiceBroker RabbitMQServiceBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("GetBinding Not Implemented"), 404, "")
	return brokerapi.GetBindingSpec{}, err
}

func (rabbitmqServiceBroker RabbitMQServiceBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	err := brokerapi.NewFailureResponse(errors.New("LastBindingOperation Not Implemented"), 404, "")
	return brokerapi.LastOperation{}, err
}
