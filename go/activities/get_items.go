package activities

import (
	"context"
	"sort"
	"temporal-order-management/resources"

	"go.temporal.io/sdk/activity"
)

func GetItems(ctx context.Context) (resources.Items, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting list of items")

	// simulate DB query
	simulateExternalOperation(100)

	itemList := resources.Items{
		{Id: 654300, Description: "Table Top", Quantity: 1},
		{Id: 654321, Description: "Table Legs", Quantity: 2},
		{Id: 654322, Description: "Keypad", Quantity: 1},
	}
	sort.Sort(itemList)

	return itemList, nil
}
