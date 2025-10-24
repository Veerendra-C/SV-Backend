package utils

import (
	"os"
	"time"

	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"github.com/golang-jwt/jwt/v5"
)

// A secret key for signing and varifing a token
var SecretShareKey = []byte(os.Getenv("SECRETSHAREKEY"))

func GenerateShareToken(fileID, ownerID, recipientID int64, canEdit bool, expiryTime time.Time) (string, error) {
	claims := modals.FileShareClaims{
		FileID:      fileID,
		OwnerID:     ownerID,
		RecipientId: recipientID,
		CanEdit:     canEdit,
		ExpiryTime:  expiryTime.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretShareKey)
}
