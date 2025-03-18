package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)
type confDevice struct {
	ID     string `toml:"id"`
	Action string `toml:"action"`
	Broker string `toml:"broker"`
}
type Config struct {
	Devices map[string]confDevice `toml:"devices"` 
}
var actionMap = map[string]DeviceAction{
	"coffeeAction": coffeeAction,
}


func loadDevicesFromFile(filePath string) (map[string]confDevice, error) {
	var config Config
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return nil, err
	}
	devicesConfig := make(map[string]confDevice)
	for id, deviceConfig := range config.Devices {
		_, found := actionMap[deviceConfig.Action]
		if !found {
			fmt.Println("Unknown action:", deviceConfig.Action)
			continue
		}
		devicesConfig[id] = confDevice{
			ID:     deviceConfig.ID,
			Action: deviceConfig.Action, 
			Broker: deviceConfig.Broker,
		}
	}

	return devicesConfig, nil
}