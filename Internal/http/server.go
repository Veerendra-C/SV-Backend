package https

import (
	"fmt"

	logger "github.com/Veerendra-C/SV-Backend.git/Internal/Log"
	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"github.com/Veerendra-C/SV-Backend.git/Internal/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// to start the server
func StartServer(cfg *modals.Config) {
	router := gin.Default()

	routes.Routes(router)
	
	logger.Log.Info("Starting server", zap.String("mode", cfg.Env), zap.String("Port", cfg.Port))

	err := router.Run(fmt.Sprintf(":%s",cfg.Port))
	if err != nil {
		fmt.Printf("\nCould not start the server on port %s",cfg.Port)
		fmt.Printf("\nError: %s",err.Error())
	}
}