package app

const ShippingServiceName = "shipping-service"
const ShippingOperationName = "ship-item"

type ShippingInput struct {
	Order OrderInput
	Item  Item
}

type ShippingOutput struct {
	Message string
}
