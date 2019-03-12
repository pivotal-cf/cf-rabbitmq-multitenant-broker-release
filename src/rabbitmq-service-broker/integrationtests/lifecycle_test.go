package integrationtests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the lifecycle of a service instance", func() {
	const (
		serviceID = "00000000-0000-0000-0000-000000000000"
		planID    = "11111111-1111-1111-1111-111111111111"
		bindingID = "22222222-2222-2222-2222-222222222222"
	)

	var provisionDetails, bindDetails []byte

	BeforeEach(func() {
		var err error
		provisionDetails, err = json.Marshal(map[string]string{
			"service_id":        serviceID,
			"plan_id":           planID,
			"organization_guid": "fake-org-guid",
			"space_guid":        "fake-space-guid",
		})
		Expect(err).NotTo(HaveOccurred())

		bindDetails, err = json.Marshal(map[string]string{
			"service_id": serviceID,
			"plan_id":    planID,
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("succeeds", func() {
		serviceInstanceID := "7ef8c450-6aa1-4d8a-a56f-a8bb4ddd1de7"

		By("sending a provision request")
		provisionResponse, provisionBody := doRequest(http.MethodPut, provisionURL(serviceInstanceID), bytes.NewReader(provisionDetails))
		Expect(provisionResponse.StatusCode).To(Equal(http.StatusCreated), string(provisionBody))

		By("checking that a dashboard URL is returned")
		var spec map[string]interface{}
		Expect(json.Unmarshal(provisionBody, &spec)).To(Succeed())
		Expect(spec).To(Equal(map[string]interface{}{
			"dashboard_url": "https://pivotal-rabbitmq.127.0.0.1/#/login/",
		}))

		By("checking that a vhost is created")
		_, err := rmqClient.GetVhost(serviceInstanceID)
		Expect(err).NotTo(HaveOccurred())

		By("checking that that rmq admin has access to the vhost")
		perms, err := rmqClient.GetPermissionsIn(serviceInstanceID, "guest")
		Expect(err).NotTo(HaveOccurred())
		Expect(perms.Configure).To(Equal(".*"))
		Expect(perms.Write).To(Equal(".*"))
		Expect(perms.Read).To(Equal(".*"))

		By("checking that vhost policies has been set")
		policies, err := rmqClient.ListPoliciesIn(serviceInstanceID)
		Expect(err).NotTo(HaveOccurred())
		Expect(policies).To(HaveLen(1))
		Expect(policies[0].Name).To(Equal("operator_set_policy"))
		Expect(policies[0].Definition).To(Equal(rabbithole.PolicyDefinition{
			"ha-mode":      "exactly",
			"ha-params":    float64(2),
			"ha-sync-mode": "automatic",
		}))
		Expect(policies[0].Priority).To(Equal(50))

		By("sending a binding request")
		bindResponse, bindBody := doRequest(http.MethodPut, bindURL(serviceInstanceID, bindingID), bytes.NewReader(bindDetails))
		Expect(bindResponse.StatusCode).To(Equal(http.StatusCreated), string(bindBody))

		By("checking the binding credentials")
		var binding map[string]interface{}
		Expect(json.Unmarshal(bindBody, &binding)).To(Succeed())
		Expect(binding).To(HaveKeyWithValue("credentials", SatisfyAll(
			HaveKeyWithValue("username", bindingID),
			HaveKeyWithValue("vhost", serviceInstanceID))))

		By("checking that that binding user has access to the vhost")
		perms, err = rmqClient.GetPermissionsIn(serviceInstanceID, bindingID)
		Expect(err).NotTo(HaveOccurred())
		Expect(perms.Configure).To(Equal(".*"))
		Expect(perms.Write).To(Equal(".*"))
		Expect(perms.Read).To(Equal(".*"))

		By("sending a deprovision request")
		deprovisionResponse, deprovisionBody := doRequest(http.MethodDelete, deprovisionURL(serviceInstanceID, serviceID, planID), nil)
		Expect(deprovisionResponse.StatusCode).To(Equal(http.StatusOK), string(deprovisionBody))

		By("checking that the vhost has been deleted")
		_, err = rmqClient.GetVhost(serviceInstanceID)
		rabbitErr := err.(rabbithole.ErrorResponse)
		Expect(rabbitErr.StatusCode).To(Equal(http.StatusNotFound))
	})

	Context("provisioning", func() {
		When("a service instance with the same ID has already been provisioned", func() {
			const serviceInstanceID = "9ef8c450-6aa1-6d8a-b56f-a8bb4ddd1de4"

			BeforeEach(func() {
				response, body := doRequest(http.MethodPut, provisionURL(serviceInstanceID), bytes.NewReader(provisionDetails))
				Expect(response.StatusCode).To(Equal(http.StatusCreated), string(body))
			})

			It("fails with HTTP 409", func() {
				response, _ := doRequest(http.MethodPut, provisionURL(serviceInstanceID), bytes.NewReader(provisionDetails))
				Expect(response.StatusCode).To(Equal(http.StatusConflict))
			})
		})
	})

	Context("deprovisioning", func() {
		When("a service instance has not been provisioned", func() {
			It("fails with HTTP 410", func() {
				serviceInstanceID := "does-not-exist"
				response, body := doRequest(http.MethodDelete, deprovisionURL(serviceInstanceID, serviceID, planID), nil)
				Expect(response.StatusCode).To(Equal(http.StatusGone), string(body))
			})
		})
	})
})

func provisionURL(serviceInstanceID string) string {
	return fmt.Sprintf("%sservice_instances/%s", baseURL, serviceInstanceID)
}

func deprovisionURL(serviceInstanceID, serviceID, planID string) string {
	return fmt.Sprintf("%s?service_id=%s&plan_id=%s", provisionURL(serviceInstanceID), serviceID, planID)
}

func bindURL(serviceInstanceID, bindingID string) string {
	return fmt.Sprintf("%sservice_instances/%s/service_bindings/%s", baseURL, serviceInstanceID, bindingID)
}
