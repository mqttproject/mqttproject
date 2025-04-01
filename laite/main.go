package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	if(runtime.GOOS!="linux"){
		fmt.Println("Unsupported OS. This program makes use of the iproute2 utility.");
		return;
	}
	generalConf,devicesConf , err := loadConf("devices.toml")
	if err != nil{
		fmt.Println("err in reading devices ",err);
		return;
	}
	deviceInterface := generalConf.Interface;
	if(!createInterface(deviceInterface)){
	 	return;
	 }
	defer cleanNetworking()

	var runningDevices []*Device

	for _, config := range devicesConf {
		action, found := actionMap[config.Action]
		if !found {
			continue
		}

		device,err := createDevice(config.Id,config.Broker, action)
		if err != nil{
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
