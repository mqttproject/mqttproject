package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

func coffeeAction(d *Device, ctx context.Context) {
	fmt.Println("Running coffee machine...")
	connectDevice(d)
	msgChannel := make(chan string)
	subscribeAndListen(d, msgChannel)

	select {
	case <-ctx.Done():
		fmt.Println("Coffee machine action stopped by client.")
	case msg := <-msgChannel:
		fmt.Println("Handling received message:", msg)
		switch msg {
		case "start":
			fmt.Println("Coffee machine started.")
			<-time.After(60 * time.Second)
			fmt.Println("Coffee machine finished.")
			send(d, "done")
		}
	}
}

func doorLockAction(d *Device, ctx context.Context) {
	fmt.Println("Running door lock...")
	connectDevice(d)

	clientID := d.client.OptionsReader()
	deviceId := clientID.ClientID()
	err := rfidStorage.createDeviceTable(deviceId)
	if err != nil {
		fmt.Printf("Error creating RFID table for device %s: %v\n", deviceId, err)
		return
	}

	rfid := doorLockGenerateRFID()
	send(d, fmt.Sprintf("RFID: %d", rfid))

	msgChannel := make(chan string)
	subscribeAndListen(d, msgChannel)

	select {
	case <-ctx.Done():
		fmt.Println("door lock action stopped by client.")
	case msg := <-msgChannel:
		fmt.Println("Handling received message:", msg)
		switch msg {
		case "unlock":
			fmt.Println("Door unlocked.")
			send(d, "unlocked")
		case "lock":
			fmt.Println("Door locked.")
			send(d, "locked")
		}
	}
}

func doorLockGenerateRFID() uint32 {
	rfid := uint32(time.Now().UnixNano())
	fmt.Println("Generated RFID:", rfid)

	return rfid
}

func roomTemperatureAction(d *Device, ctx context.Context) {
	fmt.Println("Running room temperature sensor...")
	connectDevice(d)
	msgChannel := make(chan string)
	subscribeAndListen(d, msgChannel)

	currentTemp := 20.0 + rand.Float64()*5.0
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Temperature sensor stopped.")
			return
		case msg := <-msgChannel:
			fmt.Println("Handling received message:", msg)

		case <-ticker.C:
			change := (rand.Float64() - 0.5) * 0.5
			currentTemp += change
			if currentTemp < 15.0 {
				currentTemp = 15.0
			}
			if currentTemp > 30.0 {
				currentTemp = 30.0
			}
			send(d, fmt.Sprintf("%.1f", currentTemp))
		}
	}
}
