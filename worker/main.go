package main

import (
	"crypto/tls"
	"log"
	"log/slog"
	"os"

	"simple-go/activities"
	"simple-go/workflows"

	"go.temporal.io/sdk/client"
	tlog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

const TASK_QUEUE = "simple-task-queue"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	address := getEnv("TEMPORAL_ADDRESS", "helloworld.sdvdw.tmprl.cloud:7233")
	namespace := getEnv("TEMPORAL_NAMESPACE", "helloworld.sdvdw")
	clientOptions := client.Options{
		HostPort:  address,
		Namespace: namespace,
		Logger:    tlog.NewStructuredLogger(logger),
	}

	tlsCertPath := getEnv("TEMPORAL_TLS_CERT", "/home/ktenzer/temporal/certs/ca.pem")
	tlsKeyPath := getEnv("TEMPORAL_TLS_KEY", "/home/ktenzer/temporal/certs/ca.key")
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

	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	w := worker.New(temporalClient, TASK_QUEUE, worker.Options{})

	// w.RegisterWorkflow(workflows.Simple)

	w.RegisterWorkflow(workflows.OrderManagementWorkflow)
	w.RegisterActivity(activities.ChargeCustomer)
	w.RegisterActivity(activities.CheckFraud)
	w.RegisterActivity(activities.PrepareShipment)
	w.RegisterActivity(activities.ShipOrder)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
