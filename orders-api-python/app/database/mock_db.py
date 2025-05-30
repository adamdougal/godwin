from typing import Dict, List, Optional
import time
from datetime import datetime
from uuid import uuid4

from app.models.order import Order, OrderCreate, OrderItem, OrderItemCreate, OrderStatus


class MockDatabase:
    def __init__(self):
        self.orders: Dict[str, Order] = {}

    def get_order(self, order_id: str) -> Optional[Order]:
        """Get an order by ID"""
        return self.orders.get(order_id)

    def create_order(self, order_create: OrderCreate) -> Order:
        """Create a new order"""
        # Calculate subtotals and total
        items = []
        total = 0.0

        for item_create in order_create.items:
            subtotal = item_create.quantity * item_create.unit_price
            item = OrderItem(
                id=str(uuid4()),
                product_id=item_create.product_id,
                quantity=item_create.quantity,
                unit_price=item_create.unit_price,
                subtotal=subtotal
            )
            items.append(item)
            total += subtotal

        # Use billing address if provided, otherwise use shipping address
        billing_address = order_create.billing_address or order_create.shipping_address

        # Create new order
        order = Order(
            id=str(uuid4()),
            customer_id=order_create.customer_id,
            items=items,
            status=OrderStatus.PENDING,
            total=total,
            shipping_address=order_create.shipping_address,
            billing_address=billing_address,
            created_at=datetime.utcnow(),
            updated_at=datetime.utcnow()
        )

        # Save to "database"
        self.orders[order.id] = order
        return order

    def update_order_status(self, order_id: str, status: OrderStatus) -> Optional[Order]:
        """Update the status of an order"""
        if order_id not in self.orders:
            return None

        # Create a new Order with updated status and timestamp
        old_order = self.orders[order_id]
        updated_order = Order(
            id=old_order.id,
            customer_id=old_order.customer_id,
            items=old_order.items,
            status=status,
            total=old_order.total,
            shipping_address=old_order.shipping_address,
            billing_address=old_order.billing_address,
            created_at=old_order.created_at,
            updated_at=datetime.utcnow()
        )

        # Save updated order
        self.orders[order_id] = updated_order
        return updated_order

    def delete_order(self, order_id: str) -> bool:
        """Delete an order"""
        if order_id not in self.orders:
            return False

        del self.orders[order_id]
        return True

    def list_orders(self,
                   skip: int = 0,
                   limit: int = 100,
                   customer_id: Optional[str] = None,
                   status: Optional[OrderStatus] = None) -> List[Order]:
        """List orders with optional filtering"""
        filtered_orders = list(self.orders.values())

        # Filter by customer_id if provided
        if customer_id is not None:
            filtered_orders = [o for o in filtered_orders if o.customer_id == customer_id]

        # Filter by status if provided
        if status is not None:
            filtered_orders = [o for o in filtered_orders if o.status == status]

        # Sort by created_at (newest first)
        filtered_orders.sort(key=lambda x: x.created_at, reverse=True)

        # Apply pagination
        return filtered_orders[skip:skip + limit]


# Create a singleton instance
db = MockDatabase()

# Add some sample data
def add_sample_orders():
    """Add sample orders for testing"""
    # Sample order 1
    order1 = OrderCreate(
        customer_id="cust123",
        items=[
            OrderItemCreate(product_id="prod1", quantity=2, unit_price=10.99),
            OrderItemCreate(product_id="prod2", quantity=1, unit_price=24.99)
        ],
        shipping_address="123 Main St, Anytown, USA"
    )
    db.create_order(order1)

    # Sample order 2
    order2 = OrderCreate(
        customer_id="cust456",
        items=[
            OrderItemCreate(product_id="prod3", quantity=3, unit_price=5.99)
        ],
        shipping_address="456 Oak Ave, Somewhere, USA",
        billing_address="789 Business Rd, Somewhere, USA"
    )
    db.create_order(order2)

    # Sample order 3 (with different status)
    order3 = OrderCreate(
        customer_id="cust123",
        items=[
            OrderItemCreate(product_id="prod4", quantity=1, unit_price=99.99)
        ],
        shipping_address="123 Main St, Anytown, USA"
    )
    order3_obj = db.create_order(order3)
    db.update_order_status(order3_obj.id, OrderStatus.COMPLETED)
