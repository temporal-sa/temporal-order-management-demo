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
class OrderOutput:
    trackingId: str
    address: str


@dataclass
class UpdateOrderInput:
    Address: str