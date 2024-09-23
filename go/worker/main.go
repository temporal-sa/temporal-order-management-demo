package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"temporal-order-management/activities"
	"temporal-order-management/workflows"

	"go.temporal.io/sdk/client"
	tlog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	c, err := client.Dial(getClientOptions())
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, os.Getenv("TEMPORAL_TASK_QUEUE"), worker.Options{})

	// workflows
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflow, workflow.RegisterOptions{
		Name: "OrderWorkflowHappyPath",
	})
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflowScenarios, workflow.RegisterOptions{
		Name: "OrderWorkflowAPIFailure",
	})
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflowScenarios, workflow.RegisterOptions{
		Name: "OrderWorkflowRecoverableFailure",
	})
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflowScenarios, workflow.RegisterOptions{
		Name: "OrderWorkflowNonRecoverableFailure",
	})
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflowScenarios, workflow.RegisterOptions{
		Name: "OrderWorkflowChildWorkflow",
	})
	w.RegisterWorkflow(workflows.ShippingChildWorkflow)
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflowScenarios, workflow.RegisterOptions{
		Name: "OrderWorkflowAdvancedVisibility",
	})
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflowScenarios, workflow.RegisterOptions{
		Name: "OrderWorkflowHumanInLoopSignal",
	})
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflowScenarios, workflow.RegisterOptions{
		Name: "OrderWorkflowHumanInLoopUpdate",
	})

	// activities
	w.RegisterActivity(activities.GetItems)
	w.RegisterActivity(activities.CheckFraud)
	w.RegisterActivity(activities.PrepareShipment)
	w.RegisterActivity(activities.UndoPrepareShipment)
	w.RegisterActivity(activities.ChargeCustomer)
	w.RegisterActivity(activities.UndoChargeCustomer)
	w.RegisterActivity(activities.ShipOrder)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func getClientOptions() client.Options {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	address := getEnv("TEMPORAL_HOST_URL", "localhost:7233")
	namespace := getEnv("TEMPORAL_NAMESPACE", "default")
	clientOptions := client.Options{
		HostPort:  address,
		Namespace: namespace,
		Logger:    tlog.NewStructuredLogger(logger),
	}

	apiKey := getEnv("TEMPORAL_APIKEY", "")
	tlsCertPath := getEnv("TEMPORAL_MTLS_TLS_CERT", "")
	tlsKeyPath := getEnv("TEMPORAL_MTLS_TLS_KEY", "")

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
