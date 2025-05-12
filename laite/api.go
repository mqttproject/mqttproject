
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func startAPI() {
	router := gin.New()
	router.GET("/configuration", getConfiguration)
	router.POST("/configuration", postConfiguration)
	router.POST("/device/:id", postDevice)
	router.POST("/devices", postDevices)
	router.GET("/devices", getDevices)
	router.GET("/device/:id", getDevice)
	router.POST("/device/:id/on", signalDeviceOn)
	router.POST("/device/:id/off", signalDeviceOff)
	router.POST("/device/:id/delete", deleteDevice)
	router.POST("/reboot", reboot)
	router.POST("/update", updateApp)
	
	// Add heartbeat endpoint for the server to check
	router.GET("/heartbeat", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	router.Run("0.0.0.0:8080")
}


func getDevices(c *gin.Context) {
	result := make(map[string]interface{})

	devicesMutex.Lock()
	defer devicesMutex.Unlock()
	
	for _, device := range devices {
		clientID := device.client.OptionsReader()
		newDevice := struct {
			ID     string `json:"id"`
			On     bool   `json:"on"`
			Action string `json:"action"`
		}{
			ID:     clientID.ClientID(),
			On:     device.on,
			Action: fmt.Sprintf("%T", device.action),
		}
		result[clientID.ClientID()] = newDevice
	}
	
	c.JSON(http.StatusOK, gin.H{
		"devices": result,
	})
}

func updateApp(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dst := fmt.Sprintf("./%s", file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	err = exec.Command("bash", "./compile.sh").Start()
	if err != nil {
		log.Println("Error executing compile.sh:", err)
		return
	}
	process, _ := os.FindProcess(os.Getpid())
	process.Signal(syscall.SIGTERM)
}

func reboot(c *gin.Context) {
	for _, device := range devices {
		deviceOff(device)
	}

	cleanNetworking() 
	cleanDevices()

	generalConf, devicesConf, err := loadConf("devices.toml")
	if err != nil {
		fmt.Println("err in reading devices ", err)
		return
	}
	deviceInterface := generalConf.Interface
	if (deviceInterface!="") {
		physicalInterface = deviceInterface;
		for _, config := range devicesConf {
			action, found := actionMap[config.Action]
			if !found {
				continue
			}

			device, err := createDevice(config.Id, config.Broker, action)
			if err != nil {
				continue
			}
			deviceOn(device)
		}
	}
}

func postDevices(c *gin.Context){
	type DevicesPayload struct {
		Devices map[string]confDevice `json:"devices"`
	}
	var payload DevicesPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	for _,d := range payload.Devices {
		Action := actionMap[d.Action];
		if Action == nil{
			c.JSON(400, gin.H{"error": fmt.Sprintf("Unknown action for device '%s': %s", d.Id, d.Action)})
			return
		}
		createDevice(d.Id, d.Broker,Action);
	}
	c.JSON(http.StatusOK, "It Probably worked");
}

func deleteDevice(c *gin.Context) {
	devicesMutex.Lock()
	defer devicesMutex.Unlock()
	deviceID := c.Param("id")
	found := false
	filtered := make([]*Device, 0, len(devices))
	for _, device := range devices {
		clientID := device.client.OptionsReader()
		if clientID.ClientID() == deviceID {
			found = true
			devicesMutex.Unlock()
			deviceOff(device)
			devicesMutex.Lock()
			continue 
		}
		filtered = append(filtered, device)
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	devices = filtered 
	c.JSON(http.StatusOK, gin.H{"message": "Device deleted"})
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

// postDevice with server notification
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

	_, err := createDevice(deviceID, requestData.Broker, actionFunc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Notify the server about this new device if configured and authenticated
	generalConf, _, _ := loadConf("devices.toml")
	if generalConf.ServerUrl != ""{
		go func() {
			serverPayload := map[string]interface{}{
				"action":       requestData.Action,
				"broker":       requestData.Broker,
				"simulatorUrl": fmt.Sprintf("http://%s:8080", getOutboundIP()),
				"simulatorId":  generalConf.Id,
			}
			
			jsonData, _ := json.Marshal(serverPayload)
			req, err := http.NewRequest("POST", generalConf.ServerUrl+"/api/device/"+deviceID, bytes.NewBuffer(jsonData))
			if err == nil {
				req.Header.Set("Content-Type", "application/json")
				
				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err == nil {
					defer resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						body, _ := io.ReadAll(resp.Body)
						log.Printf("Failed to notify server about device %s: status %d, body: %s", 
							deviceID, resp.StatusCode, string(body))
					} else {
						log.Printf("Successfully notified server about device %s", deviceID)
					}
				} else {
					log.Printf("Error notifying server about device %s: %v", deviceID, err)
				}
			}
		}()
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
