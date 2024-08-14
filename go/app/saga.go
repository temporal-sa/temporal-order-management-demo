package app

import "go.temporal.io/sdk/workflow"

type Saga struct {
	compensations []any
	arguments     [][]any
}

func (s *Saga) AddCompensation(activity any, parameters ...any) {
	s.compensations = append(s.compensations, activity)
	s.arguments = append(s.arguments, parameters)
}

func (s Saga) Compensate(ctx workflow.Context) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Saga compensations started")

	// Compensate in the reverse order that activies were applied.
	for i := len(s.compensations) - 1; i >= 0; i-- {
		err := workflow.ExecuteActivity(ctx, s.compensations[i], s.arguments[i]...).Get(ctx, nil)
		if err != nil {
			logger.Error("Executing compensation failed", "Error", err)
		}
	}
}
