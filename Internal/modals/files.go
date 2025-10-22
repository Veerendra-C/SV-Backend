package modals

import (
	"time"
)

type File struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID    uint      `gorm:"not null" json:"owner_id"`
	Filename   string    `gorm:"type:varchar(255);not null" json:"filename"`
	StorageURL string    `gorm:"type:text;not null" json:"storage_url"`
	Encrypted  bool      `gorm:"default:false" json:"encrypted"`
	Bucketname string    `gorm:"type:text;not null" json:"bucketname"`
	DbKey      string    `gorm:"column:db_key;type:text;not null" json:"db_key"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Owner      User        `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Shares     []FileShare `gorm:"foreignKey:FileID" json:"shares,omitempty"`
	AccessLogs []AccessLog `gorm:"foreignKey:FileID" json:"access_logs,omitempty"`
}
