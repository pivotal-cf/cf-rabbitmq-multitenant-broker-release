package broker

import (
	"context"
	"fmt"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	logger := b.logger.Session("deprovision")

	if err := b.ensureServiceInstanceExists(logger, instanceID); err != nil {
		return brokerapi.DeprovisionServiceSpec{}, err
	}

	if err := b.deleteVhost(instanceID); err != nil {
		return brokerapi.DeprovisionServiceSpec{}, err
	}

	// Users upgrading from before the introduction of the Go broker may have some dummy users still present on
	// their system. If these exist, remove them now.
	if username, err := b.FindFirstUserWithPrefix(fmt.Sprintf("mu-%v", instanceID)); err == nil && username != "" {
		logger.Info("Found mu-user", lager.Data{"username": username})
		b.deleteUser(username)
	}

	logger.Info("deprovision-succeeded")
	return brokerapi.DeprovisionServiceSpec{}, nil
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
