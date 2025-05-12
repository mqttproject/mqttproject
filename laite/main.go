// laite/main.go - Updated for token-based registration

package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("Unsupported OS. This program makes use of the iproute2 utility.")
		return
	}
	
	// create a sqlite3 database
	err := createDatabase()
	if err != nil {
		fmt.Println("Error creating database:", err)
		return
	}
	
	generalConf, devicesConf, err := loadConf("devices.toml")
	if err != nil {
		fmt.Println("err in reading devices ", err)
		return
	}
	
	config := Config{
		General: generalConf,
		Devices: devicesConf,
	}
	
	generateIdentity(&config)
	if err := saveConf("devices.toml", config); err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
	}
	
	// Debug output to check configuration
	fmt.Printf("Server URL configured: %s\n", config.General.ServerUrl)
	fmt.Printf("Token configured: %s\n", config.General.Token)
	
	// Register with server if configured
	if config.General.ServerUrl != "" && config.General.Token != "" {
		fmt.Println("Attempting to register with server...")
		if err := registerWithServer(config); err != nil {
			fmt.Printf("Warning: Failed to register with server: %v\n", err)
			fmt.Println("Continuing without server registration...")
		} else {
			fmt.Println("Registration successful!")
		}
	} else {
		if config.General.ServerUrl == "" {
			fmt.Println("No server URL configured, running in standalone mode")
		} else {
			fmt.Println("No token configured - set this from the dashboard")
		}
	}
	
	deviceInterface := generalConf.Interface
	if deviceInterface != "" {
		physicalInterface = deviceInterface
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

	go startAPI()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// close database conn
	err = closeDatabase()
	if err != nil {
		fmt.Println("Error closing database:", err)
		return
	}

	fmt.Println("Shutting down devices...")
	for _, device := range devices {
		deviceOff(device)
	}
	cleanNetworking()
	fmt.Println("Exiting program...")
}
