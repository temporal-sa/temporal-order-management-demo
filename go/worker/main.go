package main

import (
	"crypto/tls"
	"log"
	"log/slog"
	"os"
	"temporal-order-management/activities"
	"temporal-order-management/workflows"

	"go.temporal.io/sdk/client"
	tlog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
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

	// w.RegisterWorkflow(workflows.OrderWorkflowAdvancedVisibility)
	// w.RegisterWorkflow(workflows.OrderWorkflowChildWorkflow)
	// w.RegisterWorkflow(workflows.OrderWorkflowHumanInLoopSignal)
	// w.RegisterWorkflow(workflows.OrderWorkflowHumanInLoopUpdate)
	// w.RegisterWorkflow(workflows.ShippingChildWorkflow)

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

	tlsCertPath := getEnv("TEMPORAL_MTLS_TLS_CERT", "")
	tlsKeyPath := getEnv("TEMPORAL_MTLS_TLS_KEY", "")
	if tlsCertPath != "" && tlsKeyPath != "" {
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
