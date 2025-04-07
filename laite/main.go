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
	deviceInterface := generalConf.Interface
	if createInterface(deviceInterface) {
		defer cleanNetworking()
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

	fmt.Println("Exiting program...")
}
