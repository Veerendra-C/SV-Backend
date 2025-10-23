package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"time"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/repo"
	"github.com/Veerendra-C/SV-Backend.git/Internal/storage"
	"github.com/Veerendra-C/SV-Backend.git/Internal/utils"
	"go.uber.org/zap"
)

func HandleFileUpload(UserID string, file multipart.File, header *multipart.FileHeader) error {
	// Log start of upload process
	logger.Log.Info("Starting file upload process",
		zap.String("userID", UserID),
		zap.String("filename", header.Filename),
		zap.Int64("fileSize", header.Size),
		zap.String("contentType", header.Header.Get("Content-Type")))

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Error("Failed to read file content",
			zap.String("filename", header.Filename),
			zap.Error(err),
			zap.Int64("expectedSize", header.Size),
			zap.Int("bytesRead", len(fileBytes)))
		return fmt.Errorf("failed to read file content: %w", err)
	}
	logger.Log.Debug("File content read successfully",
		zap.Int("bytesRead", len(fileBytes)))

	// Encrypt file content
	encryptedData, nonce, key, err := utils.FileEncryption(fileBytes)
	if err != nil {
		logger.Log.Error("Failed to encrypt file",
			zap.String("filename", header.Filename),
			zap.Error(err))
		return fmt.Errorf("failed to encrypt file: %w", err)
	}
	logger.Log.Debug("File encrypted successfully",
		zap.Int("originalSize", len(fileBytes)),
		zap.Int("encryptedSize", len(encryptedData)))

	// Combine nonce and encrypted data for storage
	combinedData := append(nonce, encryptedData...)

	// Generate timestamped filename and get content type
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, header.Filename)
	contentType := header.Header.Get("Content-Type")

	// Default to application/octet-stream if content type is not provided
	if contentType == "" {
		contentType = "application/octet-stream"
		logger.Log.Debug("Using default content type",
			zap.String("contentType", contentType))
	}

	// Upload to MinIO
	minioKey, bucket, lastModifiedTime, err := storage.UploadFile("securevault", filename, combinedData, contentType)
	if err != nil {
		logger.Log.Error("Failed to upload file to MinIO",
			zap.String("filename", filename),
			zap.String("contentType", contentType),
			zap.Error(err))
		return fmt.Errorf("failed to upload file to storage: %w", err)
	}
	logger.Log.Debug("File uploaded to MinIO successfully",
		zap.String("bucket", bucket),
		zap.String("key", minioKey))

	// Store metadata in database
	err = repo.FileMetaData(UserID, filename, contentType, header.Filename, bucket, key, lastModifiedTime)
	if err != nil {
		logger.Log.Error("Failed to store file metadata",
			zap.String("filename", filename),
			zap.String("userID", UserID),
			zap.Error(err))

		// TODO: Implement rollback - delete file from MinIO
		// Note: Currently, orphaned files may remain in MinIO if metadata storage fails

		return fmt.Errorf("failed to store file metadata: %w", err)
	}

	// Log successful completion
	logger.Log.Info("File upload completed successfully",
		zap.String("userID", UserID),
		zap.String("filename", filename),
		zap.String("bucket", bucket),
		zap.String("contentType", contentType),
		zap.Time("lastModified", lastModifiedTime))

	return nil
}
