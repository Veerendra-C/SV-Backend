package main

import (
	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/config"
	"github.com/Veerendra-C/SV-Backend.git/Internal/db"
	https "github.com/Veerendra-C/SV-Backend.git/Internal/http"
	"github.com/Veerendra-C/SV-Backend.git/Internal/storage"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()

	logger.InitLogger(cfg.Env)

	storage.InitMinio()

	if err := db.InitDB(cfg.Env); err != nil {
		logger.Log.Fatal("Failed to initialize database", zap.Error(err))
	}

	https.StartServer(cfg)
}
