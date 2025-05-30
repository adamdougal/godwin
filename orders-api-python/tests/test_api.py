import pytest
from fastapi.testclient import TestClient
import uuid

from app.main import app
from app.models.order import OrderStatus

client = TestClient(app)


def test_get_status():
    """Test the status endpoint"""
    response = client.get("/status")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert "version" in data


def test_list_orders():
    """Test listing orders"""
    response = client.get("/api/v1/orders")
    assert response.status_code == 200
    orders = response.json()
    assert isinstance(orders, list)
    if len(orders) > 0:
        assert "id" in orders[0]
        assert "customer_id" in orders[0]


def test_create_order():
    """Test creating a new order"""
    order_data = {
        "customer_id": "test_customer",
        "items": [
            {"product_id": "test_product", "quantity": 1, "unit_price": 15.99}
        ],
        "shipping_address": "Test Address"
    }

    response = client.post("/api/v1/orders", json=order_data)
    assert response.status_code == 201
    order = response.json()
    assert order["customer_id"] == "test_customer"
    assert len(order["items"]) == 1
    assert order["items"][0]["product_id"] == "test_product"
    assert order["total"] == 15.99


def test_get_order():
    """Test getting a specific order"""
    # First create an order
    order_data = {
        "customer_id": "test_customer_get",
        "items": [
            {"product_id": "test_product", "quantity": 2, "unit_price": 10.50}
        ],
        "shipping_address": "Test Address Get"
    }

    create_response = client.post("/api/v1/orders", json=order_data)
    assert create_response.status_code == 201
    created_order = create_response.json()

    # Then get it by ID
    order_id = created_order["id"]
    get_response = client.get(f"/api/v1/orders/{order_id}")
    assert get_response.status_code == 200
    retrieved_order = get_response.json()
    assert retrieved_order["id"] == order_id
    assert retrieved_order["customer_id"] == "test_customer_get"


def test_update_order_status():
    """Test updating an order's status"""
    # First create an order
    order_data = {
        "customer_id": "test_customer_update",
        "items": [
            {"product_id": "test_product", "quantity": 1, "unit_price": 5.99}
        ],
        "shipping_address": "Test Address Update"
    }

    create_response = client.post("/api/v1/orders", json=order_data)
    assert create_response.status_code == 201
    created_order = create_response.json()

    # Then update its status
    order_id = created_order["id"]
    status_response = client.put(f"/api/v1/orders/{order_id}/status", json=OrderStatus.COMPLETED)
    assert status_response.status_code == 200
    updated_order = status_response.json()
    assert updated_order["id"] == order_id
    assert updated_order["status"] == OrderStatus.COMPLETED


def test_delete_order():
    """Test deleting an order"""
    # First create an order
    order_data = {
        "customer_id": "test_customer_delete",
        "items": [
            {"product_id": "test_product_delete", "quantity": 1, "unit_price": 9.99}
        ],
        "shipping_address": "Test Address Delete"
    }

    create_response = client.post("/api/v1/orders", json=order_data)
    assert create_response.status_code == 201
    created_order = create_response.json()

    # Then delete it
    order_id = created_order["id"]
    delete_response = client.delete(f"/api/v1/orders/{order_id}")
    assert delete_response.status_code == 204

    # Verify it's gone
    get_response = client.get(f"/api/v1/orders/{order_id}")
    assert get_response.status_code == 404


def test_get_nonexistent_order():
    """Test getting an order that doesn't exist"""
    non_existent_id = str(uuid.uuid4())
    response = client.get(f"/api/v1/orders/{non_existent_id}")
    assert response.status_code == 404


def test_filter_orders_by_customer():
    """Test filtering orders by customer_id"""
    # Create orders for a specific customer
    customer_id = "filter_test_customer"

    # Create two orders for this customer
    for i in range(2):
        order_data = {
            "customer_id": customer_id,
            "items": [
                {"product_id": f"filter_product_{i}", "quantity": 1, "unit_price": 10.0}
            ],
            "shipping_address": "Filter Test Address"
        }
        client.post("/api/v1/orders", json=order_data)

    # Create an order for a different customer
    other_order_data = {
        "customer_id": "other_customer",
        "items": [
            {"product_id": "other_product", "quantity": 1, "unit_price": 10.0}
        ],
        "shipping_address": "Other Address"
    }
    client.post("/api/v1/orders", json=other_order_data)

    # Test filtering
    response = client.get(f"/api/v1/orders?customer_id={customer_id}")
    assert response.status_code == 200
    orders = response.json()
    assert all(order["customer_id"] == customer_id for order in orders)
    assert len(orders) >= 2  # At least our two orders
