package broker

import (
	"context"
	"errors"

	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	return brokerapi.GetBindingSpec{}, brokerapi.NewFailureResponse(errors.New("GetBinding Not Implemented"), 404, "")
}
