package broker_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"

	"rabbitmq-service-broker/broker"
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
				DocumentationUrl:    "https://example.com",
				SupportUrl:          "https://support.example.com",
			},
		}

		broker := broker.New(cfg)
		services, err := broker.Services(context.Background())
		Expect(err).NotTo(HaveOccurred())

		Expect(services).To(Equal([]brokerapi.Service{brokerapi.Service{
			ID:          cfg.ServiceConfig.UUID,
			Name:        cfg.ServiceConfig.Name,
			Description: cfg.ServiceConfig.OfferingDescription,
			Bindable:    true,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         "WhiteRabbitMQ",
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", "image_icon_base64"),
				LongDescription:     "this is a long description",
				ProviderDisplayName: "SomeCompany",
				DocumentationUrl:    "https://example.com",
				SupportUrl:          "https://support.example.com",
			},
		}}))
	})
})
