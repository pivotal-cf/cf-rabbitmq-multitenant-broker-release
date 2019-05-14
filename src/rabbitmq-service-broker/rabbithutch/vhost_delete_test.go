package rabbithutch_test

import (
	"errors"
	"net/http"
	"rabbitmq-service-broker/rabbithutch/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "rabbitmq-service-broker/rabbithutch"
)

var _ = Describe("VhostDelete", func() {
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

	It("deletes a vhost and closes the body", func() {
		rabbitClient.DeleteVhostReturns(&http.Response{StatusCode: http.StatusNoContent, Body: body}, nil)
		err := rabbithutch.VHostDelete("my-vhost")

		Expect(err).NotTo(HaveOccurred())

		Expect(rabbitClient.DeleteVhostCallCount()).To(Equal(1))
		Expect(rabbitClient.DeleteVhostArgsForCall(0)).To(Equal("my-vhost"))
	})

	It("fails if it cannot delete the vhost", func() {
		rabbitClient.DeleteVhostReturns(&http.Response{StatusCode: http.StatusBadRequest, Body: body}, errors.New("fake failure to delete vhost"))
		err := rabbithutch.VHostDelete("my-vhost")

		Expect(err).To(MatchError("fake failure to delete vhost"))
	})
})
