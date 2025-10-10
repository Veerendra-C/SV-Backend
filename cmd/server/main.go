package main

import (
	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/config"
	https "github.com/Veerendra-C/SV-Backend.git/Internal/http"
)

func main() {
	cfg := config.LoadConfig()
	
	logger.InitLogger(cfg.Env)

	https.StartServer(cfg)
}