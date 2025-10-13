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
	return func (c *gin.Context){
		//fetch the JWT token from the client(cookie + Headers)
		tokenstring, err := c.Cookie("Authorization")
		if err != nil || tokenstring == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(tokenstring, "Bearer "){
				tokenstring = strings.TrimPrefix(authHeader, "Bearer")
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"Message" : "Not Authorized"})
				logger.Log.Error("User not authorized" , zap.Error(err))
				c.Abort()
				return
			}
		}

		//Parsing and varification of the JWT teken using the key
		token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid{
			c.JSON(http.StatusUnauthorized, gin.H{"message" : "Not Authorized"})
			logger.Log.Error("Invalid JWT Token", zap.Error(err))
			c.Abort()
			return 
		}

		//validating the cliams of the JWT token
		claims , ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message" : "Not Authorized"})
			logger.Log.Error("JWT Cliams are mismatched ", zap.Error(err))
			c.Abort()
			return
		}

		//adding user info to the gin context
		c.Set("user_id", claims["user_id"])
		c.Set("email",claims["email"])
		c.Set("role", claims["role"])

		c.Next()
	}
}