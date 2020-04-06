package broker

import (
	"context"

	"rabbitmq-service-broker/binding"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	logger := b.logger.Session("bind")
	logger.Info("entry", lager.Data{
		"service_instance_id": instanceID,
		"binding_id":          bindingID,
	})
	defer logger.Info("exit")

	username := bindingID
	vhost := instanceID

	if err := b.ensureServiceInstanceExists(logger, instanceID); err != nil {
		logger.Error("error-checking-service-instance-exists", err)
		return brokerapi.Binding{}, err
	}

	protocolPorts, err := b.rabbithutch.ProtocolPorts()
	if err != nil {
		logger.Error("error-retrieving-protocol-ports", err)
		return brokerapi.Binding{}, err
	}

	password, err := b.rabbithutch.CreateUserAndGrantPermissions(username, vhost, b.cfg.RabbitMQ.RegularUserTags)
	if err != nil {
		logger.Error("error-creating-user", err)
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
		logger.Error("error-building-credentials", err)
		return brokerapi.Binding{}, err
	}

	return brokerapi.Binding{Credentials: credentials}, nil
}
