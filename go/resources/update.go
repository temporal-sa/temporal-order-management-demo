package resources

import (
	"errors"
	"regexp"

	"go.temporal.io/sdk/workflow"
)

const timeout = 5

// Setup update handler to update order
func UpdateOrderWithAddress(ctx workflow.Context, address *string, isUpdate *bool) error {
	if err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		"UpdateOrder",
		func(ctx workflow.Context, update UpdateOrder) (string, error) {
			*isUpdate = true
			*address = update.Address
			return *address, nil
		},
		workflow.UpdateHandlerOptions{Validator: validateAddress},
	); err != nil {
		return err
	}

	return ctx.Err()
}

func validateAddress(ctx workflow.Context, update UpdateOrder) error {
	logger := workflow.GetLogger(ctx)

	re := regexp.MustCompile(`^\d+`)
	isMatch := re.MatchString(update.Address)

	if !isMatch {
		msg := "Rejecting order update, invalid address " + update.Address
		logger.Info(msg)
		return errors.New(msg)
	}

	logger.Info("Updating order, address " + update.Address)

	return nil
}
