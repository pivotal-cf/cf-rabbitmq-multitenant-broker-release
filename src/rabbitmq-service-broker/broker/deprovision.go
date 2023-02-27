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
	logger.Info("entry", lager.Data{
		"service_instance_id": instanceID,
	})
	defer logger.Info("exit")

	if err := b.ensureServiceInstanceExists(logger, instanceID); err != nil {
		logger.Error("error-checking-service-instance-exists", err)
		return brokerapi.DeprovisionServiceSpec{}, err
	}

	if err := b.rabbithutch.VHostDelete(instanceID); err != nil {
		logger.Error("delete-vhost-failed", err)
		return brokerapi.DeprovisionServiceSpec{}, err
	}

	// Users upgrading from before the introduction of the Go broker may have some management users
	// still present on their system. If these exist, remove them
	users, err := b.rabbithutch.UserList()
	if err != nil {
		logger.Error("user-list-failed", err)
		return brokerapi.DeprovisionServiceSpec{}, err
	}

	prefix := fmt.Sprintf("mu-%v-", instanceID)
	for _, user := range users {
		if strings.HasPrefix(user, prefix) {
			logger.Info("delete-user", lager.Data{"username": user})
			if err := b.rabbithutch.DeleteUser(user); err != nil {
				logger.Error("delete-user-failed", err)
				return brokerapi.DeprovisionServiceSpec{}, err
			}
		}
	}

	return brokerapi.DeprovisionServiceSpec{}, nil
}
