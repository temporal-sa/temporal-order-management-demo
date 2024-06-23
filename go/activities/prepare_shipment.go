package activities

import (
	"context"
	"temporal-order-management/resources"

	"go.temporal.io/sdk/activity"
)

func PrepareShipment(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Prepare Shipment activity started", "orderId", input.OrderId)

	// simulate external API call
	simulateExternalOperation(1000)

	return input.OrderId, nil
}

func UndoPrepareShipment(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Undo Prepare Shipment activity started", "orderId", input.OrderId)

	// simulate external API call
	simulateExternalOperation(1000)

	return input.OrderId, nil
}
