package modals

import (
	"time"
)

type FileShare struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FileID      uint      `gorm:"not null" json:"file_id"`
	RecipientID uint      `gorm:"not null" json:"recipient_id"`
	WrappedKey  string    `gorm:"type:blob;not null" json:"wrapped_key"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CanEdit     bool      `gorm:"default:false" json:"can_edit"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	File       File `gorm:"foreignKey:FileID" json:"file,omitempty"`
	Recipient  User `gorm:"foreignKey:RecipientID" json:"recipient,omitempty"`
}
