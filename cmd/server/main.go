package main

import (
	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/config"
	"github.com/Veerendra-C/SV-Backend.git/Internal/db"
	https "github.com/Veerendra-C/SV-Backend.git/Internal/http"
)

func main() {
	cfg := config.LoadConfig()
	
	logger.InitLogger(cfg.Env)
	db.InitDB()

	https.StartServer(cfg)
}