package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"go.uber.org/zap"
)

// Common decryption errors
var (
	ErrEmptyData    = errors.New("cipher data is empty")
	ErrEmptyKey     = errors.New("encryption key is empty")
	ErrInvalidKey   = errors.New("invalid encryption key size")
	ErrInvalidData  = errors.New("invalid cipher data format")
	ErrDecryptionOp = errors.New("decryption operation failed")
)

func DecryptFile(cipherData, key []byte) ([]byte, error) {
	// Validate inputs
	if len(cipherData) == 0 {
		logger.Log.Error("Decryption failed: empty cipher data")
		return nil, ErrEmptyData
	}
	if len(key) == 0 {
		logger.Log.Error("Decryption failed: empty key")
		return nil, ErrEmptyKey
	}
	if len(key) != 32 { // AES-256 requires 32-byte key
		logger.Log.Error("Decryption failed: invalid key size", zap.Int("keySize", len(key)))
		return nil, ErrInvalidKey
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Log.Error("Failed to create cipher block", zap.Error(err))
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logger.Log.Error("Failed to create GCM mode", zap.Error(err))
		return nil, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	// Validate cipher data length
	nonceSize := aesGCM.NonceSize()
	if len(cipherData) < nonceSize {
		logger.Log.Error("Invalid cipher data: too short",
			zap.Int("expected", nonceSize),
			zap.Int("got", len(cipherData)))
		return nil, ErrInvalidData
	}

	// Extract nonce and encrypted data
	nonce, encryptedData := cipherData[:nonceSize], cipherData[nonceSize:]

	// Decrypt data
	plaintext, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		logger.Log.Error("Decryption operation failed", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrDecryptionOp, err)
	}

	logger.Log.Info("File decryption successful", zap.Int("plaintextSize", len(plaintext)))
	return plaintext, nil
}
