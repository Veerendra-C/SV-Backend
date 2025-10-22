package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"time"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/repo"
	"github.com/Veerendra-C/SV-Backend.git/Internal/storage"
	"go.uber.org/zap"
)

func HandleFileUpload(UserID string, file multipart.File, header *multipart.FileHeader) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Error("Failed to read the file", zap.Error(err))
		return err
	}

	// get file name and content type(image/pdf/Document etc...)
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, header.Filename)
	contentType := header.Header.Get("Content-Type")

	// uploads file to MinIO bucket
	key, bucket, LastModifiedTime , err := storage.UploadFile("securevault", filename, fileBytes, contentType)
	if err != nil {
		logger.Log.Error("Failed to Upload file to MINIO Bucket", zap.Error(err))
		return err
	}

	// Store file metadata in database
	err = repo.FileMetaData(UserID, filename, contentType, header.Filename, bucket, key, LastModifiedTime)
	if err != nil {
		logger.Log.Error("Failed to store file metadata", zap.Error(err))
		// Deleting the file due to err while uploading the file to minio
		return err
	}

	return nil
}
