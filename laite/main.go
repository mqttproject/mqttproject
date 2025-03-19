package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
);


func main() {
	generalConf,devicesConf , err := loadConf("devices.toml")
	if err != nil{
		fmt.Println("err in reading devices ",err);
	}
	fmt.Println("General Config:", generalConf)
	deviceInterface := generalConf.INTERFACE;
	for _, config := range devicesConf {
		action, found := actionMap[config.Action]
		if !found {
			continue
		}
		device := createDevice(config.ID,config.Broker, action,deviceInterface)
		deviceOn(&device)
	}


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan  
	cleanNetworking(deviceInterface)
}