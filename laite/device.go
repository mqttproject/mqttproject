package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var devices []*Device
var devicesMutex sync.Mutex

type DeviceAction func(*Device, context.Context)

type Device struct {
	client  mqtt.Client
	on      bool
	action  DeviceAction
	cancel  context.CancelFunc
	context context.Context
}


func createClient(id string, broker string) (mqtt.Client,error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(id)
	virtualIP := createVirtualIP()
	if(virtualIP == ""){
		return nil, fmt.Errorf("failed to create virtual IP for device %s", id)
	}
	localIP := net.ParseIP(virtualIP)
	dialer := &net.Dialer{
		Timeout:   time.Second * 10,
		LocalAddr: &net.TCPAddr{IP: localIP},
		KeepAlive: time.Second * 30,
	}

	opts.SetDialer(dialer)
	newClient := mqtt.NewClient(opts)
	return newClient, nil
}


func createDevice(id string, broker string, action DeviceAction) (*Device,error) {
    devicesMutex.Lock()
	defer devicesMutex.Unlock()
	fmt.Println("Creating a device")
	for _, device := range devices {
		clientID := device.client.OptionsReader()
		if clientID.ClientID() == id {
			fmt.Println("Device already exists:", id)
			return &Device{}, fmt.Errorf("failed to create device: Device already exists");
		}
	}
	client, err := createClient(id, broker)
	if err != nil {
		return &Device{}, fmt.Errorf("failed to create device: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	newDevice := Device{
		client:  client,
		on:      false,
		action:  action,
		cancel:  cancel,
		context: ctx,
	}
	devices = append(devices, &newDevice)
	fmt.Println("Device created on createDevice");
	return &newDevice, nil
}

func deviceOn(d *Device) {
	devicesMutex.Lock()
	defer devicesMutex.Unlock()
	fmt.Println("Device on")
	d.on = true
	go d.action(d, d.context)
}

func deviceOff(d *Device) {
	devicesMutex.Lock()
	defer devicesMutex.Unlock()
	fmt.Println("Device off")
	d.on = false
	d.cancel()
	disconnectDevice(d)
}

func disconnectDevice(d *Device) {
	if d.client.IsConnected() {
		fmt.Println("Disconnecting from MQTT broker...")
		d.client.Disconnect(250)
	} else {
		fmt.Println("Client already disconnected. Maybe it didnt get to connect yet?")
	}
}

func connectDevice(d *Device) {
	clientID := d.client.OptionsReader()
	if token := d.client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Error connecting device %s: %s\n", clientID.ClientID(), token.Error())
	} else {
		fmt.Printf("Device %s connected successfully.\n", clientID.ClientID())
	}
}

func send(d *Device, message string) {
	clientID := d.client.OptionsReader()
	topic := fmt.Sprintf("devices/%s/message", clientID.ClientID())
	if token := d.client.Publish(topic, 0, true, message); token.Wait() && token.Error() != nil {
		fmt.Printf("Error sending message to %s: %s\n", clientID.ClientID(), token.Error())
	} else {
		fmt.Printf("Message sent to %s successfully.\n", clientID.ClientID())
	}
}

func subscribeAndListen(d *Device, msgChannel chan string) {
	clientID := d.client.OptionsReader()
	topic := fmt.Sprintf("devices/%s/message", clientID.ClientID())
	if token := d.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		msgChannel <- string(msg.Payload())
	}); token.Wait() && token.Error() != nil {
		fmt.Printf("Error subscribing to %s: %s\n", topic, token.Error())
	} else {
		fmt.Printf("Subscribed to %s successfully.\n", topic)
	}
}
