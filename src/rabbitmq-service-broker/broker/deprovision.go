package broker

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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

	// Users upgrading from before the introduction of the Go broker may have some dummy users still present on
	// their system. If these exist, remove them now.
	if username, err := b.FindFirstUserWithPrefix(fmt.Sprintf("mu-%v", instanceID)); err == nil && username != "" {
		logger.Info("Found mu-user", lager.Data{"username": username})
		b.deleteUser(username)
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

func (b *RabbitMQServiceBroker) FindFirstUserWithPrefix(prefix string) (string, error) {
	logger := b.logger.Session("find-first-user-with-prefix")

	logger.Info("find-first-user-with-prefix", lager.Data{"prefix": prefix})
	users, err := b.client.ListUsers()

	if err != nil {
		logger.Error("find-first-user-with-prefix-failed", err)
		return "", err
	}

	for _, user := range users {
		logger.Info("Attempting to match user against prefix", lager.Data{"username": user.Name, "prefix": prefix})
		if strings.HasPrefix(user.Name, prefix) {
			return user.Name, nil
		}
	}
	logger.Info("find-first-user-with-prefix")
	return "", nil
}
