package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
func startAPI(){
	router := gin.New()
	router.GET("/configuration", getConfiguration)
	router.Run("localhost:8080")

}


func getConfiguration(c *gin.Context) {
	generalConf, devicesConf, err := loadConf("devices.toml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"general": generalConf,
		"devices": devicesConf,
	})
}