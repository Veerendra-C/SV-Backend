package routes

import (
	"net/http"

	"github.com/Veerendra-C/SV-Backend.git/Internal/http/handlers"
	"github.com/Veerendra-C/SV-Backend.git/Internal/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func Routes(r *gin.Engine) {
	// Health check route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is up and running",
		})
	})

	// Public routes (no auth required)
	r.POST("/api/signup", handlers.SignUp)   // New users registration
	r.POST("/api/login", handlers.UserLogin) // User authentication

	// Protected routes group
	api := r.Group("/api")
	api.Use(middleware.Middleware()) // Apply authentication middleware
	{
		// Validation endpoint
		api.GET("/validate", handlers.Validator)

		// User routes (accessible by both users and admins)
		userRoutes := api.Group("/user")
		userRoutes.Use(middleware.RequireUser)
		{
			userRoutes.POST("/upload",handlers.UploadFileHandler)
		}

		// Admin only routes
		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.RequireAdmin)
		{
			// Add your admin-level routes here
			// adminRoutes.GET("/users", handlers.ListAllUsers)
			// adminRoutes.DELETE("/users/:id", handlers.DeleteUser)
		}
	}
}
