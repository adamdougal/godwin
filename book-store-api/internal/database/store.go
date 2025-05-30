package database

import (
	"errors"
	"sync"
	"time"

	"github.com/godwin/book-store-api/internal/models"
	"github.com/google/uuid"
)

// Store defines the methods for interacting with our data store
type Store interface {
	GetBooks() ([]models.Book, error)
	GetBookByID(id string) (models.Book, error)
	CreateBook(book models.Book) (models.Book, error)
	UpdateBook(id string, book models.Book) (models.Book, error)
	DeleteBook(id string) error
}

// MockStore is an in-memory implementation of the Store interface
type MockStore struct {
	books map[string]models.Book
	mu    sync.RWMutex
}

// NewMockStore returns a new instance of MockStore
func NewMockStore() *MockStore {
	return &MockStore{
		books: make(map[string]models.Book),
	}
}

// GetBooks returns all books in the store
func (m *MockStore) GetBooks() ([]models.Book, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	books := make([]models.Book, 0, len(m.books))
	for _, book := range m.books {
		books = append(books, book)
	}

	return books, nil
}

// GetBookByID retrieves a book by its ID
func (m *MockStore) GetBookByID(id string) (models.Book, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	book, exists := m.books[id]
	if !exists {
		return models.Book{}, errors.New("book not found")
	}

	return book, nil
}

// CreateBook adds a new book to the store
func (m *MockStore) CreateBook(book models.Book) (models.Book, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate a new ID if one wasn't provided
	if book.ID == "" {
		book.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now

	m.books[book.ID] = book

	return book, nil
}

// UpdateBook updates an existing book in the store
func (m *MockStore) UpdateBook(id string, book models.Book) (models.Book, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	existingBook, exists := m.books[id]
	if !exists {
		return models.Book{}, errors.New("book not found")
	}

	// Keep original ID, CreatedAt
	book.ID = existingBook.ID
	book.CreatedAt = existingBook.CreatedAt
	book.UpdatedAt = time.Now()

	m.books[id] = book

	return book, nil
}

// DeleteBook removes a book from the store
func (m *MockStore) DeleteBook(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.books[id]; !exists {
		return errors.New("book not found")
	}

	delete(m.books, id)
	return nil
}
