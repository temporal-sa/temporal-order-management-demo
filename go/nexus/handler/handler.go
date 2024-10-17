package handler

import (
	"context"
	"fmt"
	"temporal-order-management/app"
	"temporal-order-management/workflows"

	"go.temporal.io/sdk/client"

	"github.com/nexus-rpc/sdk-go/nexus"
	"go.temporal.io/sdk/temporalnexus"
)

var ShippingOperation = temporalnexus.NewWorkflowRunOperation(
	app.ShippingOperationName,
	workflows.ShippingWorkflow,
	func(ctx context.Context, input app.ShippingInput, soo nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
		return client.StartWorkflowOptions{ID: fmt.Sprintf("shipment-%v-%v", input.Order.OrderId, input.Item.Id)}, nil
	},
)
