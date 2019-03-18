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

	if err := b.rabbithutch.EnsureVHostExists(vhost); err != nil {
		logger.Error("bind-error-checking-vhost-present", err)
		return brokerapi.Binding{}, err
	}

	password, err := b.rabbithutch.CreateUser(username, vhost, b.cfg.RabbitMQ.RegularUserTags)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	protocolPorts, err := b.protocolPorts()
	if err != nil {
		return brokerapi.Binding{}, err
	}

	credsBuilder := binding.Builder{
		MgmtDomain:    b.cfg.RabbitMQ.ManagementDomain,
		Hostnames:     b.cfg.RabbitMQ.Hosts,
		VHost:         vhost,
		Username:      username,
		Password:      password,
		TLS:           bool(b.cfg.RabbitMQ.TLS),
		ProtocolPorts: protocolPorts,
	}

	credentials, err := credsBuilder.Build()
	if err != nil {
		return brokerapi.Binding{}, err
	}

	return brokerapi.Binding{Credentials: credentials}, nil
}

func (b *RabbitMQServiceBroker) protocolPorts() (map[string]int, error) {
	protocolPorts, err := b.client.ProtocolPorts()
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for protocol, port := range protocolPorts {
		result[protocol] = int(port)
	}

	return result, nil
}
