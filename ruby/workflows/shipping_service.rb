module Workflows
  # Nexus service contract for shipping. The names must match the handler
  # (see go/app/service.go): service "shipping-service", operation "ship-item".
  module ShippingService
    NAME = 'shipping-service'
    SHIP_ITEM = 'ship-item'
    ENDPOINT = ENV.fetch('TEMPORAL_NEXUS_SHIPPING_ENDPOINT', 'shipping-endpoint')

    # The handler expects Go's ShippingInput JSON, so send a plain hash rather than a Struct,
    # which would serialize in Ruby's json_class format
    def self.ship_item_input(order, item)
      {
        'Order' => { 'OrderId' => order.order_id, 'Address' => order.address },
        'Item' => { 'id' => item.id, 'description' => item.description, 'quantity' => item.quantity }
      }
    end
  end
end
