package resources

type OrderInput struct {
	OrderId  string `json:"orderId"`
	Scenario string `json:"scenario"`
}

type OrderOutput struct {
	TrackingId string `json:"trackingId"`
}
