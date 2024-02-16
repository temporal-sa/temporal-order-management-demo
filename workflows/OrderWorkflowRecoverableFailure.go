package workflows

import (
	"time"

	"github.com/google/uuid"
	"github.com/ktenzer/temporal-order-management/activities"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func OrderWorkflowRecoverableFailure(ctx workflow.Context, input resources.OrderInput) (*resources.OrderOutput, error) {
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

	// Side effect to generate trackingId
	generateTrackingId := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return uuid.New().String()
	})

	var trackingId string
	generateTrackingId.Get(&trackingId)

	// Expose items as query
	items, err := resources.QueryItems(ctx)
	if err != nil {
		return nil, err
	}

	// Update items
	err = workflow.ExecuteActivity(ctx, activities.GetItems).Get(ctx, &items)
	if err != nil {
		return nil, err
	}

	// CHeck Fraud
	var result1 string
	err = workflow.ExecuteActivity(ctx, activities.CheckFraud, input).Get(ctx, &result1)
	if err != nil {
		return nil, err
	}

	workflow.Sleep(ctx, 3*time.Second)

	// Prepare Shipment
	var result2 string
	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, input).Get(ctx, &result2)
	if err != nil {
		return nil, err
	}

	workflow.Sleep(ctx, 3*time.Second)

	// Charge Customer
	var result3 string
	err = workflow.ExecuteActivity(ctx, activities.ChargeCustomer, input).Get(ctx, &result3)
	if err != nil {
		return nil, err
	}

	//Divide by zero, produce recoverable exception
	//Divide(1, 0)

	workflow.Sleep(ctx, 3*time.Second)

	// Ship Order
	for _, item := range *items {
		err = workflow.ExecuteActivity(ctx, activities.ShipOrder, input, item).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	output := &resources.OrderOutput{
		TrackingId: trackingId,
		Address:    input.Address,
	}

	return output, nil
}
