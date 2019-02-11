package broker_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"rabbitmq-service-broker/broker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var config *broker.Config

	Describe("ReadConfig", func() {
		It("reads the config from file", func() {
			path, err := filepath.Abs(filepath.Join("..", "integrationtests", "fixtures", "config.yml"))
			Expect(err).NotTo(HaveOccurred())
			config, err = broker.ReadConfig(path)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.ServiceConfig.Username).To(Equal("p1-rabbit"))
			Expect(config.RabbitmqConfig.Administrator.Username).To(Equal("guest"))
			Expect(config.RabbitmqConfig.Policy.Name).To(Equal("operator_set_policy"))
			Expect(config.RabbitmqConfig.Policy.Definition["ha-mode"]).To(Equal("exactly"))
		})

		Context("when the config is not in the correct format", func() {
			It("returns an error", func() {
				tmpfile, err := ioutil.TempFile("", "wrong-config.yml")
				Expect(err).NotTo(HaveOccurred())
				fmt.Fprintf(tmpfile, "this is wrong content")
				tmpfile.Seek(0, os.SEEK_SET)
				path, err := filepath.Abs(tmpfile.Name())
				Expect(err).NotTo(HaveOccurred())

				config, err := broker.ReadConfig(path)
				Expect(config).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the config is empty", func() {
			It("returns an error", func() {
				tmpfile, err := ioutil.TempFile("", "empty-config.yml")
				Expect(err).NotTo(HaveOccurred())
				path, err := filepath.Abs(tmpfile.Name())
				Expect(err).NotTo(HaveOccurred())

				config, err := broker.ReadConfig(path)
				Expect(config).To(BeNil())
				Expect(err).To(HaveOccurred())
			})

		})

		Context("when the file does not exist", func() {
			It("returns the error", func() {
				config, err := broker.ReadConfig("this-is-missing")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("open this-is-missing: no such file or directory"))
				Expect(config).To(BeNil())
			})
		})
	})

	Describe("ValidateConfig", func() {
		var (
			config *broker.Config
		)
		BeforeEach(func() {
			path, err := filepath.Abs(filepath.Join("..", "integrationtests", "fixtures", "config.yml"))
			Expect(err).ToNot(HaveOccurred())
			config, err = broker.ReadConfig(path)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns nil if the config is valid", func() {
			err := broker.ValidateConfig(config)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error when it has an empty UUID", func() {
			config.ServiceConfig.UUID = ""
			err := broker.ValidateConfig(config)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error when it has an empty service name", func() {
			config.ServiceConfig.Name = ""
			err := broker.ValidateConfig(config)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error when it has an empty service username", func() {
			config.ServiceConfig.Username = ""
			err := broker.ValidateConfig(config)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error when it has an empty service password", func() {
			config.ServiceConfig.Password = ""
			err := broker.ValidateConfig(config)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error when it has an empty plan UUID", func() {
			config.ServiceConfig.PlanUuid = ""
			err := broker.ValidateConfig(config)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error when it has an empty hosts", func() {
			config.RabbitmqConfig.Hosts = []string{}
			err := broker.ValidateConfig(config)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error when it has an empty administrator username", func() {
			config.RabbitmqConfig.Administrator.Username = ""
			err := broker.ValidateConfig(config)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error when it has an empty administrator password", func() {
			config.RabbitmqConfig.Administrator.Password = ""
			err := broker.ValidateConfig(config)
			Expect(err).To(HaveOccurred())
		})
	})

})
