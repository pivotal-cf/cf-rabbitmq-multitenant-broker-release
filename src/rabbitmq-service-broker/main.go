package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"rabbitmq-service-broker/broker"

	"code.cloudfoundry.org/lager"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-cf/brokerapi"
)

const port = 8901

var configPath string

func init() {
	flag.StringVar(&configPath, "configPath", "", "Config file location")
}

func main() {
	flag.Parse()

	logger := lager.NewLogger("rabbitmq-service-broker")
	logger.RegisterSink(lager.NewPrettySink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewPrettySink(os.Stderr, lager.ERROR))

	config, err := broker.ReadConfig(configPath)
	if err != nil {
		logger.Fatal("read-config", err)
	}

	client, _ := rabbithole.NewClient(
		fmt.Sprintf("http://%s:15672", config.RabbitMQConfig.Hosts[0]),
		config.RabbitMQConfig.Administrator.Username,
		config.RabbitMQConfig.Administrator.Password,
	)

	broker := broker.New(config, client, logger)
	credentials := brokerapi.BrokerCredentials{
		Username: config.ServiceConfig.Username,
		Password: config.ServiceConfig.Password,
	}

	brokerAPI := brokerapi.New(broker, logger, credentials)
	http.Handle("/", brokerAPI)
	fmt.Printf("RabbitMQ Service Broker listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
