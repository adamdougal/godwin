# User Management API

A simple RESTful API for user management built with Go, Fiber, GORM, and SQLite. This project is designed for interview practice and contains intentional areas for improvement.

## Tech Stack

- **Framework**: Fiber v2 (Express-inspired Go web framework)
- **ORM**: GORM with SQLite database
- **Authentication**: JWT tokens
- **Password Hashing**: bcrypt
- **Middleware**: CORS, Logger, Custom Auth & Admin middleware

## Why Fiber?

Fiber was chosen for this project instead of Gin to demonstrate:

- **Express.js-like API** - Familiar to developers coming from Node.js
- **High Performance** - Built on top of FastHTTP for speed
- **Built-in Middleware** - Rich ecosystem of middleware (CORS, Logger, etc.)
- **Easy Testing** - Simple to write unit and integration tests
- **Context API** - Locals() for storing request-scoped data
- **Different Architecture** - Shows alternative to Gin's approach

## Project Structure

```
user-management-api/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── database/
│   │   └── db.go           # Database connection and setup
│   ├── handlers/
│   │   └── user.go         # HTTP handlers
│   └── models/
│       └── user.go         # Data models and DTOs
├── go.mod                  # Go modules file
└── README.md              # This file
```

## API Endpoints

### Public Endpoints

- `POST /register` - Register a new user
- `POST /login` - Authenticate and get JWT token

### Protected Endpoints (Require Authentication)

- `GET /api/v1/admin/users` - Get all users (Admin only)
- `POST /api/v1/updateUser/:id` - Update user information

### Utility Endpoints

- `GET /health` - Health check

## Running the API

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on port 8080.

## Usage Examples

### Register a new user
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

### Get all users (Admin only)
```bash
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update user
```bash
curl -X POST http://localhost:8080/api/v1/updateUser/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com"
  }'
```
