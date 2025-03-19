package main

import (
	"context"
	"fmt"
	"net"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)
type DeviceAction func(*Device,context.Context)

type Device struct {
	client mqtt.Client;
	on bool;
	action DeviceAction;
	cancel  context.CancelFunc; 
	context context.Context;   
}




func createClient(id string, broker string,deviceInterface string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(id)
	localIP := net.ParseIP(createVirtualIP(deviceInterface))
	dialer := &net.Dialer{
		Timeout:       time.Second * 10, 
		LocalAddr: &net.TCPAddr{IP: localIP},            
		KeepAlive:     time.Second * 30, 
	}

	opts.SetDialer(dialer)
	newClient := mqtt.NewClient(opts)
	return newClient
}



func createDevice(id string,broker string ,action DeviceAction,deviceInterface string) Device {
	fmt.Println("Creating a device");
	ctx, cancel := context.WithCancel(context.Background())
	newDevice := Device{
		client: createClient(id,broker,deviceInterface),
		on:false,
		action: action,
		cancel:cancel,
		context:ctx,
	}
	return newDevice
}

func deviceOn(d *Device){
	fmt.Println("Device on");
	d.on = true;
	go d.action(d, d.context);
}

func deviceOff(d *Device){
	fmt.Println("Device off");
	d.on = false;
	d.cancel()
	disconnectDevice(d);
}



func disconnectDevice(d *Device){
    if d.client.IsConnected() {
        fmt.Println("Disconnecting from MQTT broker...")
        d.client.Disconnect(250)
    }else{
		fmt.Println("Client already disconnected. Maybe it didnt get to connect yet?");
	}	
}


func connectDevice(d *Device) {
	clientID := d.client.OptionsReader();
	if token := d.client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Error connecting device %s: %s\n", clientID.ClientID(), token.Error())
	} else {
		fmt.Printf("Device %s connected successfully.\n", clientID.ClientID())
	}
}