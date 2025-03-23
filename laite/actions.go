package main

import (
	"context"
	"fmt"
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
	connectDevice(d)
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
