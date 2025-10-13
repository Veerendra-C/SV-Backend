package main

import (
	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/config"
	"github.com/Veerendra-C/SV-Backend.git/Internal/db"
	https "github.com/Veerendra-C/SV-Backend.git/Internal/http"
	"github.com/Veerendra-C/SV-Backend.git/Internal/storage"
)

func main() {
	cfg := config.LoadConfig()
	
	logger.InitLogger(cfg.Env)
	storage.InitMinio()
	db.InitDB(cfg.Env)

	https.StartServer(cfg)
}