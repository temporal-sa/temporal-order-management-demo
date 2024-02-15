package activities

import (
	"context"
	"time"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/activity"
)

// Basic activity definition

func CheckFraud(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Check Fraud activity started", "oderId", input.OrderId)

	time.Sleep(1 * time.Second)

	return input.OrderId, nil
}
