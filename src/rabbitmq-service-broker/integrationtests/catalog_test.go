package integrationtests_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
)

const (
	url      = "http://localhost:8901/v2/catalog"
	username = "p1-rabbit"
	password = "p1-rabbit-testpwd"
)

var _ = Describe("/v2/catalog", func() {
	When("no credentials are provided", func() {
		It("fails with HTTP 401", func() {
			response, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})

	When("credentials are provided and they are correct", func() {
		It("succeeds with HTTP 200 and returns a valid catalog", func() {
			response, body := doRequest(http.MethodGet, url, nil)
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			catalog := make(map[string][]brokerapi.Service)
			Expect(json.Unmarshal(body, &catalog)).To(Succeed())

			Expect(catalog["services"]).To(HaveLen(1))
			// match against the expectation
			Expect(catalog["services"][0]).To(Equal(brokerapi.Service{
				ID:          "00000000-0000-0000-0000-000000000000",
				Name:        "p-rabbitmq",
				Description: "this is a description",
				Bindable:    true,
				Metadata: &brokerapi.ServiceMetadata{
					DisplayName:         "WhiteRabbitMQ",
					ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", "image_icon_base64"),
					LongDescription:     "this is a long description",
					ProviderDisplayName: "SomeCompany",
					DocumentationUrl:    "https://example.com",
					SupportUrl:          "https://support.example.com",
				},
			}))
		})
	})
})

func doRequest(method, url string, body io.Reader) (*http.Response, []byte) {
	req, err := http.NewRequest(method, url, body)
	Expect(err).ToNot(HaveOccurred())

	req.SetBasicAuth(username, password)
	req.Header.Set("X-Broker-API-Version", "2.14")

	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	Expect(err).ToNot(HaveOccurred())

	bodyContent, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())

	Expect(resp.Body.Close()).To(Succeed())
	return resp, bodyContent
}
