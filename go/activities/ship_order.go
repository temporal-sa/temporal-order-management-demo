package activities

import (
	"context"
	"math/rand"
	"temporal-order-management/app"
	"time"

	"go.temporal.io/sdk/activity"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ShipOrder(ctx context.Context, input app.ShippingInput) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Ship Order activity started", "orderId", input.Order.OrderId, "ItemId", input.Item.Id, "Item Description", input.Item.Description)

	// simulate external API call
	delayMs := rand.Intn(3001) + 1000
	logger.Info("Shipping Delay Time", "delayMs", delayMs)
	simulateExternalOperation(delayMs)

	return nil
}
