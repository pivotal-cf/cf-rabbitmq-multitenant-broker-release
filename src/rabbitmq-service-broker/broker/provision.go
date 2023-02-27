package broker

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	logger := b.logger.Session("provision")
	logger.Info("entry", lager.Data{
		"service_instance_id": instanceID,
	})
	defer logger.Info("exit")

	vhost := instanceID

	ok, existsErr := b.rabbithutch.VHostExists(vhost)
	if existsErr != nil {
		logger.Error("check-vhost-exists-failed", err)
		return spec, existsErr
	}
	if ok {
		logger.Error("vhost-already-exists", err)
		return spec, brokerapi.ErrInstanceAlreadyExists
	}

	err = b.rabbithutch.VHostCreate(vhost)
	if err != nil {
		logger.Error("vhost-create-failed", err)
		return spec, err
	}

	defer func() {
		if err != nil {
			if deleteErr := b.rabbithutch.VHostDelete(vhost); deleteErr != nil {
				logger.Error("delete-vhost-failed", deleteErr)
			}
		}
	}()

	err = b.rabbithutch.AssignPermissionsTo(vhost, b.cfg.RabbitMQ.Administrator.Username)
	if err != nil {
		logger.Error("provision-admin-user-failed", err)
		return spec, err
	}

	if b.cfg.RabbitMQ.Management.Username != "" {
		err = b.rabbithutch.AssignPermissionsTo(vhost, b.cfg.RabbitMQ.Management.Username)
		if err != nil {
			logger.Error("provision-management-user-skipped", err)
		}
	}

	if b.cfg.RabbitMQ.OperatorSetPolicy.Enabled {
		err = b.rabbithutch.CreatePolicy(vhost, b.cfg.RabbitMQ.OperatorSetPolicy.Name, b.cfg.RabbitMQ.OperatorSetPolicy.Priority, b.cfg.RabbitMQ.OperatorSetPolicy.Definition)
		if err != nil {
			logger.Error("put-policy-failed", err)
			return spec, err
		}
	}

	url := fmt.Sprintf("https://%s/#/login/", b.cfg.RabbitMQ.ManagementDomain)
	return brokerapi.ProvisionedServiceSpec{DashboardURL: url}, nil
}
