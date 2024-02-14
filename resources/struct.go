package resources

type OrderInput struct {
	OrderId  string `json:"orderId"`
	Scenario string `json:"scenario"`
	Address  string `json:"address"`
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
