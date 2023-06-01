package broker

import (
	"context"
	"errors"

	"github.com/pivotal-cf/brokerapi/v10"
	"github.com/pivotal-cf/brokerapi/v10/domain"
)

func (b *RabbitMQServiceBroker) GetInstance(ctx context.Context, instanceID string, details domain.FetchInstanceDetails) (brokerapi.GetInstanceDetailsSpec, error) {
	return brokerapi.GetInstanceDetailsSpec{}, brokerapi.NewFailureResponse(errors.New("GetInstance Not Implemented"), 404, "")
}
