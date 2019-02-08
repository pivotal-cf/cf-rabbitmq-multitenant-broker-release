package broker_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"rabbitmq-service-broker/broker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var config *broker.Config

	Describe("ParseConfig", func() {
		It("parses the config", func() {
			path, err := filepath.Abs(filepath.Join("..", "integrationtests", "fixtures", "config.yml"))
			Expect(err).ToNot(HaveOccurred())
			file, err := os.Open(path)
			Expect(err).ToNot(HaveOccurred())
			config, err = broker.ParseConfig(file)
			Expect(err).ToNot(HaveOccurred())

			Expect(config.ServiceConfig.Username).To(Equal("p1-rabbit"))
		})

		It("returns an error if config is empty", func() {
			_, err := broker.ParseConfig(strings.NewReader("---\n"))
			Expect(err).To(HaveOccurred())
		})

		Context("when the configuration is not the right format", func() {
			It("returns the error", func() {
				tmpfile, err := ioutil.TempFile("", "wrong-config.yml")
				Expect(err).NotTo(HaveOccurred())
				fmt.Fprintf(tmpfile, "this is wrong content")
				tmpfile.Seek(0, os.SEEK_SET)

				config, err := broker.ParseConfig(tmpfile)
				Expect(config).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})
	})

	//	Describe("ReadConfig", func() {
	//		It("reads the config from file", func() {
	//			path, err := filepath.Abs(filepath.Join("assets", "valid.yml"))
	//			Expect(err).NotTo(HaveOccurred())
	//			config, err := broker.ReadConfig(path)
	//			Expect(err).NotTo(HaveOccurred())
	//
	//			Expect(config.ServiceConfig.Username).To(Equal("admin"))
	//			Expect(config.Rabbitmq.Administrator.Username).To(Equal("guest"))
	//			Expect(config.Rabbitmq.Policy.Name).To(Equal("operator_set_policy"))
	//			Expect(config.Rabbitmq.Policy.Definition["ha-mode"]).To(Equal("exactly"))
	//		})
	//
	//		Context("when the file does not exist", func() {
	//			It("returns the error", func() {
	//				config, err := broker.ReadConfig("this-is-missing")
	//				Expect(err).To(HaveOccurred())
	//				Expect(err.Error()).To(Equal("open this-is-missing: no such file or directory"))
	//				Expect(config).To(BeNil())
	//			})
	//		})
	//	})
	//
	//	Describe("ValidateConfig", func() {
	//		var (
	//			config *broker.Config
	//		)
	//
	//		BeforeEach(func() {
	//			path, err := filepath.Abs(filepath.Join("assets", "valid.yml"))
	//			Expect(err).ToNot(HaveOccurred())
	//			config, err = broker.ReadConfig(path)
	//			Expect(err).ToNot(HaveOccurred())
	//		})
	//
	//		It("returns nil if the config is valid", func() {
	//			err := broker.ValidateConfig(config)
	//			Expect(err).ToNot(HaveOccurred())
	//		})
	//
	//		It("returns an error when it has an empty UUID", func() {
	//			config.ServiceConfig.Uuid = ""
	//			err := broker.ValidateConfig(config)
	//			Expect(err).To(HaveOccurred())
	//		})
	//
	//		It("returns an error when it has an empty service name", func() {
	//			config.ServiceConfig.Name = ""
	//			err := broker.ValidateConfig(config)
	//			Expect(err).To(HaveOccurred())
	//		})
	//
	//		It("returns an error when it has an empty service username", func() {
	//			config.ServiceConfig.Username = ""
	//			err := broker.ValidateConfig(config)
	//			Expect(err).To(HaveOccurred())
	//		})
	//
	//		It("returns an error when it has an empty service password", func() {
	//			config.ServiceConfig.Password = ""
	//			err := broker.ValidateConfig(config)
	//			Expect(err).To(HaveOccurred())
	//		})
	//
	//		It("returns an error when it has an empty plan UUID", func() {
	//			config.ServiceConfig.PlanUuid = ""
	//			err := broker.ValidateConfig(config)
	//			Expect(err).To(HaveOccurred())
	//		})
	//
	//		It("returns an error when it has an empty hosts", func() {
	//			config.Rabbitmq.Hosts = []string{}
	//			err := broker.ValidateConfig(config)
	//			Expect(err).To(HaveOccurred())
	//		})
	//
	//		It("returns an error when it has an empty administrator username", func() {
	//			config.Rabbitmq.Administrator.Username = ""
	//			err := broker.ValidateConfig(config)
	//			Expect(err).To(HaveOccurred())
	//		})
	//
	//		It("returns an error when it has an empty administrator password", func() {
	//			config.Rabbitmq.Administrator.Password = ""
	//			err := broker.ValidateConfig(config)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})

})
