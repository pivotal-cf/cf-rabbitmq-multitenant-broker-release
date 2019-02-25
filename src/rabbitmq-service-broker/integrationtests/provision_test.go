package integrationtests_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const provisionURL = baseURL + "service_instances/"

var _ = Describe("/v2/service_instance/:id", func() {
	var provisionDetails []byte

	BeforeEach(func() {
		var err error
		provisionDetails, err = json.Marshal(map[string]string{
			"service_id":        "00000000-0000-0000-0000-000000000000",
			"plan_id":           "11111111-1111-1111-1111-111111111111",
			"organization_guid": "fake-org-guid",
			"space_guid":        "fake-space-guid",
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("succeeds with HTTP 201 and provisions a service instance", func() {
		vhost := "7ef8c450-6aa1-4d8a-a56f-a8bb4ddd1de7"
		url := provisionURL + vhost
		response, body := doRequest(http.MethodPut, url, bytes.NewReader(provisionDetails))
		Expect(response.StatusCode).To(Equal(http.StatusCreated), string(body))

		By("returning the dashboard URL")
		var spec map[string]interface{}
		Expect(json.Unmarshal(body, &spec)).To(Succeed())
		Expect(spec).To(Equal(map[string]interface{}{
			"dashboard_url": "https://pivotal-rabbitmq.127.0.0.1/#/login/",
		}))

		By("creating a vhost")
		_, err := rmqClient.GetVhost(vhost)
		Expect(err).NotTo(HaveOccurred())

		By("granting access to the vhost to the rabbitmq administrator")
		//		rmqClient.Get
		perms, err := rmqClient.GetPermissionsIn(vhost, "guest")
		Expect(err).NotTo(HaveOccurred())
		Expect(perms.Configure).To(Equal(".*"))
		Expect(perms.Write).To(Equal(".*"))
		Expect(perms.Read).To(Equal(".*"))

		By("setting vhost policies")
		policies, err := rmqClient.ListPoliciesIn(vhost)
		Expect(err).NotTo(HaveOccurred())
		Expect(policies).To(HaveLen(1))
		Expect(policies[0].Name).To(Equal("operator_set_policy"))
		Expect(policies[0].Definition).To(Equal(rabbithole.PolicyDefinition{
			"ha-mode":      "exactly",
			"ha-params":    float64(2),
			"ha-sync-mode": "automatic",
		}))
		Expect(policies[0].Priority).To(Equal(50))
	})

	When("a service instance with the same ID has already been provisioned", func() {
		var vhost, url string

		BeforeEach(func() {
			vhost = "9ef8c450-6aa1-6d8a-b56f-a8bb4ddd1de4"
			url = provisionURL + vhost
			response, body := doRequest(http.MethodPut, url, bytes.NewReader(provisionDetails))
			Expect(response.StatusCode).To(Equal(http.StatusCreated), string(body))
		})

		It("fails with HTTP 409", func() {
			response, _ := doRequest(http.MethodPut, url, bytes.NewReader(provisionDetails))
			Expect(response.StatusCode).To(Equal(http.StatusConflict))
		})
	})
})
