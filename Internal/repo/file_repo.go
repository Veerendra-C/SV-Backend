package repo

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Veerendra-C/SV-Backend.git/Internal/db"
	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
)

func FileMetaData(userID, originalFileName, contentType, FileName, bucket, key string, UpdatedTime time.Time) error {
	// Convert userID string to uint
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// store the time and a reference URL 
	now := time.Now()
	url := fmt.Sprintf("securevault/%s", FileName)

	file := modals.File{
		OwnerID:    uint(id),
		Filename:   originalFileName,
		StorageURL: url,
		Encrypted:  false,
		Bucketname: bucket,
		DbKey: key,
		CreatedAt:  now,
		UpdatedAt:  UpdatedTime,
	}

	return db.DB.Create(&file).Error
}
