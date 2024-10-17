package activities

import (
	"context"
	"temporal-order-management/app"

	"go.temporal.io/sdk/activity"
)

func ShipOrder(ctx context.Context, input app.ShippingInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Ship Order activity started", "orderId", input.Order.OrderId, "ItemId", input.Item.Id, "Item Description", input.Item.Description)

	// simulate external API call
	simulateExternalOperation(1000)

	return nil
}
