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
	Hosts             Hosts                 `yaml:"hosts"`
	DNSHost           string                `yaml:"dns_host"`
	Administrator     AdminCredentials      `yaml:"administrator" validate:"required"`
	Management        ManagementCredentials `yaml:"management"`
	ManagementTLS     ManagementTLS         `yaml:"management_tls"`
	ManagementDomain  string                `yaml:"management_domain" validate:"required"`
	OperatorSetPolicy RabbitMQPolicy        `yaml:"operator_set_policy"`
	RegularUserTags   string                `yaml:"regular_user_tags"`
	TLS               TLSEnabled            `yaml:"ssl"`
}

type ManagementCredentials struct {
	Username string `yaml:"username"`
}

type ManagementTLS struct {
	Enabled    bool   `yaml:"enabled"`
	CACert     string `yaml:"cacert"`
	SkipVerify bool   `yaml:"skip_verify"`
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

// TLSEnabled is a boolean true/false according to the job spec, but we must also support usages
// where it's set to the the value of a certificate (for true) or empty/null (for false)
type TLSEnabled bool

func (t *TLSEnabled) UnmarshalYAML(f func(interface{}) error) error {
	var s string
	if err := f(&s); err != nil {
		return err
	}

	*t = TLSEnabled(len(s) != 0 && s != "false")
	return nil
}

// Hosts can be provided in the config file as either a YAML list, or a comma-seperated list
type Hosts []string

func (h *Hosts) UnmarshalYAML(f func(interface{}) error) error {
	if err := h.unmarshalAsList(f); err == nil {
		return nil
	}

	if err := h.unmarshalAsString(f); err != nil {
		return err
	}

	return nil
}

func (h *Hosts) unmarshalAsList(f func(interface{}) error) error {
	var slice []string

	if err := f(&slice); err != nil {
		return err
	}

	*h = Hosts(slice)
	return nil
}

func (h *Hosts) unmarshalAsString(f func(interface{}) error) error {
	var s string
	if err := f(&s); err != nil {
		return err
	}

	*h = Hosts(splitCommaSeparatedList(s))
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

func (c *Config) NodeHosts() []string {
	if host := c.RabbitMQ.DNSHost; host != "" {
		return []string{host}
	}
	return c.RabbitMQ.Hosts
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

	if config.RabbitMQ.OperatorSetPolicy.Enabled && config.RabbitMQ.OperatorSetPolicy.Definition == nil {
		return fmt.Errorf("Config file has missing field: operator_set_policy.policy_definition must be provided when operator_set_policy.enabled is true")
	}

	if nodeHosts := config.NodeHosts(); len(nodeHosts) == 0 {
		return fmt.Errorf("Config file has missing fields: at least one of rabbitmq.hosts or rabbitmq.dns_host must be specified")
	}

	return nil
}

func splitCommaSeparatedList(s string) []string {
	list := strings.Split(s, ",")
	for i := range list {
		list[i] = strings.TrimSpace(list[i])
	}

	return list
}
