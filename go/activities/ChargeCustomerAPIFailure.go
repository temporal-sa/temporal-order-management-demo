package activities

import (
	"context"
	"errors"
	"time"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/activity"
)

func ChargeCustomerAPIFailure(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Charge Customer activity started", "orderId", input.OrderId)

	activityInfo := activity.GetInfo(ctx)
	activityAttempt := activityInfo.Attempt

	if activityAttempt >= 0 && activityAttempt <= 3 {
		time.Sleep(1 * time.Second)
		return input.OrderId, errors.New("Charge Customer Failed: payment service timeout")
	}

	return input.OrderId, nil
}
