package broker

import (
	"fmt"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	ServiceConfig ServiceConfig `yaml:"service"`
}
type ServiceConfig struct {
	UUID                string `yaml:"uuid"`
	Name                string `yaml:"name"`
	OfferingDescription string `yaml:"offering_description"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	DisplayName         string `yaml:"display_name"`
	IconImage           string `yaml:"icon_image"`
	LongDescription     string `yaml:"long_description"`
	ProviderDisplayName string `yaml:"provider_display_name"`
	DocumentationUrl    string `yaml:"documentation_url"`
	SupportUrl          string `yaml:"support_url"`
}

func ParseConfig(reader io.Reader) (*Config, error) {
	configBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	if err := ValidateConfig(&config); err != nil {
		return nil, err
	}

	//if err := json.Unmarshal([]byte(config.Rabbitmq.Policy.EncodedDefinition), &config.Rabbitmq.Policy.Definition); err != nil {
	//	return nil, err
	//}

	return &config, nil
}

func ValidateConfig(config *Config) error {
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

	return nil
}
