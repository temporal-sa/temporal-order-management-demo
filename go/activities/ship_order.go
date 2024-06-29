package activities

import (
	"context"
	"temporal-order-management/app"

	"go.temporal.io/sdk/activity"
)

func ShipOrder(ctx context.Context, input app.OrderInput, item app.Item) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Ship Order activity started", "orderId", input.OrderId, "ItemId", item.Id, "Item Description", item.Description)

	// simulate external API call
	simulateExternalOperation(1000)

	return nil
}
