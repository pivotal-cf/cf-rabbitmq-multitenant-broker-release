package rabbithutch_test

import (
	"net/http"

	"rabbitmq-service-broker/rabbithutch/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "rabbitmq-service-broker/rabbithutch"
)

var _ = Describe("UserPermissions", func() {

	var (
		rabbitClient *fakes.FakeAPIClient
		rabbithutch  RabbitHutch
		body         *fakeBody
	)

	BeforeEach(func() {
		rabbitClient = new(fakes.FakeAPIClient)
		rabbithutch = New(rabbitClient)
		body = &fakeBody{}
	})

	AfterEach(func() {
		Expect(body.Closed).To(BeTrue())
	})

	It("grants permissions on the vhost  user", func() {
		rabbitClient.UpdatePermissionsInReturns(&http.Response{StatusCode: http.StatusNoContent, Body: body}, nil)
		err := rabbithutch.AssignPermissionsTo("fake-vhost", "fake-user")
		Expect(err).NotTo(HaveOccurred())

		Expect(rabbitClient.UpdatePermissionsInCallCount()).To(Equal(1))
		vhost, username, permissions := rabbitClient.UpdatePermissionsInArgsForCall(0)
		Expect(vhost).To(Equal("fake-vhost"))
		Expect(username).To(Equal("fake-user"))
		Expect(permissions.Configure).To(Equal(".*"))
		Expect(permissions.Read).To(Equal(".*"))
		Expect(permissions.Write).To(Equal(".*"))
	})

	When("granting permissions fails", func() {
		BeforeEach(func() {
			rabbitClient.UpdatePermissionsInReturns(&http.Response{StatusCode: http.StatusForbidden, Body: body}, nil)
		})

		It("returns an error", func() {
			err := rabbithutch.AssignPermissionsTo("fake-vhost", "fake-user")
			Expect(err).To(MatchError("http request failed with status code: 403"))
		})
	})

})
