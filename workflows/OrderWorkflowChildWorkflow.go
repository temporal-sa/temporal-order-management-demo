package workflows

import (
	"time"

	"github.com/ktenzer/temporal-order-management/activities"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func OrderWorkflowChildWorkflow(ctx workflow.Context, input resources.OrderInput) (*resources.OrderOutput, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Processing order started", "orderId", input.OrderId)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	workflow.Sleep(ctx, 3*time.Second)

	var result1 string
	err := workflow.ExecuteActivity(ctx, activities.CheckFraud, input).Get(ctx, &result1)
	if err != nil {
		return nil, err
	}

	workflow.Sleep(ctx, 3*time.Second)

	var result2 string
	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, input).Get(ctx, &result2)
	if err != nil {
		return nil, err
	}

	workflow.Sleep(ctx, 3*time.Second)

	var result3 string
	err = workflow.ExecuteActivity(ctx, activities.ChargeCustomer, input).Get(ctx, &result3)
	if err != nil {
		return nil, err
	}

	workflow.Sleep(ctx, 3*time.Second)

	output := &resources.OrderOutput{}

	// set child workflow options
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:        "shipment-" + input.OrderId,
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_TERMINATE,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	// execute and wait on child workflow
	err = workflow.ExecuteChildWorkflow(ctx, "ShippingWorkflow", input).Get(ctx, &output)
	if err != nil {
		return nil, err
	}

	return output, nil
}
