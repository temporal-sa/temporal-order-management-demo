package app

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
	sdktally "go.temporal.io/sdk/contrib/tally"
	tlog "go.temporal.io/sdk/log"
)

func GetClientOptions() client.Options {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	clientOptions, err := envconfig.LoadDefaultClientOptions()
	if err != nil {
		log.Fatalln("error loading default client options", err)
	}

	clientOptions.Logger = tlog.NewStructuredLogger(logger)
	clientOptions.MetricsHandler = sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
		ListenAddress: "0.0.0.0:9090",
		TimerType:     "histogram",
	}))

	apiKey := GetEnv("TEMPORAL_API_KEY", "")
	if apiKey != "" {
		// "kms" service
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.Method != http.MethodPut {
				http.Error(w, "", http.StatusMethodNotAllowed)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			apiKey = string(body)
			log.Default().Println("API key updated")

			w.WriteHeader(http.StatusAccepted)
		})

		go func() {
			err := http.ListenAndServe(":3333", nil) // make this an env variable
			if err != nil {
				log.Fatalln("Unable to start webserver", err)
			}
		}()

		clientOptions.Credentials = client.NewAPIKeyDynamicCredentials(
			func(context.Context) (string, error) {
				return apiKey, nil
			},
		)
	}
	return clientOptions
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
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
