package broker

import (
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) vhostExists(vhost string) (bool, error) {
	logger := b.logger.Session("vhost-exists")
	logger.Info("get-vhost")
	if _, err := b.client.GetVhost(vhost); err != nil {
		if rabbitErr, ok := err.(rabbithole.ErrorResponse); ok && rabbitErr.StatusCode == http.StatusNotFound {
			logger.Info("vhost-not-found")
			return false, nil
		}
		logger.Info("get-vhost-failed")
		return false, err
	}
	return true, nil
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

func (b *RabbitMQServiceBroker) createUser(username string, info rabbithole.UserSettings) error {
	logger := b.logger.Session("create-user")

	logger.Info("put-user", lager.Data{"username": username})
	response, err := b.client.PutUser(username, info)
	if err != nil {
		logger.Error("put-user-failed", err)
		return err
	}
	if response != nil && response.StatusCode == http.StatusNoContent {
		logger.Error("put-user-failed", err)
		return brokerapi.ErrBindingAlreadyExists
	}
	logger.Info("put-user-succeeded")
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
