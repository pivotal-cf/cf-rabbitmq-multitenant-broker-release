package integrationtests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	client "github.com/streadway/amqp"
)

var _ = Describe("the lifecycle of a service instance", func() {
	const (
		serviceID = "00000000-0000-0000-0000-000000000000"
		planID    = "11111111-1111-1111-1111-111111111111"
		bindingID = "22222222-2222-2222-2222-222222222222"
	)

	It("succeeds for one SI and binding", func() {
		serviceInstanceID := "7ef8c450-6aa1-4d8a-a56f-a8bb4ddd1de7"

		By("sending a provision request")
		provisionResponse, provisionBody := provision(serviceInstanceID, serviceID, planID)
		Expect(provisionResponse.StatusCode).To(Equal(http.StatusCreated), string(provisionBody))

		By("creating a management user")
		managementUsername := fmt.Sprintf("mu-%v-8912348912389123", serviceInstanceID)
		res, err := rmqClient.PutUser(managementUsername, rabbithole.UserSettings{})
		if err == nil {
			res.Body.Close()
		}
		user, _ := rmqClient.GetUser(managementUsername)
		Expect(user.Name).To(Equal(managementUsername))

		By("checking that a dashboard URL is returned")
		var spec map[string]interface{}
		Expect(json.Unmarshal(provisionBody, &spec)).To(Succeed())
		Expect(spec).To(Equal(map[string]interface{}{
			"dashboard_url": "https://pivotal-rabbitmq.127.0.0.1/#/login/",
		}))

		By("checking that a vhost is created")
		_, err = rmqClient.GetVhost(serviceInstanceID)
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
		bindResponse, bindBody := bind(bindingID, serviceInstanceID, serviceID, planID)
		Expect(bindResponse.StatusCode).To(Equal(http.StatusCreated), string(bindBody))

		By("checking the binding credentials")
		var binding map[string]interface{}
		Expect(json.Unmarshal(bindBody, &binding)).To(Succeed())
		Expect(binding).To(HaveKeyWithValue("credentials", SatisfyAll(
			HaveKeyWithValue("username", bindingID),
			HaveKeyWithValue("vhost", serviceInstanceID))))

		By("checking that the user exists")
		userInfo, userErr := rmqClient.GetUser(bindingID)
		Expect(userErr).NotTo(HaveOccurred())
		Expect(userInfo.Name).To(Equal(bindingID))

		By("checking that the binding user has access to the vhost")
		perms, err = rmqClient.GetPermissionsIn(serviceInstanceID, bindingID)
		Expect(err).NotTo(HaveOccurred())
		Expect(perms.Configure).To(Equal(".*"))
		Expect(perms.Write).To(Equal(".*"))
		Expect(perms.Read).To(Equal(".*"))

		By("creating a connection")
		creds := binding["credentials"].(map[string]interface{})
		protocols := creds["protocols"].(map[string]interface{})
		amqp := protocols["amqp"].(map[string]interface{})
		uris := amqp["uris"].([]interface{})
		conn, err := client.Dial(fmt.Sprintf(uris[0].(string)))
		Expect(err).NotTo(HaveOccurred())
		_, err = conn.Channel()
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() int {
			conns, _ := rmqClient.ListConnections()
			return len(conns)
		}, 5).Should(Equal(1))

		By("sending an unbind request")
		unbindResponse, unbindBody := doRequest(http.MethodDelete, unbindURL(serviceInstanceID, bindingID, serviceID, planID), nil)
		Expect(unbindResponse.StatusCode).To(Equal(http.StatusOK), string(unbindBody))

		By("checking that the user no longer exists")
		_, userErr = rmqClient.GetUser(bindingID)
		Expect(userErr).To(MatchError(ContainSubstring("Error 404")))

		By("sending a deprovision request")
		deprovisionResponse, deprovisionBody := doRequest(http.MethodDelete, deprovisionURL(serviceInstanceID, serviceID, planID), nil)
		Expect(deprovisionResponse.StatusCode).To(Equal(http.StatusOK), string(deprovisionBody))

		By("checking that the management user has been deleted")
		_, userErr = rmqClient.GetUser(managementUsername)
		Expect(userErr).To(MatchError(ContainSubstring("Error 404")))

		By("checking that the vhost has been deleted")
		_, err = rmqClient.GetVhost(serviceInstanceID)
		rabbitErr := err.(rabbithole.ErrorResponse)
		Expect(rabbitErr.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("succeeds for many SIs and bindings", func() {
		serviceInstanceIDa := "7ef8c450-6aa1-4d8a-a56f-a8bb4ddd1ae7"
		serviceInstanceIDb := "7ef8c450-6aa1-4d8a-a56f-a8bb4ddd1be7"
		serviceInstanceIDc := "7ef8c450-6aa1-4d8a-a56f-a8bb4ddd1ce7"
		bindingIDa := "7ef8c450-8b1c-4d8a-a56f-a8bb4ddd1ae7"
		bindingIDb := "7ef8c450-8b1c-4d8a-a56f-a8bb4ddd1be7"
		bindingIDc := "7ef8c450-8b1c-4d8a-a56f-a8bb4ddd1ce7"
		bindingIDd := "7ef8c450-8b1c-4d8a-a56f-a8bb4ddd1de7"
		bindingIDe := "7ef8c450-8b1c-4d8a-a56f-a8bb4ddd1ee7"
		bindingIDf := "7ef8c450-8b1c-4d8a-a56f-a8bb4ddd1fe7"

		By("sending a provision requests")
		provisionResponse, provisionBody := provision(serviceInstanceIDa, serviceID, planID)
		Expect(provisionResponse.StatusCode).To(Equal(http.StatusCreated), string(provisionBody))
		provisionResponse, provisionBody = provision(serviceInstanceIDb, serviceID, planID)
		Expect(provisionResponse.StatusCode).To(Equal(http.StatusCreated), string(provisionBody))
		provisionResponse, provisionBody = provision(serviceInstanceIDc, serviceID, planID)
		Expect(provisionResponse.StatusCode).To(Equal(http.StatusCreated), string(provisionBody))

		By("sending a binding requests")
		bindResponse1, bindBody1 := bind(bindingIDa, serviceInstanceIDa, serviceID, planID)
		Expect(bindResponse1.StatusCode).To(Equal(http.StatusCreated), string(bindBody1))
		bindResponse2, bindBody2 := bind(bindingIDb, serviceInstanceIDa, serviceID, planID)
		Expect(bindResponse2.StatusCode).To(Equal(http.StatusCreated), string(bindBody2))
		bindResponse3, bindBody3 := bind(bindingIDc, serviceInstanceIDa, serviceID, planID)
		Expect(bindResponse3.StatusCode).To(Equal(http.StatusCreated), string(bindBody3))
		bindResponse4, bindBody4 := bind(bindingIDd, serviceInstanceIDb, serviceID, planID)
		Expect(bindResponse4.StatusCode).To(Equal(http.StatusCreated), string(bindBody4))
		bindResponse5, bindBody5 := bind(bindingIDe, serviceInstanceIDb, serviceID, planID)
		Expect(bindResponse5.StatusCode).To(Equal(http.StatusCreated), string(bindBody5))
		bindResponse6, bindBody6 := bind(bindingIDf, serviceInstanceIDb, serviceID, planID)
		Expect(bindResponse6.StatusCode).To(Equal(http.StatusCreated), string(bindBody6))

		By("checking the binding passwords are different")
		var binding1, binding2 map[string]interface{}
		Expect(json.Unmarshal(bindBody1, &binding1)).To(Succeed())
		Expect(json.Unmarshal(bindBody2, &binding2)).To(Succeed())
		Expect(binding1["credentials"].(map[string]interface{})["password"]).NotTo(Equal(binding2["credentials"].(map[string]interface{})["password"]))

		By("checking that that binding user does not have access to another instance's vhost")
		_, err := rmqClient.GetPermissionsIn(serviceInstanceIDa, bindingIDd)
		Expect(err).To(HaveOccurred())

		By("sending a deprovision request")
		deprovisionResponse1, deprovisionBody1 := doRequest(http.MethodDelete, deprovisionURL(serviceInstanceIDa, serviceID, planID), nil)
		Expect(deprovisionResponse1.StatusCode).To(Equal(http.StatusOK), string(deprovisionBody1))
		deprovisionResponse2, deprovisionBody2 := doRequest(http.MethodDelete, deprovisionURL(serviceInstanceIDb, serviceID, planID), nil)
		Expect(deprovisionResponse2.StatusCode).To(Equal(http.StatusOK), string(deprovisionBody2))
		deprovisionResponse3, deprovisionBody3 := doRequest(http.MethodDelete, deprovisionURL(serviceInstanceIDc, serviceID, planID), nil)
		Expect(deprovisionResponse3.StatusCode).To(Equal(http.StatusOK), string(deprovisionBody3))
	})

	Context("errors", func() {
		Context("provisioning", func() {
			When("a service instance with the same ID has already been provisioned", func() {
				const serviceInstanceID = "9ef8c450-6aa1-6d8a-b56f-a8bb4ddd1de4"

				BeforeEach(func() {
					response, body := provision(serviceInstanceID, serviceID, planID)
					Expect(response.StatusCode).To(Equal(http.StatusCreated), string(body))
				})

				It("fails with HTTP 409", func() {
					response, _ := provision(serviceInstanceID, serviceID, planID)
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

		Context("binding", func() {
			When("attempting to create the same binding more than once", func() {
				serviceInstanceID := "7ef8c450-6aa1-4d8a-a56f-a8bb4ddd1de5"

				BeforeEach(func() {
					response, body := provision(serviceInstanceID, serviceID, planID)
					Expect(response.StatusCode).To(Equal(http.StatusCreated), string(body))
					bind(bindingID, serviceInstanceID, serviceID, planID)
				})

				It("fails with a 409", func() {
					resentBindResponse, resentBindBody := bind(bindingID, serviceInstanceID, serviceID, planID)

					Expect(resentBindResponse.StatusCode).To(Equal(http.StatusConflict), string(resentBindBody))
				})
			})
		})

		Context("unbinding", func() {
			When("a binding doesn't exist", func() {
				It("fails with HTTP 410", func() {
					serviceInstanceID := "does-not-exist"
					response, body := doRequest(http.MethodDelete, unbindURL(bindingID, serviceInstanceID, serviceID, planID), nil)
					Expect(response.StatusCode).To(Equal(http.StatusGone), string(body))
				})
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

func unbindURL(serviceInstanceID, bindingID, serviceID, planID string) string {
	return fmt.Sprintf("%s?service_id=%s&plan_id=%s", bindURL(serviceInstanceID, bindingID), serviceID, planID)
}

func provision(serviceInstanceID, serviceID, planID string) (*http.Response, []byte) {
	provisionDetails, err := json.Marshal(map[string]string{
		"service_id":        serviceID,
		"plan_id":           planID,
		"organization_guid": "fake-org-guid",
		"space_guid":        "fake-space-guid",
	})
	Expect(err).NotTo(HaveOccurred())

	return doRequest(http.MethodPut, provisionURL(serviceInstanceID), bytes.NewReader(provisionDetails))
}

func bind(bindingID, serviceInstanceID, serviceID, planID string) (*http.Response, []byte) {
	bindDetails, err := json.Marshal(map[string]string{
		"service_id": serviceID,
		"plan_id":    planID,
	})
	Expect(err).NotTo(HaveOccurred())

	return doRequest(http.MethodPut, bindURL(serviceInstanceID, bindingID), bytes.NewReader(bindDetails))
}
