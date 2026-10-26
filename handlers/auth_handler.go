package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trikunc/go_vibe/models"
	"github.com/trikunc/go_vibe/utils"
	"gorm.io/gorm"
)

// AuthHandler handles user registration, authentication, and profile endpoints
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// RegisterInput represents the payload for user registration
type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

// LoginInput represents the payload for user login
type LoginInput struct {
	Identifier string `json:"identifier"` // Accepts username or email
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password" binding:"required"`
}

// Register handles creation of new user accounts
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Data pendaftaran tidak valid",
			"details": err.Error(),
		})
		return
	}

	username := strings.ToLower(strings.TrimSpace(input.Username))
	email := strings.ToLower(strings.TrimSpace(input.Email))

	// Check if username already exists
	var existingUser models.User
	if err := h.db.Where("LOWER(username) = ?", username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username sudah digunakan, silakan pilih yang lain",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memeriksa username di database",
		})
		return
	}

	// Check if email already exists
	if err := h.db.Where("LOWER(email) = ?", email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email sudah terdaftar, silakan gunakan email lain atau login",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memeriksa email di database",
		})
		return
	}

	user := models.User{
		Username: username,
		Email:    email,
		FullName: strings.TrimSpace(input.FullName),
		Role:     "member",
	}

	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memproses kata sandi",
		})
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal membuat akun pengguna",
			"details": err.Error(),
		})
		return
	}

	// Generate JWT
	token, err := utils.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Akun berhasil dibuat tetapi gagal membuat token login",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Pendaftaran berhasil",
		"token":        token,
		"redirect_url": "/",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"full_name":  user.FullName,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
	})
}

// Login authenticates a user and returns a signed JWT
func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Format data login tidak valid",
			"details": err.Error(),
		})
		return
	}

	identifier := strings.TrimSpace(input.Identifier)
	if identifier == "" {
		identifier = strings.TrimSpace(input.Username)
	}
	if identifier == "" {
		identifier = strings.TrimSpace(input.Email)
	}

	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username atau email harus diisi",
		})
		return
	}

	identifierLower := strings.ToLower(identifier)

	var user models.User
	err := h.db.Where("LOWER(username) = ? OR LOWER(email) = ?", identifierLower, identifierLower).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Username/Email atau password salah",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Kesalahan database saat mencari akun",
		})
		return
	}

	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Username/Email atau password salah",
		})
		return
	}

	// Generate JWT
	token, err := utils.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menerbitkan token JWT",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login berhasil",
		"token":        token,
		"redirect_url": "/",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}

// GetProfile returns the profile of the currently authenticated user
func (h *AuthHandler) GetProfile(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := val.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id in context"})
		return
	}

	var user models.User
	if err := h.db.Select("id, username, email, full_name, role, created_at, updated_at").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   user,
	})
}
