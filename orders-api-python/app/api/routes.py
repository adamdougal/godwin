from typing import List, Optional
from fastapi import APIRouter, HTTPException, Query

from app.database.mock_db import db
from app.models.order import Order, OrderCreate, OrderStatus

router = APIRouter(
    prefix="/orders",
    tags=["orders"],
)


@router.get("/", response_model=List[Order])
async def list_orders(
    skip: int = Query(0, ge=0, description="Number of orders to skip"),
    limit: int = Query(100, ge=1, le=100, description="Maximum number of orders to return"),
    customer_id: Optional[str] = Query(None, description="Filter by customer ID"),
    status: Optional[OrderStatus] = Query(None, description="Filter by order status")
):
    """
    List all orders with optional filtering and pagination.
    """
    return db.list_orders(skip=skip, limit=limit, customer_id=customer_id, status=status)


@router.post("/", response_model=Order, status_code=201)
async def create_order(order: OrderCreate):
    """
    Create a new order.
    """
    return db.create_order(order)


@router.get("/{order_id}", response_model=Order)
async def get_order(order_id: str):
    """
    Get a specific order by ID.
    """
    order = db.get_order(order_id)
    if not order:
        raise HTTPException(status_code=404, detail="Order not found")
    return order


@router.put("/{order_id}/status", response_model=Order)
async def update_order_status(order_id: str, status: OrderStatus):
    """
    Update the status of an order.
    """
    updated_order = db.update_order_status(order_id, status)
    if not updated_order:
        raise HTTPException(status_code=404, detail="Order not found")
    return updated_order


@router.delete("/{order_id}", status_code=204)
async def delete_order(order_id: str):
    """
    Delete an order.
    """
    success = db.delete_order(order_id)
    if not success:
        raise HTTPException(status_code=404, detail="Order not found")
    return None
