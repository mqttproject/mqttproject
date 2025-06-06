package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	uuid "github.com/satori/go.uuid"
)

type confDevice struct {
	Id     string `toml:"id"`
	Action string `toml:"action"`
	Broker string `toml:"broker"`
}
type confGeneral struct {
	Id string `toml:"id"`
	Interface string `toml:"interface"`
}

type Config struct {
	General confGeneral           `toml:"general"`
	Devices map[string]confDevice `toml:"devices"`
}

var actionMap = map[string]DeviceAction{
	"coffeeAction":          coffeeAction,
	"doorLockAction":        doorLockAction,
	"roomTemperatureAction": roomTemperatureAction,
}

func generateIdentity(config *Config) {
	if config.General.Id == "" {
		myuuid := uuid.NewV4()
		config.General.Id = myuuid.String() 
		fmt.Printf("Generated new UUID: %s\n", config.General.Id)
	} else {
		fmt.Printf("Existing UUID: %s\n", config.General.Id)
	}
}


func saveConf(filePath string, config Config) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func loadConf(filePath string) (confGeneral, map[string]confDevice, error) {
	var config Config
	_, err := toml.DecodeFile(filePath, &config)
	if err != nil {
		return confGeneral{}, nil, err
	}

	generalConfig := config.General

	devicesConfig := make(map[string]confDevice)
	for id, deviceConfig := range config.Devices {
		_, found := actionMap[deviceConfig.Action]
		if !found {
			fmt.Println("Unknown action:", deviceConfig.Action)
			continue
		}
		devicesConfig[id] = confDevice{
			Id:     deviceConfig.Id,
			Action: deviceConfig.Action,
			Broker: deviceConfig.Broker,
		}
	}

	return generalConfig, devicesConfig, nil
}
