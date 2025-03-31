package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	generalConf, devicesConf, err := loadConf("devices.toml")
	if err != nil {
		fmt.Println("err in reading devices ", err)
	}
	deviceInterface := generalConf.INTERFACE
	ipStart := generalConf.IPSTART
	ipEnd := generalConf.IPEND

	defer cleanNetworking(deviceInterface)

	var runningDevices []*Device

	for _, config := range devicesConf {
		action, found := actionMap[config.Action]
		if !found {
			continue
		}
		device, err := createDevice(config.ID, config.Broker, action, deviceInterface, ipStart, ipEnd)
		if err != nil {
			continue
		}
		deviceOn(&device)
		runningDevices = append(runningDevices, &device)
	}

	go startAPI()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down devices...")
	for _, device := range runningDevices {
		deviceOff(device)
	}
}
