package activities

import (
	"context"
	"sort"
	"time"

	"github.com/ktenzer/temporal-order-management/resources"
	"go.temporal.io/sdk/activity"
)

func GetItems(ctx context.Context) (resources.Items, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting list of items")

	itemList := resources.Items{
		{Id: 654300, Description: "Table Top", Quantity: 1},
		{Id: 654321, Description: "Table Legs", Quantity: 2},
		{Id: 654322, Description: "Keypad", Quantity: 1},
	}

	sort.Sort(itemList)

	time.Sleep(1 * time.Second)

	return itemList, nil
}
