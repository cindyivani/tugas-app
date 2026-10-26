package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trikunc/go_vibe/handlers"
	"github.com/trikunc/go_vibe/middleware"
	"gorm.io/gorm"
)

// SetupRouter initializes and registers all application routes, middleware, and handlers
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS middleware for web clients
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Static assets
	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")

	// Frontend View Pages
	viewHandler := func(filePath string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.File(filePath)
		}
	}

	r.Match([]string{http.MethodGet, http.MethodHead}, "/", viewHandler("./web/index.html"))
	r.Match([]string{http.MethodGet, http.MethodHead}, "/detail", viewHandler("./web/detail.html"))
	r.Match([]string{http.MethodGet, http.MethodHead}, "/compare", viewHandler("./web/compare.html"))

	// RESTful API endpoints
	proteinHandler := handlers.NewProteinHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	api := r.Group("/api")
	{
		api.GET("/health", proteinHandler.HealthCheck)
		api.GET("/proteins", proteinHandler.GetProteins)
		api.GET("/proteins/:id", proteinHandler.GetProteinByID)
		api.GET("/proteins/compare", proteinHandler.CompareProteins)
		api.GET("/stats", proteinHandler.GetStats)

		// Authentication endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetProfile)
		}
	}

	return r
}
