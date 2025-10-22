package middleware

import (
	"net/http"
	"os"
	"strings"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//fetch the JWT token from the client(cookie + Headers)
		tokenString, err := c.Cookie("Authorization")
		if err != nil || tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
				logger.Log.Error("No authorization token provided")
				c.Abort()
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				logger.Log.Error("Invalid authorization header format")
				c.Abort()
				return
			}

			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
		logger.Log.Info("The tokenstring is successfully fetched")

		//Parsing and varification of the JWT teken using the key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
			logger.Log.Error("Invalid JWT Token", zap.Error(err))
			c.Abort()
			return
		}
		logger.Log.Info("The parsing of JWT tokens was successfull")

		//validating the cliams of the JWT token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized"})
			logger.Log.Error("JWT Cliams are mismatched ", zap.Error(err))
			c.Abort()
			return
		}
		logger.Log.Info("The claims were matched")

		//adding user info to the gin context
		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])

		// Debug logging
		logger.Log.Info("Token claims",
			zap.Any("user_id", claims["user_id"]),
			zap.Any("email", claims["email"]),
			zap.Any("role", claims["role"]))

		// Verify values were set in context
		if uid, exists := c.Get("user_id"); exists {
			logger.Log.Info("User ID in context", zap.Any("user_id", uid))
		} else {
			logger.Log.Error("Failed to set user_id in context")
		}

		if role, exists := c.Get("role"); exists {
			logger.Log.Info("Role in context", zap.Any("role", role))
		} else {
			logger.Log.Error("Failed to set role in context")
		}

		c.Next()
	}
}
