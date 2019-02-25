package broker

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

func (b RabbitMQServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	return []brokerapi.Service{
		brokerapi.Service{
			ID:          b.config.ServiceConfig.UUID,
			Name:        b.config.ServiceConfig.Name,
			Description: b.config.ServiceConfig.OfferingDescription,
			Bindable:    true,
			Tags:        []string{"rabbitmq", "rabbit", "messaging", "message-queue", "amqp", "mqtt", "stomp"},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         b.config.ServiceConfig.DisplayName,
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", b.config.ServiceConfig.IconImage),
				LongDescription:     b.config.ServiceConfig.LongDescription,
				ProviderDisplayName: b.config.ServiceConfig.ProviderDisplayName,
				DocumentationUrl:    b.config.ServiceConfig.DocumentationURL,
				SupportUrl:          b.config.ServiceConfig.SupportURL,
				Shareable:           &b.config.ServiceConfig.Shareable,
			},
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					ID:          b.config.ServiceConfig.PlanUUID,
					Name:        "standard",
					Description: "Provides a multi-tenant RabbitMQ cluster",
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: "Standard",
						Bullets:     []string{"RabbitMQ", "Multi-tenant"},
					},
				},
			},
		},
	}, nil
}
