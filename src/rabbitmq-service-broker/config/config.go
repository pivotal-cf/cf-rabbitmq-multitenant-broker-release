package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("yaml")
	})
}

type Config struct {
	Service  Service  `yaml:"service"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
}

type Service struct {
	UUID                string `yaml:"uuid" validate:"required"`
	Name                string `yaml:"name" validate:"required"`
	PlanUUID            string `yaml:"plan_uuid" validate:"required"`
	Username            string `yaml:"username" validate:"required"`
	Password            string `yaml:"password" validate:"required"`
	Description         string `yaml:"offering_description"`
	DisplayName         string `yaml:"display_name"`
	LongDescription     string `yaml:"long_description"`
	ProviderDisplayName string `yaml:"provider_display_name"`
	DocumentationURL    string `yaml:"documentation_url"`
	SupportURL          string `yaml:"support_url"`
	IconImage           string `yaml:"icon_image"`
	Shareable           bool   `yaml:"shareable"`
}

type RabbitMQ struct {
	Hosts             []string              `yaml:"hosts" validate:"required,min=1"`
	Administrator     AdminCredentials      `yaml:"administrator" validate:"required"`
	Management        ManagementCredentials `yaml:"management"`
	ManagementDomain  string                `yaml:"management_domain" validate:"required"`
	OperatorSetPolicy RabbitMQPolicy        `yaml:"operator_set_policy"`
	RegularUserTags   string                `yaml:"regular_user_tags"`
	TLS               TLSEnabled            `yaml:"ssl"`
}

type ManagementCredentials struct {
	Username string `yaml:"username"`
}

type AdminCredentials struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type RabbitMQPolicy struct {
	Enabled    bool             `yaml:"enabled"`
	Name       string           `yaml:"policy_name"`
	Priority   int              `yaml:"policy_priority"`
	Definition PolicyDefinition `yaml:"policy_definition"`
}

type PolicyDefinition map[string]interface{}

func (p *PolicyDefinition) UnmarshalYAML(f func(interface{}) error) error {
	var s string
	if err := f(&s); err != nil {
		return err
	}
	return json.Unmarshal([]byte(s), p)
}

type TLSEnabled bool

func (t *TLSEnabled) UnmarshalYAML(f func(interface{}) error) error {
	var s string
	if err := f(&s); err != nil {
		return err
	}

	*t = TLSEnabled(len(s) != 0)
	return nil
}

func Read(path string) (Config, error) {
	configBytes, err := ioutil.ReadFile(filepath.FromSlash(path))
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err = yaml.Unmarshal(configBytes, &config); err != nil {
		return Config{}, err
	}

	if err := validateConfig(config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func validateConfig(config Config) error {
	if err := validate.Struct(config); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			var missing []string
			for _, err := range errs {
				missing = append(missing, strings.TrimPrefix(err.Namespace(), "Config."))
			}
			return fmt.Errorf("Config file has missing fields: " + strings.Join(missing, ", "))
		}
		return err
	}

	return nil
}
