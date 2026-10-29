package app

const ShippingServiceName = "shipping-service"
const ShippingOperationName = "ship-item"

type ShippingInput struct {
	Order OrderInput
	Item  Item
	Index int `json:",omitempty"` // item n ships in n*2 seconds; 0 means a random delay
}

type ShippingOutput struct {
	Message string
}
