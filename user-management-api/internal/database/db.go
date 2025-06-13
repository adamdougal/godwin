package database

import (
	"log"
	"user-management-api/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	// Using PostgreSQL with connection pooling
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create database tables manually
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create default admin user if not exists
	createDefaultAdmin()
}

func createDefaultAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		adminUser := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "admin123",
			Role:     "admin",
			IsActive: true,
		}

		if err := DB.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
		}
	}
}
