package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
);


func main() {
	devices , err := loadDevicesFromFile("devices.toml")
	if err != nil{
		fmt.Println("err in reading devices ",err);
	}
	
	for _, config := range devices {
		action, found := actionMap[config.Action]
		if !found {
			continue
		}
		device := createDevice(config.ID,config.Broker, action)
		deviceOn(&device)
	}


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan  //Block until ctrl+c

}