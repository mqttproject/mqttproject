package main

import (
	"context"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)



func coffeeAction(d *Device, ctx context.Context) {
	fmt.Println("Running coffee machine...")
	connectDevice(d)
	msgChannel := make(chan string)
    ClientID := d.client.OptionsReader()
    topic := fmt.Sprintf("devices/%s/message", ClientID.ClientID()) 
    d.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
        message := string(msg.Payload())
       	msgChannel<-message;
    })
	select {
	case <-ctx.Done():
		fmt.Println("Coffee machine action stopped by client.")
	case msg := <-msgChannel:
		fmt.Println("Handling received message:", msg)
	}
}