package main

import (
	"fmt"
	"os/exec"
)

var ipStart = 50;
var ipCounter = 0

func createVirtualIP(deviceInterface string) string {
    ip := fmt.Sprintf("192.168.100.%d", ipStart+ipCounter)
    ipCounter++
	cmd := exec.Command("sudo", "ip", "addr", "add", ip+"/24", "dev", deviceInterface)
    err := cmd.Run()
    if err != nil {
        fmt.Println("Failed to assign virtual IP:", err)
    }
    return ip;
}

func cleanVirtualIP(ip string,deviceInterface string) {
    cmd := exec.Command("sudo", "ip", "addr", "del", ip+"/24", "dev", deviceInterface)
    err := cmd.Run()
    if err != nil {
        fmt.Println("Failed to remove virtual IP:", err)
    } else {
        fmt.Println("Successfully removed virtual IP:", ip)
    }
}

func cleanNetworking(deviceInterface string) {

    for i := 0; i < ipCounter; i++ {
        ip := fmt.Sprintf("192.168.100.%d", ipStart+i)
        cleanVirtualIP(ip,deviceInterface)
    }
    ipCounter = 0
}