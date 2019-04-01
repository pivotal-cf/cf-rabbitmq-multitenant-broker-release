package broker_test

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"rabbitmq-service-broker/rabbithutch/fakes"

	"github.com/pivotal-cf/brokerapi"
)

var _ = Describe("Binding a RMQ service instance", func() {
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

	When("the SI does not exist", func() {
		BeforeEach(func() {
			rabbithutch.VHostExistsReturns(false, nil)
		})

		It("fails with an error saying the SI does not exist", func() {
			_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
			Expect(err).To(MatchError(brokerapi.ErrInstanceDoesNotExist))
		})
	})

	When("we fail to query the vhost", func() {
		BeforeEach(func() {
			rabbithutch.VHostExistsReturns(false, errors.New("fake vhost error"))
		})

		It("fails with an error saying the vhost could not be retrieved", func() {
			_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
			Expect(err).To(MatchError("fake vhost error"))
		})
	})

	When("the SI exists", func() {
		BeforeEach(func() {
			rabbithutch.VHostExistsReturns(true, nil)
		})

		Describe("the user", func() {
			It("creates a user", func() {
				rabbithutch.CreateUserAndGrantPermissionsReturns("fake-password", nil)
				_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(rabbithutch.CreateUserAndGrantPermissionsCallCount()).To(Equal(1))
				username, vhost, tags := rabbithutch.CreateUserAndGrantPermissionsArgsForCall(0)
				Expect(username).To(Equal("binding-id"))
				Expect(vhost).To(Equal("my-service-instance-id"))
				Expect(tags).To(Equal(""))
			})

			It("fails with an error if it cannot create a user", func() {
				rabbithutch.CreateUserAndGrantPermissionsReturns("fake-password", errors.New("foo"))
				_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
				Expect(err).To(MatchError("foo"))
			})

			When("user tags are set in the config", func() {
				BeforeEach(func() {
					broker = defaultServiceBroker(defaultConfigWithUserTags(), rabbithutch)
					ctx = context.TODO()
				})

				It("creates a user with the tags", func() {
					rabbithutch.CreateUserAndGrantPermissionsReturns("fake-password", nil)
					_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)

					Expect(err).NotTo(HaveOccurred())

					Expect(rabbithutch.CreateUserAndGrantPermissionsCallCount()).To(Equal(1))
					username, vhost, tags := rabbithutch.CreateUserAndGrantPermissionsArgsForCall(0)
					Expect(username).To(Equal("binding-id"))
					Expect(vhost).To(Equal("my-service-instance-id"))
					Expect(tags).To(Equal("administrator"))
				})
			})
		})

		Describe("the binding data", func() {
			When("it cannot read the protocol ports", func() {
				It("fails with an error", func() {
					rabbithutch.ProtocolPortsReturns(nil, fmt.Errorf("failed to read protocol ports"))
					_, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)

					Expect(err).To(MatchError("failed to read protocol ports"))
				})
			})

			When("it reads the protocol ports", func() {
				BeforeEach(func() {
					rabbithutch.ProtocolPortsReturns(fakeProtocolPorts(), nil)
					rabbithutch.CreateUserAndGrantPermissionsReturns("fake-password", nil)
				})

				It("generates the right binding", func() {
					binding, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
					Expect(err).NotTo(HaveOccurred())

					Expect(binding.Credentials).To(HaveKeyWithValue("username", "binding-id"))
					Expect(binding.Credentials).To(HaveKeyWithValue("password", "fake-password"))
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

			When("using an external load balancer", func() {
				BeforeEach(func() {
					rabbithutch.ProtocolPortsReturns(fakeProtocolPorts(), nil)
					rabbithutch.CreateUserAndGrantPermissionsReturns("fake-password", nil)
				})

				It("uses the right hosts", func() {
					broker = defaultServiceBroker(defaultConfigWithExternalLoadBalancer(), rabbithutch)
					binding, err := broker.Bind(ctx, "my-service-instance-id", "binding-id", brokerapi.BindDetails{}, false)
					Expect(err).NotTo(HaveOccurred())

					Expect(binding.Credentials).To(HaveKeyWithValue("hostnames", ConsistOf("my-dns-host.com")))
				})
			})
		})
	})
})

func fakeProtocolPorts() map[string]int {
	return map[string]int{
		"amqp/ssl":   5671,
		"clustering": 25672,
		"http":       15672,
		"mqtt":       1883,
		"mqtt/ssl":   8883,
		"stomp":      61613,
		"stomp/ssl":  61614,
	}
}
