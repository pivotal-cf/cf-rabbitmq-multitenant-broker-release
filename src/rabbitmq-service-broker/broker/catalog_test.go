package broker_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"

	"rabbitmq-service-broker/config"
	"rabbitmq-service-broker/rabbithutch/fakes"
)

var _ = Describe("Service Broker", func() {
	It("returns a valid catalog", func() {
		cfg := config.Config{
			Service: config.Service{
				UUID:                "00000000-0000-0000-0000-000000000000",
				Name:                "p-rabbitmq",
				Description:         "this is a description",
				DisplayName:         "WhiteRabbitMQ",
				IconImage:           "image_icon_base64",
				LongDescription:     "this is a long description",
				ProviderDisplayName: "SomeCompany",
				DocumentationURL:    "https://example.com",
				SupportURL:          "https://support.example.com",
				PlanUUID:            "11111111-1111-1111-1111-111111111111",
				Shareable:           false,
			},
		}
		rabbithutch := &fakes.FakeRabbitHutch{}
		broker := defaultServiceBroker(cfg, rabbithutch)
		services, err := broker.Services(context.Background())
		Expect(err).NotTo(HaveOccurred())

		Expect(services).To(Equal([]brokerapi.Service{brokerapi.Service{
			ID:          cfg.Service.UUID,
			Name:        cfg.Service.Name,
			Description: cfg.Service.Description,
			Bindable:    true,
			Tags:        []string{"rabbitmq", "rabbit", "messaging", "message-queue", "amqp", "mqtt", "stomp"},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         "WhiteRabbitMQ",
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", "image_icon_base64"),
				LongDescription:     "this is a long description",
				ProviderDisplayName: "SomeCompany",
				DocumentationUrl:    "https://example.com",
				SupportUrl:          "https://support.example.com",
				Shareable:           &cfg.Service.Shareable,
			},
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					ID:          "11111111-1111-1111-1111-111111111111",
					Name:        "standard",
					Description: "Provides a multi-tenant RabbitMQ cluster",
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: "Standard",
						Bullets:     []string{"RabbitMQ", "Multi-tenant"},
					},
				},
			},
		}}))
	})
})
