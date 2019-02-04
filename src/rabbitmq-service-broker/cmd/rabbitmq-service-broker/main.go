package main

import (
	"log"
	"net/http"

	"rabbitmq-service-broker/broker"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func main() {
	brokerLogger := lager.NewLogger("rabbitmq-multitenant-broker")
	brokerCredentials := brokerapi.BrokerCredentials{}
	serviceBroker := &broker.RabbitmqServiceBroker{}

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)

	http.Handle("/", brokerAPI)
	log.Fatal(http.ListenAndServe(":8901", nil))
}
