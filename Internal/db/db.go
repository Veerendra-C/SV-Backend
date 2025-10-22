package db

import (
	"fmt"
	"os"
	"time"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(env string) error {
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	db_host := os.Getenv("DB_HOST")

	// Check if any required env vars are missing
	if db_user == "" || db_pass == "" || db_host == "" || db_port == "" || db_name == "" {
		return fmt.Errorf("missing required database configuration. Check environment variables")
	}

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_user, db_pass, db_host, db_port, db_name)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		logger.Log.Error("Failed to connect database", zap.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Log.Error("Failed to get sqlDB", zap.Error(err))
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if env == "development" {
		if err := db.AutoMigrate(&modals.User{}, &modals.File{}, &modals.FileShare{}, &modals.AccessLog{}); err != nil {
			logger.Log.Error("AutoMigrate failed", zap.Error(err))
			return fmt.Errorf("auto migration failed: %w", err)
		}
		logger.Log.Info("Auto migration completed successfully")
	}

	DB = db
	logger.Log.Info("MySQL database connected successfully")
	return nil
}
