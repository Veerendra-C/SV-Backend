package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"go.uber.org/zap"
)

func FileEncryption(data []byte) (cipherData, nonce []byte, key string, err error) {
	// Log start of encryption process
	logger.Log.Info("Starting file encryption",
		zap.Int("data_size", len(data)))

	// Generate a random 32-byte key
	keyBytes := make([]byte, 32)
	if n, err := rand.Read(keyBytes); err != nil || n != 32 {
		logger.Log.Error("Failed to generate encryption key",
			zap.Error(err),
			zap.Int("bytes_read", n),
			zap.Int("expected_bytes", 32))
		return nil, nil, "", fmt.Errorf("failed to generate encryption key: %w", err)
	}
	logger.Log.Debug("Generated encryption key successfully")

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		logger.Log.Error("Failed to create AES cipher",
			zap.Error(err),
			zap.Int("key_size", len(keyBytes)))
		return nil, nil, "", fmt.Errorf("failed to create AES cipher: %w", err)
	}
	logger.Log.Debug("Created AES cipher successfully")

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logger.Log.Error("Failed to create GCM mode",
			zap.Error(err))
		return nil, nil, "", fmt.Errorf("failed to create GCM mode: %w", err)
	}
	logger.Log.Debug("Created GCM mode successfully")

	// Generate nonce
	nonce = make([]byte, aesGCM.NonceSize())
	if n, err := rand.Read(nonce); err != nil || n != aesGCM.NonceSize() {
		logger.Log.Error("Failed to generate nonce",
			zap.Error(err),
			zap.Int("bytes_read", n),
			zap.Int("expected_bytes", aesGCM.NonceSize()))
		return nil, nil, "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	logger.Log.Debug("Generated nonce successfully",
		zap.Int("nonce_size", len(nonce)))

	// Encrypt data
	cipherData = aesGCM.Seal(nil, nonce, data, nil)
	key = base64.StdEncoding.EncodeToString(keyBytes)

	// Log encryption success
	logger.Log.Info("File encryption completed successfully",
		zap.Int("original_size", len(data)),
		zap.Int("encrypted_size", len(cipherData)),
		zap.Int("nonce_size", len(nonce)),
		zap.Int("key_size", len(key)))

	return
}
