package handlers

import (
	"net/http"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UploadFileHandler(c *gin.Context){
	_ , exists := c.Get("user_id")

	if !exists {
		logger.Log.Error("The user ID is not found.", zap.String("message","The user id is not found in c.get"))
		c.JSON(http.StatusInternalServerError, gin.H{"Error" : "Not Authorised"})
		return
	}

	_, _ , err := c.Request.FormFile("file")

	if err != nil{
		logger.Log.Error("Failed to receive the file", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to upload file"})
		return
	}
}