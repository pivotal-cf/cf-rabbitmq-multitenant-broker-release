package broker

import (
	"context"
	"errors"

	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	return brokerapi.GetInstanceDetailsSpec{}, brokerapi.NewFailureResponse(errors.New("GetInstance Not Implemented"), 404, "")
}
