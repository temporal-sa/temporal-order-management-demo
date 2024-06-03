package resources

import (
	"go.temporal.io/sdk/workflow"
)

// Setup query handler for items
func QueryItems(ctx workflow.Context) (*Items, error) {
	log := workflow.GetLogger(ctx)

	itemList := Items{}

	err := workflow.SetQueryHandler(ctx, "getItems", func(input []byte) ([]Item, error) {
		return itemList, nil
	})
	if err != nil {
		log.Error("SetQueryHandler failed for getItems: " + err.Error())
		return &itemList, err
	}

	return &itemList, nil
}

// Setup query handler for progress
func QueryProgress(ctx workflow.Context) (*int, error) {
	log := workflow.GetLogger(ctx)

	progress := 0

	err := workflow.SetQueryHandler(ctx, "getProgress", func(input []byte) (int, error) {
		return progress, nil
	})
	if err != nil {
		log.Error("SetQueryHandler failed for getProgress: " + err.Error())
		return &progress, err
	}

	return &progress, nil
}

// Custom Len Sort Method
func (p Items) Len() int {
	return len(p)
}

// Custom Less Sort Method
func (p Items) Less(i, j int) bool {
	return p[i].Id < p[j].Id
}

// Custom Swap Sort Method
func (p Items) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
