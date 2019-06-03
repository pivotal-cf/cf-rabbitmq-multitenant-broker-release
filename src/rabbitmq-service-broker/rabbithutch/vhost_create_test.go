package rabbithutch_test

import (
	"fmt"
	"net/http"
	"rabbitmq-service-broker/rabbithutch/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "rabbitmq-service-broker/rabbithutch"
)

var _ = Describe("VhostCreate", func() {
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

	It("creates a vhost", func() {
		rabbitClient.PutVhostReturns(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

		err := rabbithutch.VHostCreate("fake-vhost")
		Expect(err).NotTo(HaveOccurred())

		Expect(rabbitClient.PutVhostCallCount()).To(Equal(1))
		Expect(rabbitClient.PutVhostArgsForCall(0)).To(Equal("fake-vhost"))
		Expect(body.Closed).To(BeTrue())
	})

	When("the vhost creation fails", func() {
		It("returns an error when the RMQ API returns an error", func() {
			rabbitClient.PutVhostReturns(nil, fmt.Errorf("vhost-creation-failed"))

			err := rabbithutch.VHostCreate("fake-vhost")
			Expect(err).To(MatchError("vhost-creation-failed"))
		})

		When("the rabbit client successfully responds", func() {

			It("returns an error when the RMQ API returns a bad HTTP response code", func() {
				rabbitClient.PutVhostReturns(&http.Response{StatusCode: http.StatusInternalServerError, Body: body}, nil)

				err := rabbithutch.VHostCreate("fake-vhost")

				Expect(err).To(MatchError("http request failed with status code: 500"))
				Expect(body.Closed).To(BeTrue())
			})
		})
	})
})
