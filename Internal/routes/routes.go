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
			//POST methods
			userRoutes.POST("/upload",handlers.UploadFileHandler) // for uploading a file
			userRoutes.POST("/share", handlers.ShareFileHandler) // for sharing a file

			//GET methods
			userRoutes.GET("/retrive/:bucket/:filename", handlers.RetriveFileHandler) // for streaming a file throught backend
		}

		// Admin only routes
		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.RequireAdmin)
		{
			
		}
	}
}
