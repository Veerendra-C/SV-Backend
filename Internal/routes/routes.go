package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// To Write the pages Route
func Routes(r *gin.Engine){
	r.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message" : "Server is up and running",
		})
	})
}