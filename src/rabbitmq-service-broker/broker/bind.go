package broker

import (
	"context"

	"rabbitmq-service-broker/binding"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	logger := b.logger.Session("bind", lager.Data{
		"service_instance_id": instanceID,
		"binding_id":          bindingID,
	})
	logger.Info("entry")

	username := bindingID
	vhost := instanceID

	ok, err := b.rabbithutch.VHostExists(vhost)
	if err != nil {
		logger.Error("bind-error-checking-vhost-present", err)
		return brokerapi.Binding{}, err
	}
	if !ok {
		return brokerapi.Binding{}, brokerapi.ErrInstanceDoesNotExist
	}

	protocolPorts, err := b.rabbithutch.ProtocolPorts()
	if err != nil {
		logger.Error("bind-error-retrieving-protocol-ports", err)
		return brokerapi.Binding{}, err
	}

	password, err := b.rabbithutch.CreateUserAndGrantPermissions(username, vhost, b.cfg.RabbitMQ.RegularUserTags)
	if err != nil {
		logger.Error("bind-error-creating-user", err)
		return brokerapi.Binding{}, err
	}

	credsBuilder := binding.Builder{
		MgmtDomain:    b.cfg.RabbitMQ.ManagementDomain,
		Hostnames:     b.cfg.NodeHosts(),
		VHost:         vhost,
		Username:      username,
		Password:      password,
		TLS:           bool(b.cfg.RabbitMQ.TLS),
		ProtocolPorts: protocolPorts,
	}

	credentials, err := credsBuilder.Build()
	if err != nil {
		logger.Error("bind-error-building-credentials", err)
		return brokerapi.Binding{}, err
	}

	return brokerapi.Binding{Credentials: credentials}, nil
}
