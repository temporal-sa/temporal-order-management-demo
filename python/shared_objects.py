from dataclasses import dataclass


@dataclass
class OrderInput:
    OrderId: str
    Address: str


@dataclass
class OrderItem:
    id: int
    description: str
    quantity: int


@dataclass
class ShippingInput:
    Order: OrderInput
    Item: OrderItem
    Index: int = 0  # item n ships in n*2 seconds; 0 means a random delay


@dataclass
class OrderOutput:
    trackingId: str
    address: str


@dataclass
class UpdateOrderInput:
    Address: str