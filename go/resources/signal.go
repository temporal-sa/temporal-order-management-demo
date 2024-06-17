package resources

import (
	"go.temporal.io/sdk/workflow"
)

const SignalOrderWithAddressChannelName = "UpdateOrder"

func (signal *UpdateOrder) SignalOrderWithAddress(ctx workflow.Context, selector workflow.Selector) {
	log := workflow.GetLogger(ctx)

	addPlayerSignalChan := workflow.GetSignalChannel(ctx, SignalOrderWithAddressChannelName)
	selector.AddReceive(addPlayerSignalChan, func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &signal)
		log.Info("Recieved signal to update order with address: " + signal.Address)
	})
}
