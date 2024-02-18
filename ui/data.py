from dataclasses import dataclass
from datetime import datetime, timedelta
from temporalio import activity, exceptions
from typing import Optional

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

