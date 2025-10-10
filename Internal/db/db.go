package db

import (
	"fmt"
	"os"
	"time"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	db_host := os.Getenv("DB_HOST")

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",db_user,db_pass,db_host,db_port,db_name)

	db ,err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil{
		logger.Log.Error("Failed to connect database: ",zap.String("Error: ", err.Error()))
		return
	}

	sqlDB , err := db.DB()
	if err != nil {
		logger.Log.Error("Failed to get sqlDB: ",zap.String("Error: ", err.Error()))
		return
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	logger.Log.Info("MySQL database connected successfully")
}