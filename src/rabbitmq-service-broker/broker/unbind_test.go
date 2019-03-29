package broker_test

import (
	"context"
	"fmt"

	"rabbitmq-service-broker/rabbithutch/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Unbind", func() {

	var (
		broker      brokerapi.ServiceBroker
		ctx         context.Context
		rabbithutch *fakes.FakeRabbitHutch
	)

	BeforeEach(func() {
		rabbithutch = &fakes.FakeRabbitHutch{}
		broker = defaultServiceBroker(defaultConfig(), rabbithutch)
		ctx = context.TODO()
	})

	It("deletes a user and its connections", func() {
		spec, err := broker.Unbind(ctx, "my-service-instance-id", "binding-id", brokerapi.UnbindDetails{}, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec).To(Equal(brokerapi.UnbindSpec{}))
		Expect(rabbithutch.DeleteUserAndConnectionsCallCount()).To(Equal(1))
		Expect(rabbithutch.DeleteUserAndConnectionsArgsForCall(0)).To(Equal("binding-id"))
	})

	When("it fails to delete a user", func() {
		BeforeEach(func() {
			rabbithutch.DeleteUserAndConnectionsReturns(fmt.Errorf("fake-error"))
		})

		It("returns an error", func() {
			_, err := broker.Unbind(ctx, "my-service-instance-id", "binding-id", brokerapi.UnbindDetails{}, false)
			Expect(err).To(MatchError("fake-error"))
		})
	})
})
