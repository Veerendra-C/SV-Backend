package services

import (
	"errors"
	"fmt"
	"time"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/db"
	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"github.com/Veerendra-C/SV-Backend.git/Internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Common errors for file sharing
var (
	ErrInvalidFileID    = errors.New("invalid file ID")
	ErrInvalidOwnerID   = errors.New("invalid owner ID")
	ErrInvalidRecipient = errors.New("invalid recipient ID")
	ErrExpiredShare     = errors.New("share duration cannot be negative")
	ErrSharingFailed    = errors.New("failed to share file")
	ErrTokenGeneration  = errors.New("failed to generate share token")
)

func ShareFiles(fileID, ownerID, recipientID int64, canEdit bool, expireTime time.Duration) (string, error) {
	// Input validation
	if fileID <= 0 {
		logger.Log.Error("Invalid file ID provided",
			zap.Int64("fileID", fileID))
		return "", ErrInvalidFileID
	}

	if ownerID <= 0 {
		logger.Log.Error("Invalid owner ID provided",
			zap.Int64("ownerID", ownerID))
		return "", ErrInvalidOwnerID
	}

	if recipientID <= 0 {
		logger.Log.Error("Invalid recipient ID provided",
			zap.Int64("recipientID", recipientID))
		return "", ErrInvalidRecipient
	}

	if expireTime <= 0 {
		logger.Log.Error("Invalid expiry duration",
			zap.Duration("expireTime", expireTime))
		return "", ErrExpiredShare
	}

	// Verify recipient exists
	var recipient modals.User
	if err := db.DB.Where("id = ?", recipientID).First(&recipient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Error("Recipient user not found",
				zap.Int64("recipientID", recipientID))
			return "", ErrInvalidRecipient
		}
		logger.Log.Error("Database error while verifying recipient",
			zap.Int64("recipientID", recipientID),
			zap.Error(err))
		return "", fmt.Errorf("database error: %w", err)
	}

	// Verify file exists and belongs to owner
	var file modals.File
	if err := db.DB.Where("id = ? AND owner_id = ?", fileID, ownerID).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Error("File not found or not owned by user",
				zap.Int64("fileID", fileID),
				zap.Int64("ownerID", ownerID))
			return "", fmt.Errorf("file not found or access denied: %w", err)
		}
		logger.Log.Error("Database error while verifying file ownership",
			zap.Int64("fileID", fileID),
			zap.Int64("ownerID", ownerID),
			zap.Error(err))
		return "", fmt.Errorf("database error: %w", err)
	}

	// Create expiry time
	expires := time.Now().Add(expireTime)

	logger.Log.Debug("Creating file share",
		zap.Int64("fileID", fileID),
		zap.Int64("ownerID", ownerID),
		zap.Int64("recipientID", recipientID),
		zap.Bool("canEdit", canEdit),
		zap.Time("expiresAt", expires))

	// Create share record
	share := modals.FileShare{
		FileID:      uint(fileID),
		RecipientID: uint(recipientID),
		CanEdit:     canEdit,
		ExpiresAt:   &expires,
	}

	// Insert into database
	if err := db.DB.Create(&share).Error; err != nil {
		logger.Log.Error("Failed to create share record",
			zap.Int64("fileID", fileID),
			zap.Int64("recipientID", recipientID),
			zap.Error(err))
		return "", fmt.Errorf("%w: %v", ErrSharingFailed, err)
	}

	// Generate JWT token
	token, err := utils.GenerateShareToken(fileID, ownerID, recipientID, canEdit, expires)
	if err != nil {
		logger.Log.Error("Failed to generate share token",
			zap.Int64("fileID", fileID),
			zap.Int64("ownerID", ownerID),
			zap.Int64("recipientID", recipientID),
			zap.Error(err))

		// Attempt to rollback share creation
		if delErr := db.DB.Delete(&share).Error; delErr != nil {
			logger.Log.Error("Failed to rollback share creation",
				zap.Error(delErr))
		}

		return "", fmt.Errorf("%w: %v", ErrTokenGeneration, err)
	}

	logger.Log.Info("File shared successfully",
		zap.Int64("fileID", fileID),
		zap.Int64("ownerID", ownerID),
		zap.Int64("recipientID", recipientID),
		zap.Time("expiresAt", expires))

	return token, nil
}
