package resources

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const timer = 5

func ApprovalTimer(ctx workflow.Context) bool {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting timer for 30 seconds")

	timerCtx, cancelTimer := workflow.WithCancel(ctx)
	approvalTimer := workflow.NewTimer(timerCtx, time.Duration(timer*time.Second))

	var timerFired bool
	timerSelector := workflow.NewSelector(timerCtx)
	timerSelector.AddFuture(approvalTimer, func(f workflow.Future) {
		err := f.Get(timerCtx, nil)
		logger := workflow.GetLogger(ctx)

		if err == nil {
			logger.Info("Timer fired, time exceeded")
			timerFired = true
		} else if ctx.Err() != nil {
			logger.Info("Timer canceled")
		}
	})

	timerSelector.Select(timerCtx)

	// Wait for either timer fired or update
	err := workflow.Await(ctx, func() bool {
		if timerFired {
			return true
		}

		return false
	})

	if err != nil {
		logger.Error("Error waiting for timer: " + err.Error())
		return true
	}

	// return back to workflow
	if timerFired {
		return true
	} else {
		// cancel timer
		cancelTimer()

		return false
	}
}

func boolPointer(b bool) *bool {
	return &b
}
