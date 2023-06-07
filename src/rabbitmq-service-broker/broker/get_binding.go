package broker

import (
	"context"
	"errors"

	"github.com/pivotal-cf/brokerapi/v10"
	"github.com/pivotal-cf/brokerapi/v10/domain"
)

func (b *RabbitMQServiceBroker) GetBinding(ctx context.Context, instanceID, bindingID string, details domain.FetchBindingDetails) (brokerapi.GetBindingSpec, error) {
	return brokerapi.GetBindingSpec{}, brokerapi.NewFailureResponse(errors.New("GetBinding Not Implemented"), 404, "")
}
