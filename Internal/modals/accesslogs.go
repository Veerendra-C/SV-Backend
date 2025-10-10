package modals

import (
	"time"
)

type AccessLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FileID    uint      `gorm:"not null" json:"file_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Action    string    `gorm:"type:enum('view','download','edit');not null" json:"action"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	File File `gorm:"foreignKey:FileID" json:"file,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
