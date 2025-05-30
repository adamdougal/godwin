# Orders API (Python Version)

This is a practice API for a job interview preparation. It's a RESTful API built with FastAPI for managing orders, which allows creating, reading, updating, and deleting orders.

## Project Structure

```
orders-api-python/
├── app/
│   ├── api/            # API routes and endpoints
│   ├── core/           # Core application config
│   ├── database/       # Database interface
│   ├── models/         # Data models
│   └── main.py         # Application entry point
├── tests/              # Test files
├── requirements.txt    # Python dependencies
└── README.md           # Documentation
```

## API Endpoints

- `GET /api/v1/orders` - Get all orders with optional filtering
- `GET /api/v1/orders/{order_id}` - Get a specific order by ID
- `POST /api/v1/orders` - Create a new order
- `PUT /api/v1/orders/{order_id}/status` - Update the status of an order
- `DELETE /api/v1/orders/{order_id}` - Delete an order
- `GET /status` - Health check endpoint

## Running the API

1. Ensure you have Python 3.8+ installed
2. Install dependencies:

```bash
pip install -r requirements.txt
```

3. Run the application:

```bash
uvicorn app.main:app --reload
```

The server will start on port 8000 by default.

## API Documentation

Once the server is running, you can access:
- Swagger UI documentation: http://localhost:8000/docs
- ReDoc documentation: http://localhost:8000/redoc

## Possible Interview Tasks

During the interview, you might be asked to:

1. **Debug Issues**: Test the endpoints and identify any potential bugs in the implementation.

2. **Add Pagination**: Modify or improve the order listings with pagination features.

3. **Implement Filtering**: Add or enhance filtering abilities by different order attributes.

4. **Add Authentication**: Implement a simple JWT-based authentication system.

5. **Implement Validation**: Add proper validation for order fields.

6. **Add Database Integration**: Replace the mock database with an actual database like PostgreSQL or SQLite.

7. **Write Tests**: Create comprehensive unit and integration tests for the API.

## Tips for the Interview

- Thoroughly test each endpoint to understand its behavior
- Think about edge cases in each API operation
- Consider ways to improve performance, security, and error handling
- Be prepared to explain RESTful design principles and best practices with Python and FastAPI
