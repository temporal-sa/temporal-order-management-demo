package resources

import (
	"strconv"
	"time"

	"go.temporal.io/sdk/workflow"
)

const timer = 60

func UpdateApprovalTimer(ctx workflow.Context) (string, bool) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting timer for " + strconv.Itoa(timer) + " seconds")

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

	address := ""
	addressPtr := &address

	isUpdate := boolPointer(false)

	// coroutine for update
	workflow.Go(ctx, func(gCtx workflow.Context) {
		err := UpdateOrderWithAddress(gCtx, addressPtr, isUpdate)

		if err != nil {
			logger.Error("Update failed.", "Error", err)
		}

		if *addressPtr != "" {
			*isUpdate = true
		}

		timerSelector.Select(gCtx)
	})

	// Wait for either timer fired or update
	err := workflow.Await(ctx, func() bool {
		if *isUpdate {
			return true
		}

		if timerFired {
			return true
		}

		return false
	})

	if err != nil {
		logger.Error("Error waiting for timer: " + err.Error())
		return *addressPtr, true
	}

	// return back to workflow
	if timerFired {
		return *addressPtr, true
	} else {
		// cancel timer
		cancelTimer()

		return *addressPtr, false
	}
}

func SignalApprovalTimer(ctx workflow.Context) (string, bool) {
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

	address := ""
	addressPtr := &address

	isSignal := boolPointer(false)

	// coroutine for signal
	workflow.Go(ctx, func(gCtx workflow.Context) {
		signal := UpdateOrder{}

		signalSelector := workflow.NewSelector(gCtx)
		signal.SignalOrderWithAddress(gCtx, signalSelector)

		signalSelector.Select(gCtx)

		if signal.Address != "" {
			*addressPtr = signal.Address
			*isSignal = true
		}
	})

	// coroutine for timer
	workflow.Go(ctx, func(gCtx workflow.Context) {
		timerSelector.Select(gCtx)
	})

	// Wait for either timer fired or signal
	err := workflow.Await(ctx, func() bool {
		if *isSignal {
			return true
		}

		if timerFired {
			return true
		}

		return false
	})

	if err != nil {
		logger.Error("Error waiting for timer: " + err.Error())
		return *addressPtr, true
	}

	// return back to workflow
	if timerFired {
		return *addressPtr, true
	} else {
		// cancel timer
		cancelTimer()

		return *addressPtr, false
	}
}

func boolPointer(b bool) *bool {
	return &b
}
