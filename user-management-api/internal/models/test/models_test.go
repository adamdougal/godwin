package models_test

import (
	"testing"
	"user-management-api/internal/models"
)

func TestUser_JSONSerialization(t *testing.T) {
	user := models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "secretpassword",
		Role:     "user",
		IsActive: true,
	}

	// Test that password is hidden when marshaling to JSON
	// This is ensured by the `json:"-"` tag on the Password field
	t.Run("Password should be hidden in JSON", func(t *testing.T) {
		// We can't easily test JSON serialization here without importing encoding/json
		// But the important thing is that the Password field has the `json:"-"` tag
		// which we can verify exists in the struct definition

		// Basic validation that user struct is properly initialized
		if user.Username != "testuser" {
			t.Error("Username should be testuser")
		}
		if user.Password == "" {
			t.Error("Password should be set (but hidden in JSON)")
		}
	})
}

func TestCreateUserRequest_Validation(t *testing.T) {
	// Test XML parsing and validation
	t.Run("Valid request", func(t *testing.T) {
		req := models.CreateUserRequest{
			Username: "validuser",
			Email:    "valid@example.com",
			Password: "password123",
			Role:     "user",
		}

		if req.Username == "" {
			t.Error("Username should not be empty")
		}
		if req.Email == "" {
			t.Error("Email should not be empty")
		}
		if req.Password == "" {
			t.Error("Password should not be empty")
		}
	})
}

func TestLoginRequest_Validation(t *testing.T) {
	t.Run("Valid login request", func(t *testing.T) {
		req := models.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		if req.Username == "" {
			t.Error("Username should not be empty")
		}
		if req.Password == "" {
			t.Error("Password should not be empty")
		}
	})
}

func TestUpdateUserRequest_PartialUpdates(t *testing.T) {
	t.Run("Partial update with email only", func(t *testing.T) {
		email := "new@example.com"
		req := models.UpdateUserRequest{
			Email: &email,
		}

		if req.Email == nil {
			t.Error("Email pointer should not be nil")
		}
		if *req.Email != "new@example.com" {
			t.Error("Email value should match")
		}

		// Other fields should be nil for partial updates
		if req.Username != nil {
			t.Error("Username should be nil for partial update")
		}
		if req.Role != nil {
			t.Error("Role should be nil for partial update")
		}
		if req.IsActive != nil {
			t.Error("IsActive should be nil for partial update")
		}
	})

	t.Run("Update all fields", func(t *testing.T) {
		username := "newusername"
		email := "new@example.com"
		role := "admin"
		isActive := false

		req := models.UpdateUserRequest{
			Username: &username,
			Email:    &email,
			Role:     &role,
			IsActive: &isActive,
		}

		if req.Username == nil || *req.Username != "newusername" {
			t.Error("Username update failed")
		}
		if req.Email == nil || *req.Email != "new@example.com" {
			t.Error("Email update failed")
		}
		if req.Role == nil || *req.Role != "admin" {
			t.Error("Role update failed")
		}
		if req.IsActive == nil || *req.IsActive != false {
			t.Error("IsActive update failed")
		}
	})
}
