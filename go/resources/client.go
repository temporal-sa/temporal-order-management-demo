package resources

import (
	"log"
	"os"
	"time"

	"crypto/tls"
	"crypto/x509"

	"go.temporal.io/sdk/client"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	sdktally "go.temporal.io/sdk/contrib/tally"
)

func GetClientOptions(clientType string) client.Options {

	var clientOptions client.Options
	if clientType == "worker" {
		clientOptions = client.Options{
			HostPort:  os.Getenv("TEMPORAL_HOST_URL"),
			Namespace: os.Getenv("TEMPORAL_NAMESPACE"),

			MetricsHandler: sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
				ListenAddress: "0.0.0.0:9090",
				TimerType:     "histogram",
			})),
		}
	} else {
		clientOptions = client.Options{
			HostPort:  os.Getenv("TEMPORAL_HOST_URL"),
			Namespace: os.Getenv("TEMPORAL_NAMESPACE"),
		}
	}

	if os.Getenv("TEMPORAL_MTLS_TLS_CERT") != "" && os.Getenv("TEMPORAL_MTLS_TLS_KEY") != "" {
		if os.Getenv("TEMPORAL_MTLS_TLS_CA") != "" {
			caCert, err := os.ReadFile(os.Getenv("TEMPORAL_MTLS_TLS_CA"))
			if err != nil {
				log.Fatalln("failed reading server CA's certificate", err)
			}

			certPool := x509.NewCertPool()
			if !certPool.AppendCertsFromPEM(caCert) {
				log.Fatalln("failed to add server CA's certificate", err)
			}

			cert, err := tls.LoadX509KeyPair(os.Getenv("TEMPORAL_MTLS_TLS_CERT"), os.Getenv("TEMPORAL_MTLS_TLS_KEY"))
			if err != nil {
				log.Fatalln("Unable to load certs", err)
			}

			var serverName string
			if os.Getenv("TEMPORAL_MTLS_TLS_ENABLE_HOST_VERIFICATION") == "true" {
				serverName = os.Getenv("TEMPORAL_MTLS_TLS_SERVER_NAME")
			}

			clientOptions.ConnectionOptions = client.ConnectionOptions{
				TLS: &tls.Config{
					RootCAs:      certPool,
					Certificates: []tls.Certificate{cert},
					ServerName:   serverName,
				},
			}
		} else {
			cert, err := tls.LoadX509KeyPair(os.Getenv("TEMPORAL_MTLS_TLS_CERT"), os.Getenv("TEMPORAL_MTLS_TLS_KEY"))
			if err != nil {
				log.Fatalln("Unable to load certs", err)
			}

			clientOptions.ConnectionOptions = client.ConnectionOptions{
				TLS: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			}
		}
	}

	return clientOptions
}

func newPrometheusScope(c prometheus.Configuration) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		log.Fatalln("error creating prometheus reporter", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
		Prefix:          "temporal_samples",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	log.Println("prometheus metrics scope created")
	return scope
}
