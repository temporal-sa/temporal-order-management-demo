import nexusrpc

from shared_objects import ShippingInput


# Nexus service contract for shipping. The names must match the handler
@nexusrpc.service(name="shipping-service")
class ShippingService:
    ship_item: nexusrpc.Operation[ShippingInput, str] = nexusrpc.Operation(name="ship-item")
