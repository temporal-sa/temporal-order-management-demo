package activities

import (
	"context"
	"errors"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

func ChargeCustomerUnrecoverableFailure(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Charge Customer activity started", "orderId", input.OrderId)

	return input.OrderId, temporal.NewNonRetryableApplicationError("Could not process payment", "activityFailure", errors.New("Credit card invalid!"))
}
