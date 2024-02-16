package workflows

import (
	"time"

	"github.com/ktenzer/temporal-order-management/activities"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ShippingWorkflow(ctx workflow.Context, input resources.OrderInput, item resources.Item) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Processing shipping started", "orderId", input.OrderId)

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
