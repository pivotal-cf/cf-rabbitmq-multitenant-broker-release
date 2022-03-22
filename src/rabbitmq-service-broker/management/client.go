package management

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"rabbitmq-service-broker/config"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

func NewClient(cfg config.Config) (*rabbithole.Client, error) {
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
