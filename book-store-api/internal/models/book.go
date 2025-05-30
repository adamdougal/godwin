package models

import (
	"time"
)

// Book represents a book in the bookstore
type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" binding:"required,min=1,max=200"`
	Author      string    `json:"author" binding:"required,min=1,max=200"`
	ISBN        string    `json:"isbn" binding:"required,isbn"`
	PublishedAt time.Time `json:"published_at" binding:"required"`
	Price       float64   `json:"price" binding:"required,gte=0"`
	Quantity    int       `json:"quantity" binding:"gte=0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BookRequest is used for book creation and update operations
type BookRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=200"`
	Author      string    `json:"author" binding:"required,min=1,max=200"`
	ISBN        string    `json:"isbn" binding:"required,isbn"`
	PublishedAt time.Time `json:"published_at" binding:"required"`
	Price       float64   `json:"price" binding:"required,gte=0"`
	Quantity    int       `json:"quantity" binding:"gte=0"`
}
