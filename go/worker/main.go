package main

import (
	"log"
	"os"
	"temporal-order-management/activities"
	"temporal-order-management/app"
	"temporal-order-management/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	c, err := client.Dial(app.GetClientOptions())
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
	w.RegisterWorkflowWithOptions(workflows.OrderWorkflowScenarios, workflow.RegisterOptions{
		Name: "OrderWorkflowNexusOperation",
	})
	w.RegisterWorkflow(workflows.ShippingWorkflow)
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
