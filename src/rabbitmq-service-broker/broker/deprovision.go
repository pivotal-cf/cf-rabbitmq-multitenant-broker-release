package broker

import (
	"context"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
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
	if err := b.deleteUser(fmt.Sprintf("mu-%v", instanceID)); err != nil {
		if rabbitErr, ok := err.(rabbithole.ErrorResponse); ok && rabbitErr.StatusCode == http.StatusNotFound {
			return spec, nil
		}

		return spec, err
	}
	logger.Info("deprovision-succeeded")
	return spec, nil
}

func (b *RabbitMQServiceBroker) deleteUser(username string) error {
	logger := b.logger.Session("delete-user")

	logger.Info("delete-user", lager.Data{"username": username})
	_, err := b.client.DeleteUser(username)
	if err != nil {
		logger.Error("delete-user-failed", err)
		return err
	}
	logger.Info("delete-user-succeeded")
	return nil
}
