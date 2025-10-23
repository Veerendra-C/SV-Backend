package handlers

import (
	"fmt"
	"net/http"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UploadFileHandler(c *gin.Context) {
	// retriving the userID from the context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		logger.Log.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Deugging: Checking if the userID was fetched or not
	userID := fmt.Sprintf("%v", userIDValue)
	logger.Log.Info("UserID Fetched Successfully", zap.String("userID", userID))

	//receiving file from the user
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.Log.Error("Failed to receive the file", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to receive file"})
		return
	}
	defer file.Close()
	logger.Log.Info("The file is received successfully")

	// calling the function to handle upload
	err = services.HandleFileUpload(userID, file, header)
	if err != nil {
		logger.Log.Error("Failed to upload file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	logger.Log.Info("File uploaded successfully", zap.String("userID", userID), zap.String("filename", header.Filename))
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
