package broker

import (
	"context"
	"net/http"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (spec brokerapi.DeprovisionServiceSpec, err error) {
	logger := b.logger.Session("deprovision")
	if _, err := b.client.GetVhost(instanceID); err != nil {
		if rabbitErr, ok := err.(rabbithole.ErrorResponse); ok && rabbitErr.StatusCode == http.StatusNotFound {
			logger.Info("vhost-not-found")
			return spec, brokerapi.ErrInstanceDoesNotExist
		}
		logger.Info("get-vhost-failed")
		return spec, err
	}
	if err := b.deleteVhost(instanceID); err != nil {
		return spec, err
	}
	logger.Info("deprovision-succeeded")
	return spec, nil
}
