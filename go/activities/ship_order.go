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

	// item n ships in n*2 seconds, so a fan-out finishes one item at a time; no index means random
	delayMs := input.Index * 2000
	if delayMs == 0 {
		delayMs = rand.Intn(3001) + 1000
	}

	// on a retry, pick up from the last second this item reported
	secondsDone := 0
	if activity.HasHeartbeatDetails(ctx) {
		if err := activity.GetHeartbeatDetails(ctx, &secondsDone); err != nil {
			return err
		}
	}
	logger.Info("Shipping Delay Time", "delayMs", delayMs, "resumingAtSeconds", secondsDone)

	// simulate external API call, heartbeating progress each second so a retry can resume
	for second := secondsDone; second < delayMs/1000; second++ {
		simulateExternalOperation(1000)
		activity.RecordHeartbeat(ctx, second+1)
	}
	simulateExternalOperation(delayMs % 1000)

	return nil
}
