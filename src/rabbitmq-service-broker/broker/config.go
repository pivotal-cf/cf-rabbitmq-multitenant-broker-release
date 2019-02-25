package broker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ServiceConfig  ServiceConfig  `yaml:"service"`
	RabbitMQConfig RabbitMQConfig `yaml:"rabbitmq"`
}

type ServiceConfig struct {
	UUID                string `yaml:"uuid"`
	Name                string `yaml:"name"`
	OfferingDescription string `yaml:"offering_description"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	PlanUUID            string `yaml:"plan_uuid"`
	DisplayName         string `yaml:"display_name"`
	IconImage           string `yaml:"icon_image"`
	LongDescription     string `yaml:"long_description"`
	ProviderDisplayName string `yaml:"provider_display_name"`
	DocumentationURL    string `yaml:"documentation_url"`
	SupportURL          string `yaml:"support_url"`
	Shareable           bool   `yaml:"shareable"`
}

type RabbitMQConfig struct {
	Hosts            []string            `yaml:"hosts"`
	DNSHost          string              `yaml:"dns_host"`
	ManagementDomain string              `yaml:"management_domain"`
	Management       RabbitMQCredentials `yaml:"management"`
	Administrator    RabbitMQCredentials `yaml:"administrator"`
	Policy           RabbitMQPolicy      `yaml:"operator_set_policy"`
}

type RabbitMQCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type PolicyDefinition map[string]interface{}

type RabbitMQPolicy struct {
	Enabled    bool             `yaml:"enabled"`
	Name       string           `yaml:"policy_name"`
	Priority   int              `yaml:"policy_priority"`
	Definition PolicyDefinition `yaml:"policy_definition"`
}

func ReadConfig(path string) (Config, error) {
	configBytes, err := ioutil.ReadFile(filepath.FromSlash(path))
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	if err = yaml.Unmarshal(configBytes, &config); err != nil {
		return Config{}, err
	}

	if err := ValidateConfig(config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func ValidateConfig(config Config) error {
	if config.ServiceConfig.UUID == "" {
		return fmt.Errorf("uuid is not set")
	}
	if config.ServiceConfig.Name == "" {
		return fmt.Errorf("service name is not set")
	}
	if config.ServiceConfig.Username == "" {
		return fmt.Errorf("service username is not set")
	}
	if config.ServiceConfig.Password == "" {
		return fmt.Errorf("service password is not set")
	}
	if config.ServiceConfig.PlanUUID == "" {
		return fmt.Errorf("plan uuid is not set")
	}
	if len(config.RabbitMQConfig.Hosts) < 1 {
		return fmt.Errorf("no rabbitmq hosts were set")
	}
	if config.RabbitMQConfig.Administrator.Username == "" {
		return fmt.Errorf("administrator username is not set")
	}
	if config.RabbitMQConfig.Administrator.Password == "" {
		return fmt.Errorf("administrator password is not set")
	}

	return nil
}

func (p *PolicyDefinition) UnmarshalYAML(f func(interface{}) error) error {
	var s string
	f(&s)
	return json.Unmarshal([]byte(s), p)
}
