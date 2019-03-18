package broker_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"rabbitmq-service-broker/rabbithutch/fakes"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Binding a RMQ service instance", func() {
	var (
		rabbitClient *fakes.FakeAPIClient
		broker       brokerapi.ServiceBroker
		ctx          context.Context
		rabbithutch  *fakes.FakeRabbitHutch
	)

	BeforeEach(func() {
		rabbitClient = &fakes.FakeAPIClient{}
		rabbithutch = &fakes.FakeRabbitHutch{}
		broker = defaultServiceBroker(defaultConfig(), rabbitClient, rabbithutch)
		ctx = context.TODO()
		rabbitClient.UpdatePermissionsInReturns(&http.Response{StatusCode: http.StatusOK}, nil)
	})

	When("the SI does not exist", func() {
		BeforeEach(func() {
			rabbithutch.EnsureVHostExistsReturns(brokerapi.ErrInstanceDoesNotExist)
		})

		It("fails with an error saying the SI does not exist", func() {
			_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
			Expect(err).To(MatchError(brokerapi.ErrInstanceDoesNotExist))
		})
	})

	When("we fail to query the vhost", func() {
		BeforeEach(func() {
			rabbithutch.EnsureVHostExistsReturns(rabbithole.ErrorResponse{StatusCode: http.StatusInternalServerError})
		})

		It("fails with an error saying the vhost could not be retrieved", func() {
			_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
			Expect(err).To(MatchError(rabbithole.ErrorResponse{StatusCode: http.StatusInternalServerError}))
		})
	})

	When("the SI exists", func() {
		Describe("the user", func() {
			It("creates a user", func() {
				_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(rabbitClient.PutUserCallCount()).To(Equal(1))
				username, info := rabbitClient.PutUserArgsForCall(0)
				Expect(username).To(Equal("binding-id"))
				Expect(info.Password).To(MatchRegexp(`[a-zA-Z0-9\-_]{24}`))
				Expect(info.Tags).To(Equal("policymaker,management"))
			})

			It("fails with an error if it cannot create a user", func() {
				rabbitClient.PutUserReturns(nil, errors.New("foo"))
				_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
				Expect(err).To(MatchError("foo"))
			})

			It("fails with an error if the user already exists", func() {
				rabbitClient.PutUserReturns(&http.Response{StatusCode: http.StatusNoContent}, nil)
				_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
				Expect(err).To(MatchError(brokerapi.ErrBindingAlreadyExists))
			})

			It("grants the user full permissions to the vhost", func() {
				_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(rabbitClient.UpdatePermissionsInCallCount()).To(Equal(1))
				vhost, username, permissions := rabbitClient.UpdatePermissionsInArgsForCall(0)
				Expect(vhost).To(Equal("my-service-instance-id"))
				Expect(username).To(Equal("binding-id"))
				Expect(permissions.Configure).To(Equal(".*"))
				Expect(permissions.Read).To(Equal(".*"))
				Expect(permissions.Write).To(Equal(".*"))
			})

			When("user tags are set in the config", func() {
				BeforeEach(func() {
					rabbitClient = new(fakes.FakeAPIClient)
					broker = defaultServiceBroker(defaultConfigWithUserTags(), rabbitClient, rabbithutch)
					ctx = context.TODO()
					rabbitClient.UpdatePermissionsInReturns(&http.Response{StatusCode: http.StatusOK}, nil)
				})

				It("creates a user with the tags", func() {
					_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
					Expect(err).NotTo(HaveOccurred())

					Expect(rabbitClient.PutUserCallCount()).To(Equal(1))
					username, info := rabbitClient.PutUserArgsForCall(0)
					Expect(username).To(Equal("binding-id"))
					Expect(info.Password).To(MatchRegexp(`[a-zA-Z0-9\-_]{24}`))
					Expect(info.Tags).To(Equal("administrator"))
				})
			})
		})

		Describe("the binding data", func() {
			When("it cannot read the protocol ports", func() {
				BeforeEach(func() {
					rabbitClient.ProtocolPortsReturns(nil, fmt.Errorf("failed to read protocol ports"))
				})

				It("fails with an error", func() {
					_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
					Expect(err).To(MatchError("failed to read protocol ports"))
				})
			})

			When("it reads the protocol ports", func() {
				BeforeEach(func() {
					rabbitClient.ProtocolPortsReturns(fakeProtocolPorts(), nil)
				})

				It("generates the right binding", func() {
					binding, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
					Expect(err).NotTo(HaveOccurred())

					Expect(binding.Credentials).To(HaveKeyWithValue("username", "binding-id"))
					Expect(binding.Credentials).To(HaveKeyWithValue("vhost", "my-service-instance-id"))
					Expect(binding.Credentials).To(HaveKeyWithValue("ssl", false))
					Expect(binding.Credentials).To(HaveKeyWithValue("hostname", "fake-hostname-1"))
					Expect(binding.Credentials).To(HaveKeyWithValue("hostnames", ConsistOf("fake-hostname-1", "fake-hostname-2")))
					Expect(binding.Credentials).To(HaveKeyWithValue("protocols", SatisfyAll(
						HaveKey("amqp+ssl"),
						HaveKey("mqtt"),
						HaveKey("mqtt+ssl"),
						HaveKey("stomp"),
						HaveKey("stomp+ssl"),
						HaveKey("management"),
					)))
				})
			})
		})
	})
})

func fakeProtocolPorts() map[string]rabbithole.Port {
	return map[string]rabbithole.Port{
		"amqp/ssl":   5671,
		"clustering": 25672,
		"http":       15672,
		"mqtt":       1883,
		"mqtt/ssl":   8883,
		"stomp":      61613,
		"stomp/ssl":  61614,
	}
}
