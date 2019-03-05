package broker

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	logger := b.logger.Session("provision")
	vhost := instanceID

	err = b.createVhost(vhost)
	if err != nil {
		return spec, err
	}

	defer func() {
		if err != nil {
			b.deleteVhost(vhost)
		}
	}()

	err = b.assignPermissionsToUser(vhost, b.cfg.RabbitMQ.Administrator.Username)
	if err != nil {
		return spec, err
	}

	if b.cfg.RabbitMQ.Management.Username != "" {
		err = b.assignPermissionsToUser(vhost, b.cfg.RabbitMQ.Management.Username)
		if err != nil {
			logger.Info("provision-management-user-skipped")
		}
	}

	if b.cfg.RabbitMQ.OperatorSetPolicy.Enabled {
		err = b.createPolicy(vhost)
		if err != nil {
			return spec, err
		}
	}

	logger.Info("provision-succeeded")
	url := fmt.Sprintf("https://%s/#/login/", b.cfg.RabbitMQ.ManagementDomain)
	return brokerapi.ProvisionedServiceSpec{DashboardURL: url}, nil
}
