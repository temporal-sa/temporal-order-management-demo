package workflows

import (
	"fmt"
	"temporal-order-management/activities"
	"temporal-order-management/app"
	"temporal-order-management/messages"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	BUG   = "OrderWorkflowRecoverableFailure"
	CHILD = "OrderWorkflowChildWorkflow"
)

func OrderWorkflowScenarios(ctx workflow.Context, input app.OrderInput) (output *app.OrderOutput, err error) {
	name := workflow.GetInfo(ctx).WorkflowType.Name
	logger := workflow.GetLogger(ctx)
	logger.Info("Processing order started", "orderId", input.OrderId)

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

	// Check fraud
	var cfresult string
	err = workflow.ExecuteActivity(ctx, activities.CheckFraud, input).Get(ctx, &cfresult)
	if err != nil {
		return nil, err
	}

	updateProgress(progress, 25, ctx, 1)

	// Prepare shipment
	saga.AddCompensation(activities.UndoPrepareShipment, input)
	var psresult string
	err = workflow.ExecuteActivity(ctx, activities.PrepareShipment, input).Get(ctx, &psresult)
	if err != nil {
		return nil, err
	}

	updateProgress(progress, 50, ctx, 1)

	// Charge customer
	saga.AddCompensation(activities.UndoChargeCustomer, input)
	var ccresult string
	err = workflow.ExecuteActivity(ctx, activities.ChargeCustomer, input, name).Get(ctx, &ccresult)
	if err != nil {
		return nil, err
	}

	updateProgress(progress, 75, ctx, 3)

	if BUG == name {
		// Simulate bug
		panic("Simulated bug - fix me!")
	}

	// Ship orders
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

	updateProgress(progress, 100, ctx, 1)

	// Generate trackingId
	trackingId := uuid.New().String()
	output = &app.OrderOutput{
		TrackingId: trackingId,
		Address:    input.Address,
	}

	return output, nil
}

func shipItemAsync(ctx workflow.Context, input app.OrderInput, item app.Item, name string) workflow.Future {
	var f workflow.Future
	if CHILD == name {
		// execute an async child wf to ship the item
		cwo := workflow.ChildWorkflowOptions{
			WorkflowID:        fmt.Sprintf("shipment-%v-%v", input.OrderId, item.Id),
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_TERMINATE,
		}
		ctx = workflow.WithChildOptions(ctx, cwo)
		f = workflow.ExecuteChildWorkflow(ctx, ShippingChildWorkflow, input)
	} else {
		// execute an async activity to ship the item
		f = workflow.ExecuteActivity(ctx, activities.ShipOrder, input, item)
	}
	return f
}
