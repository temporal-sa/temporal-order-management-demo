package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"temporal-order-management/activities"
	"temporal-order-management/app"
	"temporal-order-management/nexus/handler"
	"temporal-order-management/workflows"

	"github.com/nexus-rpc/sdk-go/nexus"
)

func main() {
	c, err := client.Dial(app.GetClientOptions())
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, app.GetEnv("TEMPORAL_NEXUS_TASK_QUEUE", "shipping"), worker.Options{})
	service := nexus.NewService(app.ShippingServiceName)
	err = service.Register(handler.ShippingOperation)
	if err != nil {
		log.Fatalln("Unable to register operations", err)
	}
	w.RegisterNexusService(service)
	w.RegisterWorkflow(workflows.ShippingWorkflow)
	w.RegisterActivity(activities.ShipOrder)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
