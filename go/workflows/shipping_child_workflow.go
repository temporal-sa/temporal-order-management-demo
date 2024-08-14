package workflows

import (
	"temporal-order-management/activities"
	"temporal-order-management/app"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ShippingChildWorkflow(ctx workflow.Context, input app.OrderInput, item app.Item) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Shipping workflow started", "orderId", input.OrderId)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	err := workflow.ExecuteActivity(ctx, activities.ShipOrder, input, item).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
