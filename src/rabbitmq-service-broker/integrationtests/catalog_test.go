package integrationtests_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
)

const catalogURL = baseURL + "catalog"

var _ = Describe("/v2/catalog", func() {
	It("succeeds with HTTP 200 and returns a valid catalog", func() {
		response, body := doRequest(http.MethodGet, catalogURL, nil)
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		catalog := make(map[string][]brokerapi.Service)
		Expect(json.Unmarshal(body, &catalog)).To(Succeed())

		Expect(catalog["services"]).To(HaveLen(1))

		shareable := false

		Expect(catalog["services"][0]).To(Equal(brokerapi.Service{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "p-rabbitmq",
			Description: "this is a description",
			Bindable:    true,
			Tags:        []string{"rabbitmq", "rabbit", "messaging", "message-queue", "amqp", "mqtt", "stomp"},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         "WhiteRabbitMQ",
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", "image_icon_base64"),
				LongDescription:     "this is a long description",
				ProviderDisplayName: "SomeCompany",
				DocumentationUrl:    "https://example.com",
				SupportUrl:          "https://support.example.com",
				Shareable:           &shareable,
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
		}))
	})
})
