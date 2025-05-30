package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godwin/book-store-api/internal/database"
	"github.com/godwin/book-store-api/internal/models"
)

// Handler holds dependencies for API handlers
type Handler struct {
	store database.Store
}

// NewHandler returns a new instance of Handler
func NewHandler(store database.Store) *Handler {
	return &Handler{
		store: store,
	}
}

// GetStatus handles GET /status endpoint
func (h *Handler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "book-store-api",
	})
}

// GetBooks handles GET /books endpoint with optional pagination and filtering
func (h *Handler) GetBooks(c *gin.Context) {
	// Get pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Get filtering parameters
	titleFilter := c.Query("title")
	authorFilter := c.Query("author")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	// Convert string parameters to integers with validation
	limit := 10
	offset := 0
	var minPrice, maxPrice float64

	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	if offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil && val >= 0 {
			offset = val
		}
	}

	if minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil && val >= 0 {
			minPrice = val
		}
	}

	if maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil && val >= 0 {
			maxPrice = val
		}
	}

	books, err := h.store.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}

	// Apply filters
	var filteredBooks []models.Book
	for _, book := range books {
		// Title filter
		if titleFilter != "" && !strings.Contains(strings.ToLower(book.Title), strings.ToLower(titleFilter)) {
			continue
		}

		// Author filter
		if authorFilter != "" && !strings.Contains(strings.ToLower(book.Author), strings.ToLower(authorFilter)) {
			continue
		}

		// Price range filter
		if minPriceStr != "" && book.Price < minPrice {
			continue
		}

		if maxPriceStr != "" && book.Price > maxPrice {
			continue
		}

		filteredBooks = append(filteredBooks, book)
	}

	// Apply pagination
	start := offset
	end := offset + limit

	if start > len(filteredBooks) {
		start = len(filteredBooks)
	}
	if end > len(filteredBooks) {
		end = len(filteredBooks)
	}

	pagedBooks := filteredBooks[start:end]

	c.JSON(http.StatusOK, gin.H{
		"total":  len(filteredBooks),
		"limit":  limit,
		"offset": offset,
		"books":  pagedBooks,
	})
}

// GetBook handles GET /books/:id endpoint
func (h *Handler) GetBook(c *gin.Context) {
	id := c.Param("id")

	book, err := h.store.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook handles POST /books endpoint
func (h *Handler) CreateBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book data"})
		return
	}

	createdBook, err := h.store.CreateBook(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, createdBook)
}

// UpdateBook handles PUT /books/:id endpoint
func (h *Handler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book data"})
		return
	}

	updatedBook, err := h.store.UpdateBook(id, book)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, updatedBook)
}

// DeleteBook handles DELETE /books/:id endpoint
func (h *Handler) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	if err := h.store.DeleteBook(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
