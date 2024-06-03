package activities

import (
	"context"
	"time"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/activity"
)

func ShipOrder(ctx context.Context, input resources.OrderInput, item resources.Item) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Ship Order activity started", "orderId", input.OrderId, "ItemId", item.Id, "Item Description", item.Description)

	time.Sleep(1 * time.Second)

	return nil
}
