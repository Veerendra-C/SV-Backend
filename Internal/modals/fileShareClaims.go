package modals

import "github.com/golang-jwt/jwt/v5"

type FileShareClaims struct {
	FileID      int64 `json:"file_id"`
	OwnerID     int64 `json:"owner_id"`
	RecipientId int64 `json:"recipient_id"`
	CanEdit     bool  `json:"can_edit"`
	ExpiryTime  int64 `json:"expiry_time"`
	jwt.RegisteredClaims
}