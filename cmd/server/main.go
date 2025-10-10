package main

import (
	"github.com/Veerendra-C/SV-Backend.git/Internal/config"
	https "github.com/Veerendra-C/SV-Backend.git/Internal/http"
)

func main() {
	cfg := config.LoadConfig()

	https.StartServer(cfg)
}