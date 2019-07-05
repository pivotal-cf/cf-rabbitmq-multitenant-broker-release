package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "rabbitmq-service-broker/config"
)

var _ = Describe("Config", func() {
	It("reads a minimal config", func() {
		conf, err := Read(fixture("minimal.yml"))
		Expect(err).NotTo(HaveOccurred())
		Expect(conf).To(Equal(Config{
			Service: Service{
				UUID:     "fake-service-uuid",
				Name:     "fake-service-name",
				PlanUUID: "fake-plan-uuid",
				Username: "fake-service-username",
				Password: "fake-service-password",
			},
			RabbitMQ: RabbitMQ{
				ManagementDomain: "fake-management-domain",
				Hosts:            []string{"fake-host-1", "fake-host-2"},
				Administrator:    AdminCredentials{"fake-rmq-user", "fake-rmq-password"},
				TLS:              false,
			},
		}))
	})

	It("reads a complete config", func() {
		conf, err := Read(fixture("complete.yml"))
		Expect(err).NotTo(HaveOccurred())
		Expect(conf).To(Equal(Config{
			Service: Service{
				UUID:                "00000000-0000-0000-0000-000000000000",
				Name:                "p-rabbitmq",
				PlanUUID:            "11111111-1111-1111-1111-111111111111",
				Username:            "p1-rabbit",
				Password:            "p1-rabbit-test",
				Description:         "this is a description",
				DisplayName:         "WhiteRabbitMQ",
				LongDescription:     "this is a long description",
				ProviderDisplayName: "SomeCompany",
				DocumentationURL:    "https://example.com",
				SupportURL:          "https://support.example.com",
				IconImage:           "image_icon_base64",
				Shareable:           true,
			},
			RabbitMQ: RabbitMQ{
				RegularUserTags:  "policymaker,management",
				Hosts:            []string{"127.0.0.1", "127.0.0.2"},
				Administrator:    AdminCredentials{"fake-rmq-user", "fake-rmq-password"},
				Management:       ManagementCredentials{"management-username"},
				ManagementDomain: "pivotal-rabbitmq.127.0.0.1",
				TLS:              true,
				OperatorSetPolicy: RabbitMQPolicy{
					Enabled: true,
					Name:    "operator_set_policy",
					Definition: PolicyDefinition{
						"ha-mode":      "exactly",
						"ha-params":    float64(2),
						"ha-sync-mode": "automatic",
					},
					Priority: 50,
				},
			},
		}))
	})

	It("fails when one required field is missing", func() {
		_, err := Read(fixture("missing-uuid.yml"))
		Expect(err).To(MatchError("Config file has missing fields: service.uuid"))
	})

	It("fails when the policy definition is empty", func() {
		_, err := Read(fixture("empty-policy-definition.yml"))
		Expect(err).To(MatchError("Config file has missing field: operator_set_policy.policy_definition must be provided when operator_set_policy.enabled is true"))
	})

	It("fails when many required fields are missing", func() {
		_, err := Read(fixture("missing-everything.yml"))
		Expect(err).To(MatchError(ContainSubstring("missing fields:")))
		Expect(err).To(MatchError(ContainSubstring("service.username")))
		Expect(err).To(MatchError(ContainSubstring("service.password")))
		Expect(err).To(MatchError(ContainSubstring("service.name")))
		Expect(err).To(MatchError(ContainSubstring("service.plan_uuid")))
		Expect(err).To(MatchError(ContainSubstring("service.uuid")))
		Expect(err).To(MatchError(ContainSubstring("rabbitmq.administrator.username")))
		Expect(err).To(MatchError(ContainSubstring("rabbitmq.administrator.password")))
	})

	Describe("hosts", func() {
		It("fails when both `hosts` and the `dns_host` are empty", func() {
			_, err := Read(fixture("empty-hosts.yml"))
			Expect(err).To(MatchError("Config file has missing fields: at least one of rabbitmq.hosts or rabbitmq.dns_host must be specified"))
		})

		When("an external load balancer hostname is specified", func() {
			It("does not fail when the list of hosts is empty", func() {
				_, err := Read(fixture("external-load-balancer.yml"))
				Expect(err).NotTo(HaveOccurred())
			})

			It("will use the external load balancer hostname in bindings", func() {
				conf, err := Read(fixture("external-load-balancer.yml"))
				Expect(err).NotTo(HaveOccurred())
				Expect(conf.NodeHosts()).To(Equal([]string{"my-dns-host.com"}))
			})

		})

		When("there is not an external load balancer", func() {
			It("uses rabbitmq.hosts", func() {
				conf, err := Read(fixture("minimal.yml"))
				Expect(err).NotTo(HaveOccurred())
				Expect(conf.NodeHosts()).To(Equal([]string{
					"fake-host-1",
					"fake-host-2",
				}))
			})
		})

		When("both hosts and an external load balancer hostname are specified", func() {
			It("favours load balancer `dns_host` over `hosts`", func() {
				conf, err := Read(fixture("hosts-and-external-load-balancer.yml"))
				Expect(err).NotTo(HaveOccurred())
				Expect(conf.NodeHosts()).To(Equal([]string{"my-dns-host.com"}))
			})
		})

		When("`hosts` is a comma-separated string of IPs (rather than a YAML list)", func() {
			It("reads it as a list of hosts", func() {
				conf, err := Read(fixture("comma-list-hosts.yml"))
				Expect(err).NotTo(HaveOccurred())
				Expect(conf.NodeHosts()).To(BeEquivalentTo([]string{"127.0.0.1", "127.0.0.2"}))
			})
		})
	})

	Describe("the `ssl` field", func() {
		It("interprets a `false` as disabling TLS", func() {
			conf, err := Read(fixture("ssl-false.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(conf.RabbitMQ.TLS).To(BeEquivalentTo(false))
		})

		It("interprets an empty field as disabling TLS", func() {
			conf, err := Read(fixture("ssl-empty.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(conf.RabbitMQ.TLS).To(BeEquivalentTo(false))
		})
		It("interprets `null` as disabling TLS", func() {
			conf, err := Read(fixture("ssl-null.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(conf.RabbitMQ.TLS).To(BeEquivalentTo(false))
		})

		It("interprets `true` as enabling TLS", func() {
			conf, err := Read(fixture("ssl-true.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(conf.RabbitMQ.TLS).To(BeEquivalentTo(true))
		})

		It("interprets a certificate as enabling TLS", func() {
			conf, err := Read(fixture("ssl-cert.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(conf.RabbitMQ.TLS).To(BeEquivalentTo(true))
		})
	})
})
