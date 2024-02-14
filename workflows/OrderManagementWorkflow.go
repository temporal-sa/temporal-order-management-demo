package workflows

import (
	"errors"
	"time"

	"github.com/ktenzer/temporal-order-management/activities"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func OrderManagementWorkflow(ctx workflow.Context, input resources.OrderInput) (*resources.OrderOutput, error) {
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

	var result1 string
	err := workflow.ExecuteActivity(ctx, activities.CheckFraud, input).Get(ctx, &result1)
	if err != nil {
		return nil, err
	}

	logger.Info("Sleeping for 1 second...")
	workflow.Sleep(ctx, 1*time.Second)

	var result2 string
	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, input).Get(ctx, &result2)
	if err != nil {
		return nil, err
	}

	var result3 string
	err = workflow.ExecuteActivity(ctx, activities.ChargeCustomer, input).Get(ctx, &result3)
	if err != nil {
		return nil, err
	}

	if input.Scenario == "UNRECOVERABLE_FAILURE" {
		//Divide by zero exception
		//Divide(1, 0)
	}

	// Start timer and wait for timer to fire or signal
	if input.Scenario == "HUMAN_IN_THE_LOOP_SIGNAL" {
		address, isCancelled := resources.SignalApprovalTimer(ctx)
		if isCancelled {
			return nil, errors.New("Time limit for approval has been exceeded!")
		}

		input.Address = address
	}

	// Start timer and wait for timer to fire or update
	if input.Scenario == "HUMAN_IN_THE_LOOP_UPDATE" {
		address, isCancelled := resources.UpdateApprovalTimer(ctx)
		if isCancelled {
			return nil, errors.New("Time limit for approval has been exceeded!")
		}

		input.Address = address
	}

	var trackingId string
	err = workflow.ExecuteActivity(ctx, activities.ShipOrder, input).Get(ctx, &trackingId)
	if err != nil {
		return nil, err
	}

	output := &resources.OrderOutput{
		TrackingId: trackingId,
		Address:    input.Address,
	}
	return output, nil
}

func Divide(a, b int) (int, error) {
	return a / b, nil
}

func boolPointer(b bool) *bool {
	return &b
}
