package workflows

import (
	"fmt"
	"temporal-order-management/activities"
	"temporal-order-management/app"
	"temporal-order-management/messages"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	BUG        = "OrderWorkflowRecoverableFailure"
	CHILD      = "OrderWorkflowChildWorkflow"
	NEXUS      = "OrderWorkflowNexusOperation"
	SIGNAL     = "OrderWorkflowHumanInLoopSignal"
	UPDATE     = "OrderWorkflowHumanInLoopUpdate"
	VISIBILITY = "OrderWorkflowAdvancedVisibility"
)

var orderStatusKey = temporal.NewSearchAttributeKeyKeyword("OrderStatus")

func OrderWorkflowScenarios(ctx workflow.Context, args converter.EncodedValues) (output *app.OrderOutput, err error) {
	var input app.OrderInput
	err = args.Get(&input)
	if err != nil {
		return nil, fmt.Errorf("failed to decode arguments: %w", err)
	}

	name := workflow.GetInfo(ctx).WorkflowType.Name
	logger := workflow.GetLogger(ctx)
	logger.Info("Dynamic Order workflow started", "type", name, "orderId", input.OrderId)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	localActivityOptions := workflow.LocalActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	laCtx := workflow.WithLocalActivityOptions(ctx, localActivityOptions)

	// Create saga to manage order compensations
	var saga app.Saga
	defer func() {
		if err != nil {
			disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)
			saga.Compensate(disconnectedCtx)
		}
	}()

	// Expose progress as query
	progress, err := messages.SetQueryHandlerForProgress(ctx)
	if err != nil {
		return nil, err
	}

	// Get items
	items := app.Items{}
	err = workflow.ExecuteLocalActivity(laCtx, activities.GetItems).Get(ctx, &items)
	if err != nil {
		return nil, err
	}

	updateProgress("Check Fraud", progress, 0, ctx, 0)

	// Check fraud
	err = workflow.ExecuteActivity(ctx, activities.CheckFraud, input).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	updateProgress("Prepare Shipment", progress, 25, ctx, 1)

	// Prepare shipment
	saga.AddCompensation(activities.UndoPrepareShipment, input)
	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, input).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	updateProgress("Charge Customer", progress, 50, ctx, 1)

	// Charge customer
	saga.AddCompensation(activities.UndoChargeCustomer, input)
	err = workflow.ExecuteActivity(ctx, activities.ChargeCustomer, input, name).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	updateProgress("Ship Order", progress, 75, ctx, 3)

	if BUG == name {
		// Simulate bug
		panic("Simulated bug - fix me!")
	}

	if SIGNAL == name {
		// Await signal message to update address
		logger.Info("Waiting up to 60 seconds for updated address")
		var updateInput messages.UpdateOrderInput
		c := messages.GetSignalChannelForUpdateOrder(ctx)
		ok, _ := c.ReceiveWithTimeout(ctx, time.Minute, &updateInput)
		if ok {
			input.Address = updateInput.Address
		}
	}

	if UPDATE == name {
		// Await update message to update address
		logger.Info("Waiting up to 60 seconds for updated address")
		updatedAddress, err := messages.SetUpdateHandlerForUpdateOrder(ctx)
		if err != nil {
			return nil, err
		}
		ok, _ := workflow.AwaitWithTimeout(ctx, time.Minute, func() bool {
			return *updatedAddress != ""
		})
		if ok {
			input.Address = *updatedAddress
		}
	}

	// Ship order items
	var shipFutures []workflow.Future
	for _, item := range items {
		logger.Info("Shipping item " + item.Description)
		shipFutures = append(shipFutures, shipItemAsync(ctx, input, item, name))
	}

	// Wait for all items to ship
	for _, f := range shipFutures {
		err = f.Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	updateProgress("Order Completed", progress, 100, ctx, 0)

	// Generate trackingId
	trackingId := uuid.New().String()
	output = &app.OrderOutput{
		TrackingId: trackingId,
		Address:    input.Address,
	}

	return output, nil
}

func updateProgress(orderStatus string, progress *int, value int, ctx workflow.Context, seconds int) {
	sleep(ctx, seconds, progress, value)
	if VISIBILITY == workflow.GetInfo(ctx).WorkflowType.Name {
		workflow.UpsertTypedSearchAttributes(ctx, orderStatusKey.ValueSet(orderStatus))
	}
}

func shipItemAsync(ctx workflow.Context, input app.OrderInput, item app.Item, name string) workflow.Future {
	logger := workflow.GetLogger(ctx)
	var f workflow.Future

	shippingInput := app.ShippingInput{
		Order: input,
		Item:  item,
	}

	if CHILD == name {
		// execute an async child wf to ship the item
		cwo := workflow.ChildWorkflowOptions{
			WorkflowID:        fmt.Sprintf("shipment-%v-%v", input.OrderId, item.Id),
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_TERMINATE,
		}
		ctx = workflow.WithChildOptions(ctx, cwo)
		f = workflow.ExecuteChildWorkflow(ctx, ShippingWorkflow, shippingInput)
		logger.Info("Started Child Workflow: " + cwo.WorkflowID)
	} else if NEXUS == name {
		client := workflow.NewNexusClient(app.GetEnv("TEMPORAL_NEXUS_SHIPPING_ENDPOINT", "shipping-endpoint"), app.ShippingServiceName)

		fut := client.ExecuteOperation(ctx, app.ShippingOperationName, shippingInput, workflow.NexusOperationOptions{})
		f = fut

		var exec workflow.NexusOperationExecution
		fut.GetNexusOperationExecution().Get(ctx, &exec)
		logger.Info("Started Nexus Operation: " + exec.OperationToken)
	} else {
		// execute an async activity to ship the item
		f = workflow.ExecuteActivity(ctx, activities.ShipOrder, shippingInput)
		logger.Info("Started Activity: ShipOrder ")
	}
	return f
}
