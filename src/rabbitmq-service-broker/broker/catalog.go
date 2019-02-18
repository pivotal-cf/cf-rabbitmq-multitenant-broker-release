package broker

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

func (b RabbitMQServiceBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	return []brokerapi.Service{
		brokerapi.Service{
			ID:          b.Config.ServiceConfig.UUID,
			Name:        b.Config.ServiceConfig.Name,
			Description: b.Config.ServiceConfig.OfferingDescription,
			Bindable:    true,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         b.Config.ServiceConfig.DisplayName,
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", b.Config.ServiceConfig.IconImage),
				LongDescription:     b.Config.ServiceConfig.LongDescription,
				ProviderDisplayName: b.Config.ServiceConfig.ProviderDisplayName,
				DocumentationUrl:    b.Config.ServiceConfig.DocumentationUrl,
				SupportUrl:          b.Config.ServiceConfig.SupportUrl,
				Shareable:           &b.Config.ServiceConfig.Shareable,
			},
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					ID:          b.Config.ServiceConfig.PlanUuid,
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
