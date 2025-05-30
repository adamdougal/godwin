# Book Store API

This is a practice API for a job interview preparation. It's a RESTful API for a bookstore that allows creating, reading, updating, and deleting books.

## Project Structure

```
book-store-api/
├── cmd/
│   └── server/           # Main application entry point
├── internal/             # Private application code
│   ├── database/         # Database interface and implementations
│   ├── handlers/         # HTTP handlers for the API
│   └── models/           # Data models
├── .env                  # Environment variables
└── go.mod                # Go module definition
```

## API Endpoints

- `GET /status` - Check API status
- `GET /books` - Get all books
- `GET /books/:id` - Get a specific book by ID
- `POST /books` - Create a new book
- `PUT /books/:id` - Update an existing book
- `DELETE /books/:id` - Delete a book

## Running the API

1. Ensure you have Go 1.18 or higher installed
2. Clone this repository
3. Run the following commands:

```bash
# Navigate to the project directory
cd book-store-api

# Run the server
go run cmd/server/main.go
```

The server will start on port 8080 by default, or on the port specified in the .env file.
