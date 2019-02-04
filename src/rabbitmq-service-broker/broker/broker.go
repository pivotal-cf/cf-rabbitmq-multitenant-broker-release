package broker

import (
	"context"
	"errors"

	"github.com/pivotal-cf/brokerapi"
)

type RabbitmqServiceBroker struct{}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	err := brokerapi.NewFailureResponse(errors.New("Services Not Implemented"), 404, "")
	return []brokerapi.Service{}, err
}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("Provision Not Implemented"), 404, "")
	return brokerapi.ProvisionedServiceSpec{}, err
}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("Deprovision Not Implemented"), 404, "")
	return brokerapi.DeprovisionServiceSpec{}, err
}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("GetInstance Not Implemented"), 404, "")
	return brokerapi.GetInstanceDetailsSpec{}, err
}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("Update Not Implemented"), 404, "")
	return brokerapi.UpdateServiceSpec{}, err
}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	err := brokerapi.NewFailureResponse(errors.New("LastOperation Not Implemented"), 404, "")
	return brokerapi.LastOperation{}, err
}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	err := brokerapi.NewFailureResponse(errors.New("Bind Not Implemented"), 404, "")
	return brokerapi.Binding{}, err
}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("Unbind Not Implemented"), 404, "")
	return brokerapi.UnbindSpec{}, err
}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	err := brokerapi.NewFailureResponse(errors.New("GetBinding Not Implemented"), 404, "")
	return brokerapi.GetBindingSpec{}, err
}

func (rabbitmqServiceBroker *RabbitmqServiceBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	err := brokerapi.NewFailureResponse(errors.New("LastBindingOperation Not Implemented"), 404, "")
	return brokerapi.LastOperation{}, err
}
