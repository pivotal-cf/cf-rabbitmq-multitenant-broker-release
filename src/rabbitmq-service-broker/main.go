package main

import (
	"fmt"
	"log"
	"net/http"

	"rabbitmq-service-broker/broker"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

const port = 8901

func main() {
	brokerLogger := lager.NewLogger("rabbitmq-multitenant-broker")
	brokerAPI := brokerapi.New(broker.RabbitMQServiceBroker{}, brokerLogger, brokerapi.BrokerCredentials{})
	http.Handle("/", brokerAPI)
	fmt.Printf("RabbitMQ Service Broker listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
