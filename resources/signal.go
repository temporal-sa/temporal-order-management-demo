package resources

import (
	"go.temporal.io/sdk/workflow"
)

const SignalOrderWithAddressChannelName = "UpdateOrder"

func SignalOrderWithAddress(ctx workflow.Context) string {
	log := workflow.GetLogger(ctx)

	var signal UpdateOrder
	selector := workflow.NewSelector(ctx)
	addPlayerSignalChan := workflow.GetSignalChannel(ctx, SignalOrderWithAddressChannelName)
	selector.AddReceive(addPlayerSignalChan, func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &signal)
		log.Info("Recieved signal to update order with address: " + signal.Address)
	})

	selector.Select(ctx)

	return signal.Address
}
