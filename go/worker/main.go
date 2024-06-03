package main

import (
	"log"
	"os"

	"github.com/ktenzer/temporal-order-management/activities"
	"github.com/ktenzer/temporal-order-management/resources"
	"github.com/ktenzer/temporal-order-management/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(resources.GetClientOptions("worker"))
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, os.Getenv("TEMPORAL_TASK_QUEUE"), worker.Options{})

	// workflows
	w.RegisterWorkflow(workflows.OrderWorkflowHappyPath)
	w.RegisterWorkflow(workflows.OrderWorkflowAPIFailure)
	w.RegisterWorkflow(workflows.OrderWorkflowAdvancedVisibility)
	w.RegisterWorkflow(workflows.OrderWorkflowChildWorkflow)
	w.RegisterWorkflow(workflows.OrderWorkflowHumanInLoopSignal)
	w.RegisterWorkflow(workflows.OrderWorkflowHumanInLoopUpdate)
	w.RegisterWorkflow(workflows.OrderWorkflowRecoverableFailure)
	w.RegisterWorkflow(workflows.OrderWorkflowNonRecoverableFailure)
	w.RegisterWorkflow(workflows.ShippingChildWorkflow)

	// activities
	w.RegisterActivity(activities.ChargeCustomerRollback)
	w.RegisterActivity(activities.ChargeCustomerAPIFailure)
	w.RegisterActivity(activities.ChargeCustomerNonRecoverableFailure)
	w.RegisterActivity(activities.ChargeCustomer)
	w.RegisterActivity(activities.CheckFraud)
	w.RegisterActivity(activities.PrepareShipment)
	w.RegisterActivity(activities.ShipOrder)
	w.RegisterActivity(activities.GetItems)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
