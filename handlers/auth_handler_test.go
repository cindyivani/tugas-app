package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trikunc/go_vibe/config"
	"github.com/trikunc/go_vibe/handlers"
	"github.com/trikunc/go_vibe/middleware"
	"github.com/trikunc/go_vibe/models"
	"github.com/trikunc/go_vibe/utils"
)

func TestUserModelPasswordHashing(t *testing.T) {
	u := models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	err := u.HashPassword("secret123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if u.Password == "secret123" || len(u.Password) == 0 {
		t.Fatalf("Password was not properly hashed")
	}

	if !u.CheckPassword("secret123") {
		t.Fatalf("CheckPassword failed with correct password")
	}

	if u.CheckPassword("wrongpassword") {
		t.Fatalf("CheckPassword returned true for incorrect password")
	}
}

func TestJWTGenerationAndValidation(t *testing.T) {
	user := &models.User{
		ID:       99,
		Username: "jwtuser",
		Email:    "jwtuser@example.com",
		Role:     "member",
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatalf("Generated token is empty")
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != 99 {
		t.Fatalf("Expected UserID 99, got %d", claims.UserID)
	}

	if claims.Username != "jwtuser" {
		t.Fatalf("Expected Username jwtuser, got %s", claims.Username)
	}

	if claims.Email != "jwtuser@example.com" {
		t.Fatalf("Expected Email jwtuser@example.com, got %s", claims.Email)
	}

	// Test invalid token
	_, err = utils.ValidateToken(token + "invalid")
	if err == nil {
		t.Fatalf("Expected error for tampered token, got nil")
	}
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/protected", middleware.AuthMiddleware(), func(c *gin.Context) {
		userID := c.MustGet("userID").(uint)
		username := c.MustGet("username").(string)
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"user_id":  userID,
			"username": username,
		})
	})

	// 1. Missing header
	req1, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for missing header, got %d", w1.Code)
	}

	// 2. Invalid format
	req2, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
	req2.Header.Set("Authorization", "Basic 123456")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for non-Bearer auth, got %d", w2.Code)
	}

	// 3. Valid Bearer Token
	user := &models.User{ID: 42, Username: "middleware_user", Email: "mw@test.com", Role: "member"}
	token, _ := utils.GenerateToken(user)

	req3, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("Expected 200 for valid token, got %d. Body: %s", w3.Code, w3.Body.String())
	}
}

func TestAuthIntegrationRegisterLoginProfile(t *testing.T) {
	db, err := config.InitDB()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	// Auto-migrate User
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("AutoMigrate User failed: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	authHandler := handlers.NewAuthHandler(db)

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)
	r.GET("/api/auth/me", middleware.AuthMiddleware(), authHandler.GetProfile)

	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	testUsername := fmt.Sprintf("testuser_%s", uniqueSuffix)
	testEmail := fmt.Sprintf("user_%s@example.com", uniqueSuffix)
	testPassword := "P@ssword123!"

	// 1. Test Register
	regPayload := handlers.RegisterInput{
		Username: testUsername,
		Email:    testEmail,
		Password: testPassword,
		FullName: "Testing User",
	}
	body, _ := json.Marshal(regPayload)
	reqReg, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)

	if wReg.Code != http.StatusCreated {
		t.Fatalf("Register failed: expected 201, got %d. Body: %s", wReg.Code, wReg.Body.String())
	}

	var regResp struct {
		Message     string `json:"message"`
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
		User        struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	}
	if err := json.Unmarshal(wReg.Body.Bytes(), &regResp); err != nil {
		t.Fatalf("Failed to parse register response: %v", err)
	}

	if regResp.Token == "" {
		t.Fatalf("Expected token in register response")
	}

	if regResp.RedirectURL != "/" {
		t.Fatalf("Expected redirect_url '/', got '%s'", regResp.RedirectURL)
	}

	// 2. Test Duplicate Register (Should fail with 409)
	reqDup, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	reqDup.Header.Set("Content-Type", "application/json")
	wDup := httptest.NewRecorder()
	r.ServeHTTP(wDup, reqDup)

	if wDup.Code != http.StatusConflict {
		t.Fatalf("Expected 409 Conflict for duplicate registration, got %d", wDup.Code)
	}

	// 3. Test Login with Username
	loginPayload := handlers.LoginInput{
		Identifier: testUsername,
		Password:   testPassword,
	}
	loginBody, _ := json.Marshal(loginPayload)
	reqLogin, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)

	if wLogin.Code != http.StatusOK {
		t.Fatalf("Login failed: expected 200, got %d. Body: %s", wLogin.Code, wLogin.Body.String())
	}

	var loginResp struct {
		Message     string `json:"message"`
		Token       string `json:"token"`
		RedirectURL string `json:"redirect_url"`
	}
	if err := json.Unmarshal(wLogin.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	if loginResp.Token == "" {
		t.Fatalf("Expected token in login response")
	}

	if loginResp.RedirectURL != "/" {
		t.Fatalf("Expected redirect_url '/', got '%s'", loginResp.RedirectURL)
	}

	// 4. Test Login with Wrong Password (Should fail with 401)
	wrongLoginPayload := handlers.LoginInput{
		Identifier: testUsername,
		Password:   "wrong-password",
	}
	wrongBody, _ := json.Marshal(wrongLoginPayload)
	reqWrong, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(wrongBody))
	reqWrong.Header.Set("Content-Type", "application/json")
	wWrong := httptest.NewRecorder()
	r.ServeHTTP(wWrong, reqWrong)

	if wWrong.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for wrong password, got %d", wWrong.Code)
	}

	// 5. Test Access Protected /api/auth/me using token
	reqMe, _ := http.NewRequest(http.MethodGet, "/api/auth/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+loginResp.Token)
	wMe := httptest.NewRecorder()
	r.ServeHTTP(wMe, reqMe)

	if wMe.Code != http.StatusOK {
		t.Fatalf("Expected 200 from /api/auth/me, got %d. Body: %s", wMe.Code, wMe.Body.String())
	}

	var meResp struct {
		Status string      `json:"status"`
		User   models.User `json:"user"`
	}
	if err := json.Unmarshal(wMe.Body.Bytes(), &meResp); err != nil {
		t.Fatalf("Failed to parse /api/auth/me response: %v", err)
	}

	if meResp.User.Username != testUsername {
		t.Fatalf("Expected username %s, got %s", testUsername, meResp.User.Username)
	}
}
