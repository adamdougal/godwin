package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"
	"user-management-api/internal/database"
	"user-management-api/internal/handlers"
	"user-management-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup test app and MongoDB connection
func setupTestApp() (*fiber.App, *gorm.DB) {
	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	// Set the global database
	database.DB = db

	// Setup Fiber app
	app := fiber.New()

	// Setup routes
	app.Post("/register", handlers.RegisterUser)
	app.Post("/login", handlers.LoginUser)

	api := app.Group("/api/v1", handlers.AuthMiddleware)
	admin := api.Group("/admin", handlers.AdminMiddleware)
	admin.Get("/users", handlers.GetUsers)
	api.Post("/updateUser/:id", handlers.UpdateUser)

	return app, db
}

func TestRegisterUser_Success(t *testing.T) {
	app, db := setupTestApp()
	defer db.Exec("DELETE FROM users")

	reqBody := models.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("Expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}

	var user models.User
	bodyBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		t.Fatal(err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
	if user.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", user.Role)
	}
	if !user.IsActive {
		t.Error("Expected user to be active")
	}
	if user.Password != "" {
		t.Error("Password should be hidden in JSON response")
	}
}

func TestRegisterUser_DuplicateUsername(t *testing.T) {
	app, db := setupTestApp()
	defer db.Exec("DELETE FROM users")

	// Create first user
	user := models.User{
		Username: "duplicate",
		Email:    "first@example.com",
		Password: "hashedpassword",
		Role:     "user",
		IsActive: true,
	}
	db.Create(&user)

	// Try to create second user with same username
	reqBody := models.CreateUserRequest{
		Username: "duplicate",
		Email:    "second@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusConflict {
		t.Errorf("Expected status %d for duplicate username, got %d", fiber.StatusConflict, resp.StatusCode)
	}
}

func TestRegisterUser_BugNoValidation(t *testing.T) {
	app, db := setupTestApp()
	defer db.Exec("DELETE FROM users")

	req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}

	t.Log("BUG DEMONSTRATED: API accepts invalid JSON but doesn't properly validate the parsed data structure")
}

func TestLoginUser_BugInactiveUserCanLogin(t *testing.T) {
	app, db := setupTestApp()
	defer db.Exec("DELETE FROM users")

	// First register a user properly to get correct password hash
	reqBody := models.CreateUserRequest{
		Username: "inactiveuser",
		Email:    "inactive@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatal("Failed to create test user")
	}

	// Now deactivate the user
	var user models.User
	db.Where("username = ?", "inactiveuser").First(&user)
	user.IsActive = false
	db.Save(&user)

	// Try to login with inactive user
	loginReq := models.LoginRequest{
		Username: "inactiveuser",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginHttpReq := httptest.NewRequest("POST", "/login", bytes.NewReader(loginBody))
	loginHttpReq.Header.Set("Content-Type", "application/json")

	loginResp, err := app.Test(loginHttpReq, -1)
	if err != nil {
		t.Fatal(err)
	}

	if loginResp.StatusCode == fiber.StatusOK {
		t.Error("BUG FOUND: Inactive users should not be able to login, but they can!")
		t.Log("This demonstrates the security bug where user activation status is not checked during login")
	} else {
		t.Log("Expected behavior: inactive user login was rejected")
	}
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	app, db := setupTestApp()
	defer db.Exec("DELETE FROM users")

	reqBody := models.LoginRequest{
		Username: "nonexistent",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d for invalid credentials, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestGetUsers_BugNoPagination(t *testing.T) {
	app, db := setupTestApp()
	defer db.Exec("DELETE FROM users")

	// Create many users
	for i := 0; i < 50; i++ {
		user := models.User{
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		db.Create(&user)
	}

	// Create admin token
	token := createValidAdminToken()

	req := httptest.NewRequest("GET", "/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	var users []models.User
	bodyBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &users)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) >= 50 {
		t.Error("BUG FOUND: API returns all users without pagination!")
		t.Logf("Returned %d users, should be limited with pagination", len(users))
	}
}

func TestGetUsers_RequiresAuth(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/admin/users", nil)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d for missing auth, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestGetUsers_RequiresAdminRole(t *testing.T) {
	app, _ := setupTestApp()

	// Create regular user token
	token := createValidUserToken()

	req := httptest.NewRequest("GET", "/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusForbidden {
		t.Errorf("Expected status %d for non-admin user, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
}

func TestUpdateUser_BugNoAuthorizationCheck(t *testing.T) {
	app, db := setupTestApp()
	defer db.Exec("DELETE FROM users")

	// Create two users
	user1 := models.User{
		Username: "user1",
		Email:    "user1@example.com",
		Password: "hashedpassword",
		Role:     "user",
		IsActive: true,
	}
	db.Create(&user1)

	user2 := models.User{
		Username: "user2",
		Email:    "user2@example.com",
		Password: "hashedpassword",
		Role:     "user",
		IsActive: true,
	}
	db.Create(&user2)

	// Create token for user2
	token := createValidUserTokenForUser(user2.ID)

	// Try to update user1 using user2's token
	updateReq := models.UpdateUserRequest{
		Email: stringPtr("hacked@example.com"),
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/updateUser/%d", user1.ID), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == fiber.StatusOK {
		t.Error("BUG FOUND: User was able to update another user's information!")
		t.Log("This demonstrates the authorization bug where any user can update any other user")

		// Verify the update actually happened
		var updatedUser models.User
		db.First(&updatedUser, user1.ID)
		if updatedUser.Email == "hacked@example.com" {
			t.Log("Confirmed: Email was actually changed by unauthorized user")
		}
	} else {
		t.Log("Expected behavior: user update was properly rejected")
	}
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	app, _ := setupTestApp()

	token := createValidUserToken()

	updateReq := models.UpdateUserRequest{
		Email: stringPtr("new@example.com"),
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("POST", "/api/v1/updateUser/99999", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status %d for non-existent user, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d for invalid token, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/admin/users", nil)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d for missing token, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

// Helper functions
func createValidAdminToken() string {
	claims := &handlers.Claims{
		UserID: 1,
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte("asd4323eghk!FL'") // Same as in handlers
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

func createValidUserToken() string {
	claims := &handlers.Claims{
		UserID: 1,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte("asd4323eghk!FL'") // Same as in handlers
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

func createValidUserTokenForUser(userID uint) string {
	claims := &handlers.Claims{
		UserID: userID,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte("asd4323eghk!FL'") // Same as in handlers
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

func stringPtr(s string) *string {
	return &s
}
