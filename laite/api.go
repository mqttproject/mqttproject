package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)
func startAPI(){
	router := gin.New()
	router.GET("/configuration", getConfiguration)
	router.POST("/configuration",postConfiguration);
	router.POST("/device/:id", postDevice)
	router.GET("/device/:id", getDevice)
	router.POST("/device/:id/on",signalDeviceOn)
	router.POST("/device/:id/off",signalDeviceOff)
	router.Run("localhost:8080")
}

func getDevice(c *gin.Context) {
	deviceID := c.Param("id")

	for _, device := range devices {
		clientID := device.client.OptionsReader()

		if clientID.ClientID() == deviceID {

			deviceInfo := struct {
				ID     string `json:"id"`
				On     bool   `json:"on"`
				Action string `json:"action"`
			}{
				ID:     deviceID,
				On:     device.on,
				Action: fmt.Sprintf("%T", device.action), 
			}
			c.JSON(http.StatusOK, deviceInfo)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
}
func signalDeviceOn(c *gin.Context){
	deviceID := c.Param("id");

	for _,device := range devices {
		clientID := device.client.OptionsReader()

		if clientID.ClientID() == deviceID {
			if(device.on == true){
				c.JSON(http.StatusOK, "device already on");
				return;
			}
			deviceOn(device);
			c.JSON(http.StatusOK, "device on");
			return;
		}
	}
}
func signalDeviceOff(c *gin.Context){
	deviceID := c.Param("id");

	for _,device := range devices {
		clientID := device.client.OptionsReader()

		if clientID.ClientID() == deviceID {
			if(device.on == false){
				c.JSON(http.StatusOK, "device already off");
				return;
			}
			deviceOff(device);
			c.JSON(http.StatusOK, "device off");
			return;
		}
	}
}
func postDevice(c *gin.Context) {
	deviceID := c.Param("id") 

	var requestData struct {
		Action string `json:"action"`
		Broker string `json:"broker"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		fmt.Println("Failed to parse JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	fmt.Println("Device ID:", deviceID)
	fmt.Println("Action:", requestData.Action)
	fmt.Println("Broker:", requestData.Broker)


	actionFunc, exists := actionMap[requestData.Action]
	if !exists {
		fmt.Println("Unknown action:", requestData.Action)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown action"})
		return
	}


	_,err := createDevice(deviceID,requestData.Broker,actionFunc)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
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
	if newConfig.General.Interface != "" {
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