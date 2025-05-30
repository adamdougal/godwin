package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/godwin/book-store-api/internal/database"
	"github.com/godwin/book-store-api/internal/handlers"
	"github.com/godwin/book-store-api/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file (if it exists)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	// Set production mode if needed
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set up the store
	store := database.NewMockStore()

	// Add some sample data to the store
	addSampleBooks(store)

	// Set up the router
	r := setupRouter(store)

	// Determine port for HTTP service
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT environment variable not set, defaulting to %s", port)
	}

	// Create an HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// setupRouter configures the Gin router with routes and middleware
// addSampleBooks adds some sample data to the store for demonstration purposes
func addSampleBooks(store database.Store) {
	sampleBooks := []models.Book{
		{
			Title:       "Clean Code",
			Author:      "Robert C. Martin",
			ISBN:        "978-0132350884",
			PublishedAt: time.Date(2008, 8, 1, 0, 0, 0, 0, time.UTC),
			Price:       37.49,
			Quantity:    15,
		},
		{
			Title:       "Design Patterns",
			Author:      "Erich Gamma, Richard Helm, Ralph Johnson, John Vlissides",
			ISBN:        "978-0201633610",
			PublishedAt: time.Date(1994, 11, 10, 0, 0, 0, 0, time.UTC),
			Price:       44.99,
			Quantity:    12,
		},
		{
			Title:       "The Pragmatic Programmer",
			Author:      "Andrew Hunt, David Thomas",
			ISBN:        "978-0201616224",
			PublishedAt: time.Date(1999, 10, 30, 0, 0, 0, 0, time.UTC),
			Price:       34.99,
			Quantity:    10,
		},
	}

	for _, book := range sampleBooks {
		_, err := store.CreateBook(book)
		if err != nil {
			log.Printf("Error adding sample book: %v", err)
		}
	}
}

func setupRouter(store database.Store) *gin.Engine {
	r := gin.Default()

	// Create handler with store dependency
	h := handlers.NewHandler(store)

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Routes
	r.GET("/status", h.GetStatus)

	// Book routes
	books := r.Group("/books")
	{
		books.GET("", h.GetBooks)
		books.GET("/:id", h.GetBook)
		books.POST("", h.CreateBook)
		books.PUT("/:id", h.UpdateBook)
		books.DELETE("/:id", h.DeleteBook)
	}

	return r
}
