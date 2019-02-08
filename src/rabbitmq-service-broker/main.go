package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"rabbitmq-service-broker/broker"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
	yaml "gopkg.in/yaml.v2"
)

const port = 8901

var configPath string

func init() {
	flag.StringVar(&configPath, "configPath", "", "Config file location")
}

func main() {
	flag.Parse()

	brokerLogger := lager.NewLogger("rabbitmq-multitenant-broker")

	configBytes, err := ioutil.ReadFile(filepath.FromSlash(configPath))
	if err != nil {
		brokerLogger.Fatal("read-config", err)
	}

	config := broker.Config{}
	if err = yaml.Unmarshal(configBytes, &config); err != nil {
		brokerLogger.Fatal("config-unmarshal", err)
	}

	broker := broker.New(config)
	credentials := brokerapi.BrokerCredentials{
		Username: config.ServiceConfig.Username,
		Password: config.ServiceConfig.Password,
	}

	brokerAPI := brokerapi.New(broker, brokerLogger, credentials)
	http.Handle("/", brokerAPI)
	fmt.Printf("RabbitMQ Service Broker listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
