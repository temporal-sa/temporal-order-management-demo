package activities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

// Basic activity definition

func CheckFraud(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Check Fraud activity started", "oderId", input.OrderId)

	time.Sleep(1 * time.Second)

	return input.OrderId, nil
}

func PrepareShipment(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Prepare Shipment activity started", "orderId", input.OrderId)

	time.Sleep(1 * time.Second)

	return input.OrderId, nil
}

func ChargeCustomer(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Charge Customer activity started", "orderId", input.OrderId)

	activityInfo := activity.GetInfo(ctx)
	activityAttempt := activityInfo.Attempt

	if input.Scenario == "API_DOWNTIME" {
		if activityAttempt >= 0 && activityAttempt <= 3 {
			time.Sleep(1 * time.Second)
			return input.OrderId, errors.New("Charge Customer Failed: payment service timeout")
		}
	}

	if input.Scenario == "UNRECOVERABLE_FAILURE" {
		return input.OrderId, temporal.NewNonRetryableApplicationError("Could not process payment", "activityFailure", errors.New("Credit card invalid!"))
	}

	return input.OrderId, nil
}

func ShipOrder(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Ship Order activity started", "oderId", input.OrderId)

	time.Sleep(1 * time.Second)

	result := uuid.New().String()
	return result, nil
}
