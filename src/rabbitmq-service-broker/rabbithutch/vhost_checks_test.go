package rabbithutch_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "rabbitmq-service-broker/rabbithutch"
	"rabbitmq-service-broker/rabbithutch/fakes"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Binding a RMQ service instance", func() {
	var (
		rabbitClient *fakes.FakeAPIClient
		rabbithutch  RabbitHutch
	)

	BeforeEach(func() {
		rabbitClient = new(fakes.FakeAPIClient)
		rabbithutch = New(rabbitClient)
	})

	Describe("EnsureVHostExists()", func() {
		AfterEach(func() {
			Expect(rabbitClient.GetVhostArgsForCall(0)).To(Equal("fake-vhost"))
		})

		When("the vhost does not exist", func() {
			BeforeEach(func() {
				rabbitClient.GetVhostReturns(nil, rabbithole.ErrorResponse{StatusCode: http.StatusNotFound})
			})

			It("fails with an error saying the service instance does not exist", func() {
				err := rabbithutch.EnsureVHostExists("fake-vhost")
				Expect(err).To(MatchError(brokerapi.ErrInstanceDoesNotExist))
			})
		})

		When("we fail to query the vhost", func() {
			BeforeEach(func() {
				rabbitClient.GetVhostReturns(nil, rabbithole.ErrorResponse{StatusCode: http.StatusInternalServerError})
			})

			It("fails with an error saying the vhost could not be retrieved", func() {
				err := rabbithutch.EnsureVHostExists("fake-vhost")
				Expect(err).To(MatchError(rabbithole.ErrorResponse{StatusCode: http.StatusInternalServerError}))
			})
		})

		When("the vhost exists", func() {
			BeforeEach(func() {
				rabbitClient.GetVhostReturns(&rabbithole.VhostInfo{}, nil)
			})

			It("returns nil", func() {
				err := rabbithutch.EnsureVHostExists("fake-vhost")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
