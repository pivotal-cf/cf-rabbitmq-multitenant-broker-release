package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"rabbitmq-service-broker/broker"
	"rabbitmq-service-broker/config"

	"code.cloudfoundry.org/lager"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-cf/brokerapi"
)

var (
	configPath string
	port       int
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "Config file location")
	flag.IntVar(&port, "port", 8901, "Port to listen on")
}

func main() {
	flag.Parse()

	logger := lager.NewLogger("rabbitmq-service-broker")
	logger.RegisterSink(lager.NewPrettySink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewPrettySink(os.Stderr, lager.ERROR))

	cfg, err := config.Read(configPath)
	if err != nil {
		logger.Fatal("read-config", err)
	}

	client, _ := rabbithole.NewClient(
		fmt.Sprintf("http://%s:15672", cfg.RabbitMQ.Hosts[0]),
		cfg.RabbitMQ.Administrator.Username,
		cfg.RabbitMQ.Administrator.Password,
	)

	broker := broker.New(cfg, client, logger)
	credentials := brokerapi.BrokerCredentials{
		Username: cfg.Service.Username,
		Password: cfg.Service.Password,
	}

	brokerAPI := brokerapi.New(broker, logger, credentials)
	http.Handle("/", brokerAPI)
	logger.Info("main-serving", lager.Data{"port": port})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
