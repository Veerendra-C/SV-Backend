package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/services"
	"github.com/Veerendra-C/SV-Backend.git/Internal/storage"
	"github.com/Veerendra-C/SV-Backend.git/Internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func RetriveFileHandler(c *gin.Context) {
	// Get parameters from the request
	bucket := c.Param("bucket")
	filename := c.Param("filename")

	userID, exists := c.Get("user_id")
	if !exists {
		logger.Log.Error("User ID not found in context",
			zap.String("bucket", bucket),
			zap.String("filename", filename))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	var ID int64
	switch v := userID.(type) {
	case float64:
		ID = int64(v)
	case int:
		ID = int64(v)
	case int64:
		ID = v
	case string:
		var err error
		ID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			logger.Log.Error("Failed to parse user ID string",
				zap.String("userID", v),
				zap.String("bucket", bucket),
				zap.String("filename", filename),
				zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
	default:
		logger.Log.Error("Invalid user ID type in context",
			zap.String("bucket", bucket),
			zap.String("filename", filename),
			zap.Any("userID", userID),
			zap.String("type", fmt.Sprintf("%T", userID)))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Validate input parameters
	if bucket == "" || filename == "" {
		logger.Log.Error("Missing required parameters",
			zap.String("bucket", bucket),
			zap.String("filename", filename))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters",
		})
		return
	}

	// Log request details
	logger.Log.Info("File retrieval request",
		zap.String("bucket", bucket),
		zap.String("filename", filename))

	// Create context for MinIO operations
	ctx := context.Background()

	// Get the object from MinIO
	obj, err := storage.MinIOGUI.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		if err.Error() == "The specified key does not exist." || err.Error() == "The specified bucket does not exist" {
			logger.Log.Error("File or bucket not found",
				zap.String("bucket", bucket),
				zap.String("filename", filename),
				zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{
				"error": "The requested file was not found",
			})
			return
		}
		logger.Log.Error("Failed to fetch file from MinIO",
			zap.String("bucket", bucket),
			zap.String("filename", filename),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve file",
		})
		return
	}
	defer obj.Close()

	cipherData, err := io.ReadAll(obj)
	if err != nil {
		logger.Log.Error("Failed to read file content from MinIO",
			zap.String("bucket", bucket),
			zap.String("filename", filename),
			zap.Int64("userID", ID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	key, err := services.FetchFileByID(ID, filename)
	if err != nil {
		logger.Log.Error("Failed to fetch encryption key from database",
			zap.String("bucket", bucket),
			zap.String("filename", filename),
			zap.Int64("userID", ID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch file"})
		return
	}

	fileData, err := utils.DecryptFile(cipherData, key)
	if err != nil {
		logger.Log.Error("File decryption failed",
			zap.String("bucket", bucket),
			zap.String("filename", filename),
			zap.Int64("userID", ID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}

	// Set proper content type for PNG files
	contentType := "application/octet-stream"
	if strings.HasSuffix(strings.ToLower(filename), ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		contentType = "application/pdf"
	}

	// Set response headers
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	// Write the file data directly to response
	if _, err = c.Writer.Write(fileData); err != nil {
		logger.Log.Error("Failed to write file to response",
			zap.String("bucket", bucket),
			zap.String("filename", filename),
			zap.Int64("userID", ID),
			zap.String("contentType", contentType),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream file"})
		return
	}

	logger.Log.Info("File retrieved and streamed successfully",
		zap.String("bucket", bucket),
		zap.String("filename", filename),
		zap.Int64("userID", ID),
		zap.Int("fileSize", len(fileData)),
		zap.Int("originalSize", len(cipherData)))
}
