package services

import (
	"net/http"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetSharedFiles(c *gin.Context){
	token := c.Query("token")
	if token == "" {
		logger.Log.Error("Token string is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}

	_, err := utils.ValidateTokens(token)
	if err != nil {
		logger.Log.Error("Failed to validate tokens",zap.Error(err))
		c.JSON(http.StatusBadGateway, gin.H{"message": "Not authorised"})
		return
	}

	//fetch
}
