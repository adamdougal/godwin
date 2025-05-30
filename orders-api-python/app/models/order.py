from datetime import datetime
from enum import Enum
from typing import List, Optional
from uuid import uuid4

from pydantic import BaseModel, Field


class OrderStatus(str, Enum):
    PENDING = "pending"
    PROCESSING = "processing"
    COMPLETED = "completed"
    CANCELLED = "cancelled"


class OrderItemCreate(BaseModel):
    product_id: str
    quantity: int
    unit_price: float


class OrderItem(OrderItemCreate):
    id: str = Field(default_factory=lambda: str(uuid4()))
    subtotal: float

    class Config:
        frozen = True


class OrderCreate(BaseModel):
    customer_id: str
    items: List[OrderItemCreate]
    shipping_address: str
    billing_address: Optional[str] = None


class Order(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    customer_id: str
    items: List[OrderItem]
    status: OrderStatus = OrderStatus.PENDING
    total: float
    shipping_address: str
    billing_address: str
    created_at: datetime = Field(default_factory=datetime.utcnow)
    updated_at: datetime = Field(default_factory=datetime.utcnow)

    class Config:
        frozen = True
