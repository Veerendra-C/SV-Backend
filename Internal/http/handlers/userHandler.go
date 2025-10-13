package handlers

import (
	"net/http"
	"os"
	"time"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/db"
	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func UserLogin(c *gin.Context) {
	//receive the email and password from the user
	var user modals.User
	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		logger.Log.Error("Failed to bind user request during login", zap.Error(err))
		return
	}

	//check for the email's existiance
	var usr_cred modals.User
	db.DB.First(&usr_cred, "email = ?", user.Email)

	if usr_cred.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Could not varify the Email address"})
		logger.Log.Error("Fialed to varify Email", zap.String("User does not exist: ", user.Email))
		return
	}

	//check if the password is correct
	err := bcrypt.CompareHashAndPassword([]byte(usr_cred.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Invalid Password"})
		logger.Log.Error("Failed to vaeify the password", zap.Error(err))
		return
	}

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": usr_cred.ID,
		"email":   usr_cred.Email,
		"role":    usr_cred.Role,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		// Log server-side; do not include secrets or token value
		logger.Log.Error("Failed to create jwt token",
			zap.Error(err),
			zap.Uint("user_id", usr_cred.ID),
			zap.String("email", usr_cred.Email),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create jwt token"})
		return
	}

	// Setting the security of the site
	c.SetSameSite(http.SameSiteStrictMode)

	// cookie 1: For sensitive data.
	c.SetCookie(
		"Authorization", // name
		tokenString,     // value
		3600*12,         // maxAge (12 hours)
		"/",             // path
		"",              // domain (empty = current domain)
		false,            // secure (requires HTTPS)
		true,            // httpOnly (prevents JavaScript access)
	)

	// cookie 2: For non-sensitive data.
	c.SetCookie(
		"user_preferences",
		usr_cred.Role, // just the role for frontend UI decisions
		3600*24*30,    // 30 days
		"/",
		"",
		true,
		false, // not httpOnly - frontend can read
	)

	//response to the client side
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user": gin.H{
			"name": usr_cred.Name, // for display purposes
			"role": usr_cred.Role, // for UI rendering decisions
		},
	})
}
