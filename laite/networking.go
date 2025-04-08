package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)
type VirtualDevice struct {
    InterfaceName string
    IPAddress     string
	Used          bool
}
var physicalInterface string;
var virtualDevices = make(map[string]VirtualDevice)
var currentIP uint32 = ipToInt("10.0.0.1")
var deviceCount int = 0

func ipToInt(ip string) uint32 {
	parts := strings.Split(ip, ".") 
	var ipInt uint32	
	for i, partStr := range parts { 
		partInt, _ := strconv.Atoi(partStr)  
		shift := uint(24-8*i); // 24 , 16 , 8 , 0
		shiftedVal := uint32(partInt) << shift;
		
		ipInt |= shiftedVal;
	}
	return ipInt
}
func intToIP(ipInt uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(ipInt>>24)&255,
		(ipInt>>16)&255,
		(ipInt>>8)&255,
		ipInt&255)
}
func getNextAvailableIP() (string, bool) {
	for attempts := 0; attempts < 16777214; attempts++ {
		ipStr := intToIP(currentIP)
		device, exists := virtualDevices[ipStr]
		if !exists || !device.Used {
			currentIP++
			return ipStr, true
		}

		currentIP++
	}
	return "", false
}


func createVirtualDevice() string {
	ip, ok := getNextAvailableIP()
	if !ok {
		fmt.Println("No available IPs.")
		return ""
	}

	ifaceName := fmt.Sprintf("vdev%d", deviceCount)
	deviceCount++

	addCmd := exec.Command("sudo", "ip", "link", "add", ifaceName, "link", physicalInterface, "type", "macvlan", "mode", "bridge")
	addOutput, err := addCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to create macvlan interface %s: %v\nCommand output: %s\n", ifaceName, err, string(addOutput))
		return ""
	}

	upCmd := exec.Command("sudo", "ip", "link", "set", ifaceName, "up")
	if err := upCmd.Run(); err != nil {
		fmt.Println("Failed to bring interface up:", err)
		return ""
	}

	assignCmd := exec.Command("sudo", "ip", "addr", "add", ip+"/8", "dev", ifaceName)
	if err := assignCmd.Run(); err != nil {
		fmt.Println("Failed to assign IP:", err)
		return ""
	}

	virtualDevices[ip] = VirtualDevice{
		InterfaceName: ifaceName,
		IPAddress:     ip,
		Used:          true,
	}

	fmt.Println("Created virtual device", ifaceName, "with IP", ip)
	return ip
}


func cleanNetworking() {
	for _, device := range virtualDevices {
		deleteCmd := exec.Command("sudo", "ip", "link", "delete", device.InterfaceName)
		if err := deleteCmd.Run(); err != nil {
			fmt.Printf("Failed to delete macvlan interface %s: %v\n", device.InterfaceName, err)
		} else {
			fmt.Printf("Deleted macvlan interface %s\n", device.InterfaceName)
		}
	}
	virtualDevices = make(map[string]VirtualDevice)
	currentIP = ipToInt("10.0.0.1")
	deviceCount = 0
}
