package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"github.com/Veerendra-C/SV-Backend.git/Internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ShareFileHandler(c *gin.Context) {
	// Bind JSON request
	var req modals.Sharerequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid share request format",
			zap.Error(err),
			zap.Any("request", req))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request fields
	if req.FileID <= 0 {
		logger.Log.Error("Invalid file ID in request",
			zap.Int64("fileID", req.FileID))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file ID",
		})
		return
	}

	if req.RecipientID <= 0 {
		logger.Log.Error("Invalid recipient ID in request",
			zap.Int64("recipientID", req.RecipientID))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid recipient ID",
		})
		return
	}

	if req.ExpiresIn <= 0 {
		logger.Log.Error("Invalid expiration time in request",
			zap.Int64("expiresIn", req.ExpiresIn))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Expiration time must be positive",
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Log.Error("User ID missing from context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return
	}

	// Type assert user ID to float64 (JSON number type)
	userIDFloat, ok := userID.(float64)
	if !ok {
		logger.Log.Error("Invalid user ID type in context",
			zap.Any("userID", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Convert float64 to int64
	ID := int64(userIDFloat)

	// Create share and generate token
	logger.Log.Debug("Attempting to share file",
		zap.Int64("fileID", req.FileID),
		zap.Int64("ownerID", ID),
		zap.Int64("recipientID", req.RecipientID),
		zap.Bool("canEdit", req.CanEdit),
		zap.Int64("expiresIn", req.ExpiresIn))

	token, err := services.ShareFiles(
		req.FileID,
		ID,
		req.RecipientID,
		req.CanEdit,
		time.Duration(req.ExpiresIn)*time.Second,
	)

	if err != nil {
		var statusCode int
		var response gin.H

		switch {
		case errors.Is(err, services.ErrInvalidFileID):
			statusCode = http.StatusNotFound
			response = gin.H{"error": "File not found"}
		case errors.Is(err, services.ErrInvalidRecipient):
			statusCode = http.StatusBadRequest
			response = gin.H{"error": "Invalid recipient"}
		case errors.Is(err, services.ErrExpiredShare):
			statusCode = http.StatusBadRequest
			response = gin.H{"error": "Invalid expiration time"}
		default:
			statusCode = http.StatusInternalServerError
			response = gin.H{"error": "Failed to share file"}
		}

		logger.Log.Error("File sharing failed",
			zap.Int64("fileID", req.FileID),
			zap.Int64("ownerID", ID),
			zap.Int64("recipientID", req.RecipientID),
			zap.Error(err))

		c.JSON(statusCode, response)
		return
	}

	shareURL := fmt.Sprintf("http://localhost:8080/api/user/share/%s", token)

	logger.Log.Info("File shared successfully",
		zap.Int64("fileID", req.FileID),
		zap.Int64("ownerID", ID),
		zap.Int64("recipientID", req.RecipientID),
		zap.Bool("canEdit", req.CanEdit),
		zap.Int64("expiresIn", req.ExpiresIn))

	c.JSON(http.StatusOK, gin.H{
		"message":  "File shared successfully",
		"shareURL": shareURL,
	})
}
