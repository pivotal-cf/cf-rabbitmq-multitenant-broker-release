package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"rabbitmq-service-broker/broker"
	"rabbitmq-service-broker/config"
	"rabbitmq-service-broker/rabbithutch"

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
	flag.IntVar(&port, "port", 4567, "Port to listen on")
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

	client, err := newManagementClient(cfg)
	if err != nil {
		logger.Fatal("create-rmq-management-client", err)
	}

	broker := broker.New(cfg, rabbithutch.New(client), logger)
	credentials := brokerapi.BrokerCredentials{
		Username: cfg.Service.Username,
		Password: cfg.Service.Password,
	}

	brokerAPI := brokerapi.New(broker, logger, credentials)
	http.Handle("/", brokerAPI)
	logger.Info("main-serving", lager.Data{"port": port})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func newManagementClient(cfg config.Config) (*rabbithole.Client, error) {
	if !cfg.RabbitMQ.ManagementTLS.Enabled {
		return rabbithole.NewClient(
			fmt.Sprintf("http://%s:15672", cfg.NodeHosts()[0]),
			cfg.RabbitMQ.Administrator.Username,
			cfg.RabbitMQ.Administrator.Password,
		)
	}

	//If we're here, configure a TLS client
	var caPool *x509.CertPool

	if cfg.RabbitMQ.ManagementTLS.CACert != "" {
		caPool = x509.NewCertPool()
		caPool.AppendCertsFromPEM([]byte(cfg.RabbitMQ.ManagementTLS.CACert))
	} else {
		caPool, _ = x509.SystemCertPool()
	}

	return rabbithole.NewTLSClient(
		fmt.Sprintf("https://%s:15671", cfg.NodeHosts()[0]),
		cfg.RabbitMQ.Administrator.Username,
		cfg.RabbitMQ.Administrator.Password,
		&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.RabbitMQ.ManagementTLS.SkipVerify,
				RootCAs:            caPool,
			},
		},
	)
}
