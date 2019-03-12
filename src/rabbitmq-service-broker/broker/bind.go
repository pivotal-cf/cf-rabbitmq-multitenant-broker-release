package broker

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	rabbithole "github.com/michaelklishin/rabbit-hole"

	"rabbitmq-service-broker/binding"

	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	logger := b.logger.Session("bind")

	username := bindingID
	vhost := instanceID

	exists, err := b.vhostExists(instanceID)
	if err != nil {
		logger.Error("bind-error-getting-vhost", err)
		return brokerapi.Binding{}, err
	}
	if !exists {
		logger.Error("bind-service-does-not-exist", brokerapi.ErrInstanceDoesNotExist)
		return brokerapi.Binding{}, brokerapi.ErrInstanceDoesNotExist
	}

	tags := b.cfg.RabbitMQ.RegularUserTags
	if tags == "" {
		tags = "policymaker,management"
	}

	password, err := generatePassword()
	if err != nil {
		return brokerapi.Binding{}, err
	}

	userSettings := rabbithole.UserSettings{
		Password: password,
		Tags:     tags,
	}
	err = b.createUser(username, userSettings)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	err = b.assignPermissionsToUser(vhost, username)
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

func generatePassword() (string, error) {
	rb := make([]byte, 24)
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rb), nil
}
