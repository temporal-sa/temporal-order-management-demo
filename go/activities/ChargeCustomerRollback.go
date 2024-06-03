package activities

import (
	"context"
	"time"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/activity"
)

func ChargeCustomerRollback(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Rollback Charge Customer activity started", "orderId", input.OrderId)

	time.Sleep(1 * time.Second)

	return input.OrderId, nil
}
