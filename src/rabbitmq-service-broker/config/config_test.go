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

	It("fails when many required fields are missing", func() {
		_, err := Read(fixture("missing-everything.yml"))
		Expect(err).To(MatchError(ContainSubstring("missing fields:")))
		Expect(err).To(MatchError(ContainSubstring("service.username")))
		Expect(err).To(MatchError(ContainSubstring("service.password")))
		Expect(err).To(MatchError(ContainSubstring("service.name")))
		Expect(err).To(MatchError(ContainSubstring("service.plan_uuid")))
		Expect(err).To(MatchError(ContainSubstring("service.uuid")))
		Expect(err).To(MatchError(ContainSubstring("rabbitmq.host")))
		Expect(err).To(MatchError(ContainSubstring("rabbitmq.administrator.username")))
		Expect(err).To(MatchError(ContainSubstring("rabbitmq.administrator.password")))
	})

	It("fails when the list of hosts is empty", func() {
		_, err := Read(fixture("empty-hosts.yml"))
		Expect(err).To(MatchError("Config file has missing fields: rabbitmq.hosts"))
	})

	It("interprets an empty field as disabling TLS", func() {
		conf, err := Read(fixture("ssl-empty.yml"))
		Expect(err).NotTo(HaveOccurred())
		Expect(conf.RabbitMQ.TLS).To(BeEquivalentTo(false))
	})
})
