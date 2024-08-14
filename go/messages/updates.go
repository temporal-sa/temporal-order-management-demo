package messages

import (
	"errors"
	"regexp"

	"go.temporal.io/sdk/workflow"
)

// "UpdateOrder" update handler
func SetUpdateHandlerForUpdateOrder(ctx workflow.Context) (*string, error) {
	logger := workflow.GetLogger(ctx)

	var updatedAddress string

	err := workflow.SetUpdateHandlerWithOptions(
		ctx,
		"UpdateOrder",
		func(ctx workflow.Context, updateInput UpdateOrderInput) (string, error) {
			updatedAddress = updateInput.Address
			return updatedAddress, nil
		},
		workflow.UpdateHandlerOptions{Validator: validateAddress},
	)

	if err != nil {
		logger.Error("SetUpdateHandler failed for UpdateOrder: " + err.Error())
		return nil, err
	}

	return &updatedAddress, nil
}

func validateAddress(ctx workflow.Context, update UpdateOrderInput) error {
	logger := workflow.GetLogger(ctx)

	re := regexp.MustCompile(`^\d+`)
	if !re.MatchString(update.Address) {
		msg := "Rejecting order update, invalid address " + update.Address
		logger.Info(msg)
		return errors.New(msg)
	}

	logger.Info("Updating order, address " + update.Address)
	return nil
}
