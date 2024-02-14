package resources

import (
	"regexp"

	"go.temporal.io/sdk/workflow"
)

const timeout = 5

// Setup update handler to update order
func UpdateOrderWithAddress(ctx workflow.Context, address *string, isUpdate *bool) error {
	if err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		"UpdateOrder",
		func(ctx workflow.Context, update UpdateOrder) error {
			*isUpdate = true
			*address = update.Address
			return nil
		},
		workflow.UpdateHandlerOptions{Validator: validateAddress},
	); err != nil {
		return err
	}

	return ctx.Err()
}

func validateAddress(ctx workflow.Context, update UpdateOrder) error {
	log := workflow.GetLogger(ctx)

	re := regexp.MustCompile(`^\d+`)
	isMatch := re.MatchString(update.Address)

	if !isMatch {
		log.Debug("Rejecting order update, invalid address", update.Address)
	} else {
		log.Debug("Updating order, address", update.Address)
	}

	return nil
}
