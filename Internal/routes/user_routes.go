package routes

import (
	"github.com/Veerendra-C/SV-Backend.git/Internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/api/signup", handlers.SignUp) //api for new users to Singin
	r.POST("/api/login",handlers.UserLogin) //api for registered users
}
