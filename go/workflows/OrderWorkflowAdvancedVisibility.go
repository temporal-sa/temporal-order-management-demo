package workflows

// import (
// 	"time"

// 	"temporal-order-management/activities"

// 	"github.com/google/uuid"

// 	"temporal-order-management/resources"

// 	"go.temporal.io/sdk/temporal"
// 	"go.temporal.io/sdk/workflow"
// )

// func OrderWorkflowAdvancedVisibility(ctx workflow.Context, input resources.OrderInput) (*resources.OrderOutput, error) {
// 	logger := workflow.GetLogger(ctx)
// 	logger.Info("Processing order started", "orderId", input.OrderId)

// 	// activity options
// 	activityOptions := workflow.ActivityOptions{
// 		StartToCloseTimeout: 5 * time.Second,
// 		RetryPolicy: &temporal.RetryPolicy{
// 			InitialInterval:    1 * time.Second,
// 			BackoffCoefficient: 2.0,
// 			MaximumInterval:    30 * time.Second,
// 		},
// 	}
// 	ctx = workflow.WithActivityOptions(ctx, activityOptions)

// 	// local activity options
// 	localActivityOptions := workflow.LocalActivityOptions{
// 		StartToCloseTimeout: 5 * time.Second,
// 	}
// 	laCtx := workflow.WithLocalActivityOptions(ctx, localActivityOptions)

// 	// Expose progress as query
// 	progress, err := resources.SetQueryHandlerForProgress(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get items
// 	items := resources.Items{}
// 	err = workflow.ExecuteLocalActivity(laCtx, activities.GetItems).Get(ctx, &items)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Upsert Check Fraud
// 	orderStatus := map[string]interface{}{
// 		"OrderStatus": "Check Fraud",
// 	}
// 	workflow.UpsertSearchAttributes(ctx, orderStatus)

// 	// Check Fraud
// 	var result1 string
// 	err = workflow.ExecuteActivity(ctx, activities.CheckFraud, input).Get(ctx, &result1)
// 	if err != nil {
// 		return nil, err
// 	}

// 	*progress = 25
// 	workflow.Sleep(ctx, 3*time.Second)

// 	// Upsert Prepare Shipment
// 	orderStatus = map[string]interface{}{
// 		"OrderStatus": "Prepare Shipment",
// 	}
// 	workflow.UpsertSearchAttributes(ctx, orderStatus)

// 	// Prepare Shipment
// 	var result2 string
// 	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, input).Get(ctx, &result2)
// 	if err != nil {
// 		return nil, err
// 	}

// 	*progress = 50
// 	workflow.Sleep(ctx, 3*time.Second)

// 	// Upsert Charge Customer
// 	orderStatus = map[string]interface{}{
// 		"OrderStatus": "Charge Customer",
// 	}
// 	workflow.UpsertSearchAttributes(ctx, orderStatus)

// 	// Charge Customer
// 	var result3 string
// 	err = workflow.ExecuteActivity(ctx, activities.ChargeCustomer, input).Get(ctx, &result3)
// 	if err != nil {
// 		return nil, err
// 	}

// 	*progress = 75
// 	workflow.Sleep(ctx, 3*time.Second)

// 	// Upsert Ship Order
// 	orderStatus = map[string]interface{}{
// 		"OrderStatus": "Ship Order",
// 	}
// 	workflow.UpsertSearchAttributes(ctx, orderStatus)

// 	// Ship Orders
// 	var shipItems []workflow.Future
// 	for _, item := range items {
// 		logger.Info("Shipping item " + item.Description)
// 		shipItem := workflow.ExecuteActivity(ctx, activities.ShipOrder, input, item)
// 		shipItems = append(shipItems, shipItem)
// 	}

// 	// Wait for all items to ship
// 	for _, shipItem := range shipItems {
// 		err = shipItem.Get(ctx, nil)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	*progress = 100

// 	// Upsert Order Completed
// 	orderStatus = map[string]interface{}{
// 		"OrderStatus": "Order Completed",
// 	}
// 	workflow.UpsertSearchAttributes(ctx, orderStatus)

// 	// Generate Tracking Id
// 	trackingId := uuid.New().String()

// 	output := &resources.OrderOutput{
// 		TrackingId: trackingId,
// 		Address:    input.Address,
// 	}

// 	return output, nil
// }
