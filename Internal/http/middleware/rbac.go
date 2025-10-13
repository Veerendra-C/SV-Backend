package middleware

import (
	"strings"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequireRole middleware checks if the user has the required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by auth middleware)
		role, exists := c.Get("role")
		if !exists {
			logger.Log.Error("No role found in context. Auth middleware must be used before RequireRole")
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}

		userRole, ok := role.(string)
		if !ok {
			logger.Log.Error("Role is not a string", zap.Any("role", role))
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}

		// Check if user's role is in the allowed roles
		for _, r := range roles {
			if strings.EqualFold(userRole, r) {
				c.Next()
				return
			}
		}

		// If we get here, user's role is not allowed
		logger.Log.Warn("Access denied",
			zap.String("user_role", userRole),
			zap.Strings("required_roles", roles),
			zap.String("path", c.Request.URL.Path),
		)
		c.AbortWithStatusJSON(403, gin.H{
			"error": "You don't have permission to access this resource",
		})
	}
}

// Common role check functions
var (
	RequireAdmin = RequireRole("admin")
	RequireUser  = RequireRole("user", "admin") // admin can do everything a user can
)
