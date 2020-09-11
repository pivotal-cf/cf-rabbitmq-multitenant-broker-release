package management_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"rabbitmq-service-broker/config"
	. "rabbitmq-service-broker/management"
)

var _ = Describe("NewClient", func() {
	var brokerConfig config.Config
	When("Management over TLS is not configured", func() {
		BeforeEach(func() {
			brokerConfig = config.Config{RabbitMQ: config.RabbitMQ{
				Hosts: []string{"127.0.0.1", "127.0.0.2"},
				Administrator: config.AdminCredentials{
					Username: "alexandra",
					Password: "ardnaxela",
				},
			}}
		})
		It("produces a RabbitMQ management client pointing at the non-TLS endpoint", func() {
			client, err := NewClient(brokerConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.Endpoint).To(Equal("http://127.0.0.1:15672"))
			Expect(client.Username).To(Equal("alexandra"))
			Expect(client.Password).To(Equal("ardnaxela"))
		})
	})

	When("Management over TLS is configured", func() {
		BeforeEach(func() {
			brokerConfig = config.Config{RabbitMQ: config.RabbitMQ{
				Hosts: []string{"127.0.0.1", "127.0.0.2"},
				Administrator: config.AdminCredentials{
					Username: "alexandra",
					Password: "ardnaxela",
				},
				ManagementTLS: config.ManagementTLS{
					Enabled: true,
				},
			}}
		})
		It("produces a RabbitMQ management client pointing at the non-TLS endpoint", func() {
			client, err := NewClient(brokerConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.Endpoint).To(Equal("https://127.0.0.1:15671"))
			Expect(client.Username).To(Equal("alexandra"))
			Expect(client.Password).To(Equal("ardnaxela"))
		})
	})

	When("The broker is configured with an external load balancer", func() {
		When("Management over TLS is not configured", func() {
			BeforeEach(func() {
				brokerConfig = config.Config{RabbitMQ: config.RabbitMQ{
					Hosts:   []string{"127.0.0.1", "127.0.0.2"},
					DNSHost: "abc.123.com",
					Administrator: config.AdminCredentials{
						Username: "alexandra",
						Password: "ardnaxela",
					},
				}}
			})
			It("uses the URL of the load balancer in the client endpoint", func() {
				client, err := NewClient(brokerConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(client.Endpoint).To(Equal("http://abc.123.com:15672"))
				Expect(client.Username).To(Equal("alexandra"))
				Expect(client.Password).To(Equal("ardnaxela"))
			})
		})

		When("Management over TLS is configured", func() {
			BeforeEach(func() {
				brokerConfig = config.Config{RabbitMQ: config.RabbitMQ{
					Hosts:   []string{"127.0.0.1", "127.0.0.2"},
					DNSHost: "abc.123.com",
					Administrator: config.AdminCredentials{
						Username: "alexandra",
						Password: "ardnaxela",
					},
					ManagementTLS: config.ManagementTLS{
						Enabled: true,
					},
				}}
			})
			It("uses the URL of the load balancer in the client endpoint", func() {
				client, err := NewClient(brokerConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(client.Endpoint).To(Equal("https://abc.123.com:15671"))
				Expect(client.Username).To(Equal("alexandra"))
				Expect(client.Password).To(Equal("ardnaxela"))
			})
		})
	})

})
