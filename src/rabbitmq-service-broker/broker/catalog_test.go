package broker_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"

	"rabbitmq-service-broker/broker"
	"rabbitmq-service-broker/broker/fakes"

	"code.cloudfoundry.org/lager/lagertest"
)

var _ = Describe("Service Broker", func() {
	It("returns a valid catalog", func() {
		cfg := broker.Config{
			ServiceConfig: broker.ServiceConfig{
				UUID:                "00000000-0000-0000-0000-000000000000",
				Name:                "p-rabbitmq",
				OfferingDescription: "this is a description",
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
		client := new(fakes.FakeAPIClient)
		logger := lagertest.NewTestLogger("test")
		broker := broker.New(cfg, client, logger)
		services, err := broker.Services(context.Background())
		Expect(err).NotTo(HaveOccurred())

		Expect(services).To(Equal([]brokerapi.Service{brokerapi.Service{
			ID:          cfg.ServiceConfig.UUID,
			Name:        cfg.ServiceConfig.Name,
			Description: cfg.ServiceConfig.OfferingDescription,
			Bindable:    true,
			Tags:        []string{"rabbitmq", "rabbit", "messaging", "message-queue", "amqp", "mqtt", "stomp"},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         "WhiteRabbitMQ",
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", "image_icon_base64"),
				LongDescription:     "this is a long description",
				ProviderDisplayName: "SomeCompany",
				DocumentationUrl:    "https://example.com",
				SupportUrl:          "https://support.example.com",
				Shareable:           &cfg.ServiceConfig.Shareable,
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
