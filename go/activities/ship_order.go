package activities

import (
	"context"
	"temporal-order-management/resources"

	"go.temporal.io/sdk/activity"
)

func ShipOrder(ctx context.Context, input resources.OrderInput, item resources.Item) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Ship Order activity started", "orderId", input.OrderId, "ItemId", item.Id, "Item Description", item.Description)

	// simulate external API call
	simulateExternalOperation(1000)

	return nil
}
