package handlers

import (
	"net/http"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/db"
	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	//get user email/password off req body
	var user modals.User
	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		logger.Log.Error("Failed to bind user request", zap.Error(err))
		return
	}
	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		logger.Log.Error("Failed to hash the user passwords: ", zap.Error(err))
		return
	}
	//create user
	if user.Role == "" {
		user.Role = "user"
	}

	User := modals.User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  string(hash),
		PublicKey: "NOT APPLICAPLE",
		Role:      user.Role,
	}

	result := db.DB.Create(&User)

	if result.Error != nil {
		logger.Log.Error("Failed to create user: ", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user! Try Again"})
		return
	}
	//respond
	c.JSON(http.StatusAccepted, gin.H{"Message": "User created successfully"})
}