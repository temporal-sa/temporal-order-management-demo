package app

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"go.temporal.io/sdk/client"
	tlog "go.temporal.io/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetClientOptions() client.Options {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	address := GetEnv("TEMPORAL_ADDRESS", "localhost:7233")
	namespace := GetEnv("TEMPORAL_NAMESPACE", "default")
	clientOptions := client.Options{
		HostPort:  address,
		Namespace: namespace,
		Logger:    tlog.NewStructuredLogger(logger),
	}

	apiKey := GetEnv("TEMPORAL_APIKEY", "")
	tlsCertPath := GetEnv("TEMPORAL_CERT_PATH", "")
	tlsKeyPath := GetEnv("TEMPORAL_KEY_PATH", "")

	switch {
	case apiKey != "":
		serverName := strings.Split(address, ":")[0]

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
		clientOptions.ConnectionOptions = client.ConnectionOptions{
			TLS: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         serverName,
			},
			DialOptions: []grpc.DialOption{
				grpc.WithUnaryInterceptor(
					func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn,
						invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
						return invoker(
							metadata.AppendToOutgoingContext(ctx, "temporal-namespace", namespace),
							method,
							req,
							reply,
							cc,
							opts...,
						)
					},
				),
			},
		}
	case tlsCertPath != "" && tlsKeyPath != "":
		cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsKeyPath)
		if err != nil {
			log.Fatalln("Unable to load cert and key pair", err)
		}
		clientOptions.ConnectionOptions = client.ConnectionOptions{
			TLS: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
	}

	return clientOptions
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
