package activities

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/activity"
)

func ShipOrder(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Ship Order activity started", "oderId", input.OrderId)

	time.Sleep(1 * time.Second)

	result := uuid.New().String()
	return result, nil
}
