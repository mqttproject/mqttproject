package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

var assignedIPs = make(map[string]bool)

func ipToInt(ip string) int {
	parts := strings.Split(ip, ".")
	o1, _ := strconv.Atoi(parts[0])
	o2, _ := strconv.Atoi(parts[1])
	o3, _ := strconv.Atoi(parts[2])
	o4, _ := strconv.Atoi(parts[3])
	return (o1 << 24) | (o2 << 16) | (o3 << 8) | o4
}

func intToIP(ipInt int) string {
	return fmt.Sprintf("%d.%d.%d.%d", 
		(ipInt>>24)&255, 
		(ipInt>>16)&255, 
		(ipInt>>8)&255, 
		ipInt&255)
}

func getNextAvailableIP(startInt, endInt int) (string, bool) {
	for ip := startInt; ip <= endInt; ip++ {
		ipStr := intToIP(ip)
		if !assignedIPs[ipStr] {
			assignedIPs[ipStr] = true 
			return ipStr, true
		}
	}
	return "", false 
}

func validateRange(ipStart string, ipEnd string) bool {
	startIP := net.ParseIP(ipStart)
	endIP := net.ParseIP(ipEnd)

	if startIP == nil || endIP == nil {
		fmt.Println("Invalid IP format")
		return false
	}

	startParts := strings.Split(ipStart, ".")
	endParts := strings.Split(ipEnd, ".")
	if startParts[0] != endParts[0] || startParts[1] != endParts[1] || startParts[2] != endParts[2] {
		fmt.Println("IP range is not within the same /24 subnet")
		return false
	}

	startInt := int(startIP[12])
	endInt := int(endIP[12])

	if startInt > endInt {
		fmt.Println("Start IP is greater than End IP")
		return false
	}

	return true
}


func createVirtualIP(deviceInterface string, ipStart string, ipEnd string) string {
	startInt := ipToInt(ipStart)
	endInt := ipToInt(ipEnd)

	if !validateRange(ipStart,ipEnd) {
		fmt.Println("Invalid IP range")
		return ""
	}

	ipStr, found := getNextAvailableIP(startInt, endInt)
	if !found {
		fmt.Println("No available IPs in the range")
		return ""
	}

	cmd := exec.Command("sudo", "ip", "addr", "add", ipStr+"/24", "dev", deviceInterface)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to assign virtual IP:", err)
		return ""
	}

	return ipStr
}

func cleanVirtualIP(ip string, deviceInterface string) {
	cmd := exec.Command("sudo", "ip", "addr", "del", ip+"/24", "dev", deviceInterface)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to remove virtual IP:", err)
	} else {
		fmt.Println("Successfully removed virtual IP:", ip)
	}
}

func cleanNetworking(deviceInterface string) {
	for ip := range assignedIPs {
		cleanVirtualIP(ip, deviceInterface)
		delete(assignedIPs, ip)
	}
}
