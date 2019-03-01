package broker

import (
	"context"
	"fmt"
	"net/http"

	rabbithole "github.com/michaelklishin/rabbit-hole"

	"github.com/pivotal-cf/brokerapi"

	"code.cloudfoundry.org/lager"
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

func (b *RabbitMQServiceBroker) createVhost(vhost string) error {
	logger := b.logger.Session("create-vhost")
	logger.Info("get-vhost")
	if _, err := b.client.GetVhost(vhost); err == nil {
		err = brokerapi.ErrInstanceAlreadyExists
		logger.Error("get-vhost-failed", err)
		return err
	}
	logger.Info("get-vhost-succeeded")

	logger.Info("put-vhost")
	err := validateResponse(b.client.PutVhost(vhost, rabbithole.VhostSettings{}))
	if err != nil {
		logger.Error("put-vhost-failed", err)
		return err
	}
	logger.Info("put-vhost-succeeded")

	return nil
}

func (b *RabbitMQServiceBroker) deleteVhost(vhost string) error {
	logger := b.logger.Session("delete-vhost")
	logger.Info("delete-vhost")
	err := validateResponse(b.client.DeleteVhost(vhost))
	if err != nil {
		logger.Error("delete-vhost-failed", err)
		return err
	}
	logger.Info("delete-vhost-succeeded")
	return nil
}

func (b *RabbitMQServiceBroker) assignPermissionsToUser(vhost, username string) error {
	logger := b.logger.Session("assign-persmissions-to-user", lager.Data{"username": username})
	logger.Info("update-permissions")

	permissions := rabbithole.Permissions{Configure: ".*", Write: ".*", Read: ".*"}
	err := validateResponse(b.client.UpdatePermissionsIn(vhost, username, permissions))
	if err != nil {
		logger.Error("update-permissions-failed", err)
		return err
	}

	logger.Info("update-permissions-succeeded")
	return nil
}

func (b *RabbitMQServiceBroker) createPolicy(vhost string) error {
	logger := b.logger.Session("create-policy")

	policy := rabbithole.Policy{
		Definition: rabbithole.PolicyDefinition(b.cfg.RabbitMQ.OperatorSetPolicy.Definition),
		Priority:   b.cfg.RabbitMQ.OperatorSetPolicy.Priority,
		Vhost:      vhost,
		Pattern:    ".*",
		ApplyTo:    "all",
		Name:       b.cfg.RabbitMQ.OperatorSetPolicy.Name,
	}

	logger.Info("put-policy", lager.Data{"policy": policy})
	err := validateResponse(b.client.PutPolicy(vhost, b.cfg.RabbitMQ.OperatorSetPolicy.Name, policy))
	if err != nil {
		logger.Error("put-policy-failed", err)
		return err
	}

	logger.Info("put-policy-succeeded")
	return nil
}

func validateResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		return fmt.Errorf("http request failed with status code: %v", resp.StatusCode)
	}

	return nil
}
