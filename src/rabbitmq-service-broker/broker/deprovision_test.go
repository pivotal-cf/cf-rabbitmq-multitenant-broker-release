package broker_test

import (
	"context"
	"errors"

	"rabbitmq-service-broker/rabbithutch/fakes"

	"github.com/pivotal-cf/brokerapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deprovisioning a RMQ service instance", func() {
	var (
		rabbithutch *fakes.FakeRabbitHutch
		broker      brokerapi.ServiceBroker
		ctx         context.Context
	)

	BeforeEach(func() {
		rabbithutch = &fakes.FakeRabbitHutch{}
	})

	When("the instance exists", func() {
		BeforeEach(func() {
			broker = defaultServiceBroker(defaultConfig(), rabbithutch)
			rabbithutch.VHostExistsReturns(true, nil)
			ctx = context.TODO()
		})

		It("deletes a vhost", func() {
			spec, err := broker.Deprovision(ctx, "my-service-instance-id", brokerapi.DeprovisionDetails{}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.IsAsync).To(BeFalse())

			Expect(rabbithutch.VHostDeleteCallCount()).To(Equal(1))
			Expect(rabbithutch.VHostDeleteArgsForCall(0)).To(Equal("my-service-instance-id"))
		})

		It("fails if it cannot delete the vhost", func() {
			rabbithutch.VHostDeleteReturns(errors.New("fake failure to delete vhost"))
			_, err := broker.Deprovision(ctx, "my-service-instance-id", brokerapi.DeprovisionDetails{}, false)
			Expect(err).To(MatchError("fake failure to delete vhost"))
		})

		It("deletes any management users if they exists", func() {
			rabbithutch.UserListReturns([]string{
				"mu-not-my-service-instance-asdasdasd",
				"mu-my-service-instance-id-qweqweqwe",
				"mu-my-service-instance-id-lfdsahjhh",
			}, nil)

			spec, err := broker.Deprovision(ctx, "my-service-instance-id", brokerapi.DeprovisionDetails{}, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(spec.IsAsync).To(BeFalse())

			Expect(rabbithutch.DeleteUserCallCount()).To(Equal(2))
			Expect(rabbithutch.DeleteUserArgsForCall(0)).To(Equal("mu-my-service-instance-id-qweqweqwe"))
			Expect(rabbithutch.DeleteUserArgsForCall(1)).To(Equal("mu-my-service-instance-id-lfdsahjhh"))
		})

		When("the management user doesn't exist", func() {
			It("does not try to delete it", func() {
				rabbithutch.UserListReturns([]string{"mu-not-my-service-instance-asdasdasd"}, nil)

				_, err := broker.Deprovision(ctx, "my-service-instance-id", brokerapi.DeprovisionDetails{}, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(rabbithutch.DeleteUserCallCount()).To(Equal(0))
			})
		})

		When("deleting the management user fails", func() {
			It("returns an error", func() {
				rabbithutch.UserListReturns([]string{"mu-my-service-instance-id-qweqweqwe"}, nil)
				rabbithutch.DeleteUserReturns(errors.New("fake delete user error"))

				_, err := broker.Deprovision(ctx, "my-service-instance-id", brokerapi.DeprovisionDetails{}, false)
				Expect(err).To(MatchError("fake delete user error"))
			})
		})
	})

	When("the SI does not exist", func() {
		BeforeEach(func() {
			rabbithutch.VHostExistsReturns(false, nil)
			broker = defaultServiceBroker(defaultConfig(), rabbithutch)
			ctx = context.TODO()
		})

		It("returns an error if vhost does not exist", func() {
			_, err := broker.Deprovision(ctx, "my-service-instance-id", brokerapi.DeprovisionDetails{}, false)
			Expect(err).To(MatchError(brokerapi.ErrInstanceDoesNotExist))
		})
	})

	When("we fail to query the vhost", func() {
		BeforeEach(func() {
			rabbithutch.VHostExistsReturns(false, errors.New("fake failure to query vhost"))
			broker = defaultServiceBroker(defaultConfig(), rabbithutch)
			ctx = context.TODO()
		})

		It("fails with an error saying the vhost could not be retrieved", func() {
			_, err := broker.Deprovision(ctx, "my-service-instance-id", brokerapi.DeprovisionDetails{}, false)
			Expect(err).To(MatchError("fake failure to query vhost"))
		})
	})
})
