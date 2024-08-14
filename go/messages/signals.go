package messages

import "go.temporal.io/sdk/workflow"

// "UpdateOrder" signal channel
func GetSignalChannelForUpdateOrder(ctx workflow.Context) workflow.ReceiveChannel {
	return workflow.GetSignalChannel(ctx, "UpdateOrder")
}
