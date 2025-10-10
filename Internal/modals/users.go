package modals

import (
	"time"
)

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Email      string    `gorm:"type:varchar(150);unique;not null" json:"email"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"`
	PublicKey  string    `gorm:"type:text" json:"public_key,omitempty"`
	Role       string    `gorm:"type:enum('user','admin');default:'user'" json:"role"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Files      []File       `gorm:"foreignKey:OwnerID" json:"files,omitempty"`
	Shares     []FileShare  `gorm:"foreignKey:RecipientID" json:"shares,omitempty"`
	AccessLogs []AccessLog  `gorm:"foreignKey:UserID" json:"access_logs,omitempty"`
}
