package services

import (
	"encoding/base64"
	"errors"
	"fmt"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/db"
	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Common errors for file fetching
var (
	ErrEmptyFilename    = errors.New("filename is empty")
	ErrEmptyUserID      = errors.New("user ID is invalid")
	ErrFileNotFound     = errors.New("file not found")
	ErrFileAccessDenied = errors.New("access denied to file")
	ErrEmptyEncrKey     = errors.New("encryption key not found")
)

func FetchFileByID(userID int64, filename string) (key []byte, err error) {
	// Input validation
	if userID <= 0 {
		logger.Log.Error("Invalid user ID provided", zap.Int64("userID", userID))
		return nil, ErrEmptyUserID
	}

	if filename == "" {
		logger.Log.Error("Empty filename provided", zap.Int64("userID", userID))
		return nil, ErrEmptyFilename
	}

	logger.Log.Debug("Attempting to fetch file",
		zap.Int64("userID", userID),
		zap.String("filename", filename))

	var file modals.File
	if err := db.DB.Where("owner_id = ? AND filename = ?", userID, filename).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Error("File not found in database",
				zap.Int64("userID", userID),
				zap.String("filename", filename))
			return nil, ErrFileNotFound
		}
		logger.Log.Error("Database query failed",
			zap.Int64("userID", userID),
			zap.String("filename", filename),
			zap.Error(err))
		return nil, fmt.Errorf("database query error: %w", err)
	}

	// Validate encryption key presence
	if file.EncrKey == "" {
		logger.Log.Error("File found but missing encryption key",
			zap.Int64("userID", userID),
			zap.String("filename", filename),
			zap.Uint("fileID", file.ID))
		return nil, ErrEmptyEncrKey
	}

	// Decode the base64 encoded key
	keyBytes, err := base64.StdEncoding.DecodeString(file.EncrKey)
	if err != nil {
		logger.Log.Error("Failed to decode base64 encryption key",
			zap.Int64("userID", userID),
			zap.String("filename", filename),
			zap.Uint("fileID", file.ID),
			zap.Error(err))
		return nil, fmt.Errorf("invalid encryption key format: %w", err)
	}

	// Validate decoded key length
	if len(keyBytes) != 32 {
		logger.Log.Error("Invalid decoded encryption key length",
			zap.Int64("userID", userID),
			zap.String("filename", filename),
			zap.Uint("fileID", file.ID),
			zap.Int("keyLength", len(keyBytes)))
		return nil, fmt.Errorf("invalid encryption key length after decoding: expected 32 bytes, got %d bytes", len(keyBytes))
	}

	logger.Log.Info("File encryption key retrieved and decoded successfully",
		zap.Int64("userID", userID),
		zap.String("filename", filename),
		zap.Uint("fileID", file.ID),
		zap.Int("keyLength", len(keyBytes)))

	return keyBytes, nil
}
