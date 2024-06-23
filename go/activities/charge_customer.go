package activities

import (
	"context"
	"errors"
	"temporal-order-management/resources"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

const (
	ErrorChargeAPIUnavailable = "OrderWorkflowAPIFailure"
	ErrorInvalidCreditCard    = "OrderWorkflowNonRecoverableFailure"
)

func ChargeCustomer(ctx context.Context, input resources.OrderInput, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Charge Customer activity started", "orderId", input.OrderId)
	attempt := activity.GetInfo(ctx).Attempt

	// simulate external API call
	error := simulateExternalOperationWithError(1000, name, attempt)
	logger.Info("Simulated call complete", "name", name, "error", error)
	switch error {
	case ErrorChargeAPIUnavailable:
		// a transient error, which can be retried
		logger.Info("Charge Customer API unavailable", "attempt", attempt)
		return "", errors.New("charge customer activity failed, API unavailable")
	case ErrorInvalidCreditCard:
		// a business error, which cannot be retried
		return "", temporal.NewNonRetryableApplicationError("charge customer activity failed", "activityFailure", errors.New("credit card invalid"))
	default:
		// pass through, no error
	}

	return input.OrderId, nil
}

func UndoChargeCustomer(ctx context.Context, input resources.OrderInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Undo Charge Customer activity started", "orderId", input.OrderId)

	// simulate external API call
	simulateExternalOperation(1000)

	return input.OrderId, nil
}
