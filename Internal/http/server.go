package https

import (
	"fmt"

	"github.com/Veerendra-C/SV-Backend.git/Internal/modals"
	"github.com/Veerendra-C/SV-Backend.git/Internal/routes"
	"github.com/gin-gonic/gin"
)

// to start the server
func StartServer(cfg *modals.Config) {
	router := gin.Default()

	routes.Routes(router)
	
	fmt.Printf("Starting server in %s mode on port %s\n", cfg.Env, cfg.Port)

	err := router.Run(fmt.Sprintf(":%s",cfg.Port))
	if err != nil {
		fmt.Printf("\nCould not start the server on port %s",cfg.Port)
		fmt.Printf("\nError: %s",err.Error())
	}
}