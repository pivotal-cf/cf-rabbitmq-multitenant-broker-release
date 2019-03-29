package broker_test

import (
	"context"
	"errors"
	"fmt"

	"rabbitmq-service-broker/rabbithutch/fakes"

	"github.com/pivotal-cf/brokerapi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Provisioning a RMQ service instance", func() {
	var (
		rabbithutch *fakes.FakeRabbitHutch
		broker      brokerapi.ServiceBroker
		ctx         context.Context
	)

	BeforeEach(func() {
		rabbithutch = &fakes.FakeRabbitHutch{}
		broker = defaultServiceBroker(defaultConfig(), rabbithutch)
		ctx = context.TODO()
		rabbithutch.VHostCreateReturns(nil)
		rabbithutch.VHostDeleteReturns(nil)
		rabbithutch.AssignPermissionsToReturns(nil)
	})

	It("creates a vhost", func() {
		_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
		Expect(err).NotTo(HaveOccurred())

		Expect(rabbithutch.VHostCreateCallCount()).To(Equal(1))
		Expect(rabbithutch.VHostCreateArgsForCall(0)).To(Equal("my-service-instance-id"))
	})

	It("grants permissions on the vhost to the service broker RMQ admin user", func() {
		_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
		Expect(err).NotTo(HaveOccurred())

		Expect(rabbithutch.AssignPermissionsToCallCount()).To(Equal(1))
		vhost, username := rabbithutch.AssignPermissionsToArgsForCall(0)
		Expect(vhost).To(Equal("my-service-instance-id"))
		Expect(username).To(Equal("default-admin-username"))
	})

	It("returns the dashboard URL", func() {
		spec, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
		Expect(err).NotTo(HaveOccurred())

		Expect(spec.DashboardURL).To(Equal("https://foo.bar.com/#/login/"))
		Expect(spec.IsAsync).To(BeFalse())
		Expect(spec.OperationData).To(Equal(""))
	})

	When("granting permissions fails", func() {
		BeforeEach(func() {
			rabbithutch.AssignPermissionsToReturns(fmt.Errorf("fake-error"))
		})

		It("cleans up the vhost", func() {
			_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
			Expect(err).To(MatchError("fake-error"))
			Expect(rabbithutch.VHostDeleteCallCount()).To(Equal(1))
		})
	})

	Context("management user permissions", func() {
		When("the management user is set", func() {
			BeforeEach(func() {
				cfg := defaultConfig()
				cfg.RabbitMQ.Management.Username = "default-management-username"
				broker = defaultServiceBroker(cfg, rabbithutch)
			})

			It("grants permissions on the vhost to the management RMQ user", func() {
				_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(rabbithutch.AssignPermissionsToCallCount()).To(Equal(2))
				vhost, username := rabbithutch.AssignPermissionsToArgsForCall(0)
				Expect(vhost).To(Equal("my-service-instance-id"))
				Expect(username).To(Equal("default-admin-username"))
			})

			When("granting permissions fails", func() {
				It("ignores the error", func() {
					_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
					Expect(err).NotTo(HaveOccurred())
					Expect(rabbithutch.VHostDeleteCallCount()).To(Equal(0))
				})
			})
		})

		When("the management RMQ user has not been set", func() {
			It("does not attempt to grant it permissions", func() {
				_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(rabbithutch.AssignPermissionsToCallCount()).To(Equal(1))
			})
		})
	})

	Context("vhost policies", func() {
		When("vhost policies are configured", func() {
			BeforeEach(func() {
				cfg := defaultConfig()
				cfg.RabbitMQ.OperatorSetPolicy.Enabled = true
				cfg.RabbitMQ.OperatorSetPolicy.Name = "fake-policy-name"
				cfg.RabbitMQ.OperatorSetPolicy.Definition = map[string]interface{}{"fake-policy-key": "fake-policy-value"}
				cfg.RabbitMQ.OperatorSetPolicy.Priority = 42
				broker = defaultServiceBroker(cfg, rabbithutch)
			})

			It("sets policies for the new instance", func() {
				_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(rabbithutch.CreatePolicyCallCount()).To(Equal(1))
				vhost, policyName, priority, definition := rabbithutch.CreatePolicyArgsForCall(0)
				Expect(vhost).To(Equal("my-service-instance-id"))
				Expect(policyName).To(Equal("fake-policy-name"))
				Expect(definition).To(BeEquivalentTo(map[string]interface{}{"fake-policy-key": "fake-policy-value"}))
				Expect(priority).To(Equal(42))
			})

			When("setting policies fails", func() {
				BeforeEach(func() {
					rabbithutch.CreatePolicyReturns(fmt.Errorf("some-error"))
				})

				It("cleans up the vhost", func() {
					_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
					Expect(err).To(MatchError("some-error"))
					Expect(rabbithutch.VHostDeleteCallCount()).To(Equal(1))
				})
			})
		})

		When("there are no policies configured", func() {
			It("does not attempt to configure any", func() {
				_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(rabbithutch.CreatePolicyCallCount()).To(Equal(0))
			})
		})
	})

	When("the vhost already exists", func() {
		It("returns ErrInstanceAlreadyExists error", func() {
			rabbithutch.VHostExistsReturns(true, nil)

			_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
			Expect(err).To(MatchError(brokerapi.ErrInstanceAlreadyExists))

			Expect(rabbithutch.VHostExistsCallCount()).To(Equal(1))
			Expect(rabbithutch.VHostExistsArgsForCall(0)).To(Equal("my-service-instance-id"))
			Expect(rabbithutch.VHostCreateCallCount()).To(Equal(0))
		})
	})

	When("checking whether the VHost exists fails", func() {
		It("returns an error", func() {
			rabbithutch.VHostExistsReturns(false, errors.New("fake check vhost error"))

			_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)

			Expect(err).To(MatchError("fake check vhost error"))
		})
	})

	When("the vhost creation fails", func() {
		It("returns an error when the RMQ API returns an error", func() {
			rabbithutch.VHostCreateReturns(fmt.Errorf("vhost-creation-failed"))

			_, err := broker.Provision(ctx, "my-service-instance-id", brokerapi.ProvisionDetails{}, false)
			Expect(err).To(MatchError("vhost-creation-failed"))
		})
	})
})
