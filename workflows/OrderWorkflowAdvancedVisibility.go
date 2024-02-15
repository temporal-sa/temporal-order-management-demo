package workflows

import (
	"time"

	"github.com/ktenzer/temporal-order-management/activities"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func OrderWorkflowAdvancedVisibility(ctx workflow.Context, input resources.OrderInput) (*resources.OrderOutput, error) {
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

	// upsert check fraud
	orderStatus := map[string]interface{}{
		"OrderStatus": "Check Fraud",
	}
	workflow.UpsertSearchAttributes(ctx, orderStatus)

	workflow.Sleep(ctx, 3*time.Second)

	var result1 string
	err := workflow.ExecuteActivity(ctx, activities.CheckFraud, input).Get(ctx, &result1)
	if err != nil {
		return nil, err
	}

	// upsert prepare shipment
	orderStatus = map[string]interface{}{
		"OrderStatus": "Prepare Shipment",
	}
	workflow.UpsertSearchAttributes(ctx, orderStatus)

	workflow.Sleep(ctx, 3*time.Second)

	var result2 string
	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, input).Get(ctx, &result2)
	if err != nil {
		return nil, err
	}

	// upsert charge customer
	orderStatus = map[string]interface{}{
		"OrderStatus": "Charge Customer",
	}
	workflow.UpsertSearchAttributes(ctx, orderStatus)

	workflow.Sleep(ctx, 3*time.Second)

	var result3 string
	err = workflow.ExecuteActivity(ctx, activities.ChargeCustomer, input).Get(ctx, &result3)
	if err != nil {
		return nil, err
	}

	// upsert ship order
	orderStatus = map[string]interface{}{
		"OrderStatus": "Ship Order",
	}
	workflow.UpsertSearchAttributes(ctx, orderStatus)

	workflow.Sleep(ctx, 3*time.Second)

	var trackingId string
	err = workflow.ExecuteActivity(ctx, activities.ShipOrder, input).Get(ctx, &trackingId)
	if err != nil {
		return nil, err
	}

	// upsert complete
	orderStatus = map[string]interface{}{
		"OrderStatus": "Complete",
	}
	workflow.UpsertSearchAttributes(ctx, orderStatus)

	output := &resources.OrderOutput{
		TrackingId: trackingId,
		Address:    input.Address,
	}

	return output, nil
}
