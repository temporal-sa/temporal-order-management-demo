package resources

type OrderInput struct {
	OrderId string `json:"orderId"`
	Address string `json:"address"`
}

type OrderOutput struct {
	TrackingId string `json:"trackingId"`
	Address    string `json:"address"`
}

type UpdateOrder struct {
	Address string `json:"address"`
}

type UpdateOrderInput struct {
	Address string `json:"address"`
}

type Items []Item

type Item struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}
