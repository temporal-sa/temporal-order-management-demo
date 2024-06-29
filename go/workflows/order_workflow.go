package workflows

import (
	"temporal-order-management/activities"
	"temporal-order-management/app"
	"temporal-order-management/messages"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func OrderWorkflow(ctx workflow.Context, input app.OrderInput) (*app.OrderOutput, error) {
	name := workflow.GetInfo(ctx).WorkflowType.Name
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

	localActivityOptions := workflow.LocalActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	laCtx := workflow.WithLocalActivityOptions(ctx, localActivityOptions)

	// Expose progress as query
	progress, err := messages.SetQueryHandlerForProgress(ctx)
	if err != nil {
		return nil, err
	}

	// Get items
	items := app.Items{}
	err = workflow.ExecuteLocalActivity(laCtx, activities.GetItems).Get(ctx, &items)
	if err != nil {
		return nil, err
	}

	// Check fraud
	var cfresult string
	err = workflow.ExecuteActivity(ctx, activities.CheckFraud, input).Get(ctx, &cfresult)
	if err != nil {
		return nil, err
	}

	updateProgress(progress, 25, ctx, 1)

	// Prepare shipment
	var psresult string
	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, input).Get(ctx, &psresult)
	if err != nil {
		return nil, err
	}

	updateProgress(progress, 50, ctx, 1)

	// Charge customer
	var ccresult string
	err = workflow.ExecuteActivity(ctx, activities.ChargeCustomer, input, name).Get(ctx, &ccresult)
	if err != nil {
		return nil, err
	}

	updateProgress(progress, 75, ctx, 3)

	// Ship orders
	var futures []workflow.Future
	for _, item := range items {
		logger.Info("Shipping item " + item.Description)
		f := workflow.ExecuteActivity(ctx, activities.ShipOrder, input, item)
		futures = append(futures, f)
	}

	// Wait for all items to ship
	for _, future := range futures {
		err = future.Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	updateProgress(progress, 100, ctx, 1)

	// Generate trackingId
	trackingId := uuid.New().String()
	output := &app.OrderOutput{
		TrackingId: trackingId,
		Address:    input.Address,
	}

	return output, nil
}

func updateProgress(progress *int, value int, ctx workflow.Context, seconds int) {
	*progress = value
	if seconds > 0 {
		duration := time.Duration(seconds) * time.Second
		workflow.Sleep(ctx, duration)
	}
}
