package main

import (
	"log"
	"temporal-order-management/activities"
	"temporal-order-management/app"
	"temporal-order-management/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	co := app.GetClientOptions()
	c, err := client.Dial(co)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	log.Printf("✅ Client connected to %v in namespace '%v'", co.HostPort, co.Namespace)
	defer c.Close()

	w := worker.New(c, app.GetEnv("TEMPORAL_TASK_QUEUE", "orders"), worker.Options{})

	// workflows
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflow, workflow.RegisterOptions{
		Name: "OrderWorkflowHappyPath",
	})
	w.RegisterDynamicWorkflow(workflows.OrderWorkflowScenarios, workflow.DynamicRegisterOptions{})
	w.RegisterWorkflow(workflows.ShippingWorkflow)

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
