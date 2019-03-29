package broker

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

func (b *RabbitMQServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	return []brokerapi.Service{
		brokerapi.Service{
			ID:          b.cfg.Service.UUID,
			Name:        b.cfg.Service.Name,
			Description: b.cfg.Service.Description,
			Bindable:    true,
			Tags:        []string{"rabbitmq", "rabbit", "messaging", "message-queue", "amqp", "mqtt", "stomp"},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         b.cfg.Service.DisplayName,
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", b.cfg.Service.IconImage),
				LongDescription:     b.cfg.Service.LongDescription,
				ProviderDisplayName: b.cfg.Service.ProviderDisplayName,
				DocumentationUrl:    b.cfg.Service.DocumentationURL,
				SupportUrl:          b.cfg.Service.SupportURL,
				Shareable:           &b.cfg.Service.Shareable,
			},
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					ID:          b.cfg.Service.PlanUUID,
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
