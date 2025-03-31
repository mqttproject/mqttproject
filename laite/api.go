package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
func startAPI(){
	router := gin.New()
	router.GET("/configuration", getConfiguration)
	router.POST("/configuration",postConfiguration);
	router.Run("localhost:8080")

}


func postConfiguration(c *gin.Context) {
	var newConfig Config

	if err := c.ShouldBindJSON(&newConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	generalConf, devicesConf, err := loadConf("devices.toml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load existing configuration"})
		return
	}
	existingConfig := Config{
		General: generalConf,
		Devices: devicesConf,
	}
	if newConfig.General.INTERFACE != "" {
		existingConfig.General = newConfig.General
	}
	for key, device := range newConfig.Devices {
		existingConfig.Devices[key] = device
	}
	if err := saveConf("devices.toml", existingConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
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