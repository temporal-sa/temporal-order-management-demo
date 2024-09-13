from dataclasses import dataclass

@dataclass
class OrderInput:
    OrderId: str
    Address: str

@dataclass
class OrderOutput:
    TrackingId: str
    Address: str

@dataclass
class UpdateOrder:
    Address: str
