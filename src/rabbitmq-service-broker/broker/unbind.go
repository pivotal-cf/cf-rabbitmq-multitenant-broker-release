package broker

import (
	"context"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	logger := b.logger.Session("unbind")
	logger.Info("entry", lager.Data{
		"service_instance_id": instanceID,
		"binding_id":          bindingID,
	})
	defer logger.Info("exit")

	err := b.rabbithutch.DeleteUserAndConnections(bindingID)
	if err != nil {
		logger.Error("unbind-error", err)
		return brokerapi.UnbindSpec{}, err
	}

	return brokerapi.UnbindSpec{}, nil
}
