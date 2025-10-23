package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

func RetriveFileHandler(c *gin.Context) {
	// Get parameters from the request
	bucket := c.Param("bucket")
	filename := c.Param("filename")

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

	// Get object stats
	info, err := obj.Stat()
	if err != nil {
		logger.Log.Error("Failed to get file stats",
			zap.String("bucket", bucket),
			zap.String("filename", filename),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve file information",
		})
		return
	}

	// Set response headers
	c.Header("Content-Type", info.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", info.Size))
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))

	// Stream the file to the response
	if _, err := io.Copy(c.Writer, obj); err != nil {
		logger.Log.Error("Failed to stream file",
			zap.String("bucket", bucket),
			zap.String("filename", filename),
			zap.String("content_type", info.ContentType),
			zap.Int64("size", info.Size),
			zap.Error(err))
		// Note: At this point headers are already sent, so we can't send a JSON response
		// Just log the error and let the connection close
		return
	}

	// Log successful retrieval
	logger.Log.Info("File retrieved successfully",
		zap.String("bucket", bucket),
		zap.String("filename", filename),
		zap.String("content_type", info.ContentType),
		zap.Int64("size", info.Size))
}
