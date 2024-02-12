package activities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"go.temporal.io/sdk/activity"
)

// Basic activity definition

func CheckFraud(ctx context.Context, input OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Check Fraud activity started", "oderId", input)

	time.Sleep(1 * time.Second)

	return input, nil
}

func PrepareShipment(ctx context.Context, input OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Prepare SHipment activity started", "orderId", input)

	time.Sleep(1 * time.Second)

	return input, nil
}

func ChargeCustomer(ctx context.Context, input OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Charge Customer activity started", "orderId", input)

	activityInfo := activity.GetInfo(ctx)
	activityAttempt := activityInfo.Attempt

	if activityAttempt >= 0 && activityAttempt <= 3 {
		time.Sleep(1 * time.Second)
		return input, errors.New("Charge Customer Failed: payment service timeout")
	}

	return input, nil
}

func ShipOrder(ctx context.Context, input OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Ship Order activity started", "oderId", input)

	time.Sleep(1 * time.Second)

	result := uuid.New().String()
	return result, nil
}
