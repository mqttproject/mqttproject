package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var assignedIPs = make(map[string]bool)
var currentIP uint32 = ipToInt("10.0.0.1")
var vlanInterface = "vlan0"
var vlanID = "10"

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
	for {
		ipStr := intToIP(currentIP)
		if !assignedIPs[ipStr] {
			assignedIPs[ipStr] = true
			currentIP++
			return ipStr, true
		}
		currentIP++
	}
}

func cleanVirtualIP(ip string) {
	cmd := exec.Command("sudo", "ip", "addr", "del", ip+"/24", "dev", vlanInterface)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to remove virtual IP %s: %v\n", ip, err)
	} else {
		fmt.Printf("Successfully removed virtual IP %s\n", ip)
	}
}


func createVirtualIP() string {
	ipStr, found := getNextAvailableIP()
	if !found {
		fmt.Println("Failed to assign a virtual IP: No available IPs.")
		return ""
	}

	cmd := exec.Command("sudo", "ip", "addr", "add", ipStr+"/24", "dev",vlanInterface)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to assign virtual IP:"+ ipStr , err)
		return ""
	}
	return ipStr
}

func createInterface(physicalInterface string ) bool {
	addIface := exec.Command("sudo", "ip", "link", "add", vlanInterface, "link", physicalInterface, "type", "vlan", "id", vlanID)
	err := addIface.Run()
	if err != nil {
		fmt.Println("Failed to create VLAN interface:", err)
		return false
	}
	upIface := exec.Command("sudo", "ip", "link", "set", vlanInterface, "up")
	err = upIface.Run()
	if err != nil {
		fmt.Println("Failed to bring up VLAN interface:", err)
		return false
	}

	fmt.Println("VLAN interface", vlanInterface, "created successfully on", physicalInterface)
	return true
}

func cleanNetworking() {
	delIface := exec.Command("sudo", "ip", "link", "delete", vlanInterface)
	err := delIface.Run()
	if err != nil {
		fmt.Println("Failed to delete VLAN interface:", err)
	} else {
		fmt.Println("Successfully deleted VLAN interface", vlanInterface)
	}
}

