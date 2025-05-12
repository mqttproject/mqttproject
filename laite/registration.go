
// laite/registration.go - Fixed syntax errors

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Structure for registration request (simplified for the main branch discovery.js)
type RegistrationRequest struct {
	SimulatorId string `json:"simulatorId"`
	Port        int    `json:"port"`
}

// Simple response structure 
type RegistrationResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Simulator struct {
		Id     string `json:"id"`
		Url    string `json:"url"`
		Status string `json:"status"`
	} `json:"simulator"`
}

// Register simulator with the server using the token
func registerWithServer(config Config) error {
	if config.General.ServerUrl == "" {
		return fmt.Errorf("server_url not configured")
	}
	
	if config.General.Token == "" {
		return fmt.Errorf("token not configured - set this from the dashboard")
	}
	
	// Create registration request
	regRequest := RegistrationRequest{
		SimulatorId: config.General.Id,
		Port:        8080,
	}
	
	jsonData, err := json.Marshal(regRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal registration request: %v", err)
	}
	
	// Send registration request with token in header
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", 
		fmt.Sprintf("%s/api/simulators/register", config.General.ServerUrl),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	// Set token in header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.General.Token))
	
	fmt.Printf("Sending registration request to %s\n", req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send registration request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse the response
	var regResponse RegistrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&regResponse); err != nil {
		return fmt.Errorf("failed to decode registration response: %v", err)
	}
	
	if !regResponse.Success {
		return fmt.Errorf("registration failed: %s", regResponse.Message)
	}
	
	fmt.Printf("Successfully registered with server. Simulator ID: %s\n", regResponse.Simulator.Id)
	fmt.Printf("Server recognizes simulator at: %s\n", regResponse.Simulator.Url)
	fmt.Printf("Simulator status: %s\n", regResponse.Simulator.Status)
	
	return nil
}

// Get local IP address (unchanged)
func getLocalIP() (string, error) {
	return getOutboundIP(), nil
}

func getOutboundIP() string {
	// Use the physical interface IP if available
	if physicalInterface != "" {
		cmd := exec.Command("ip", "addr", "show", physicalInterface)
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "inet ") && !strings.Contains(line, "127.0.0.1") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						ip := strings.Split(fields[1], "/")[0]
						return ip
					}
				}
			}
		}
	}
	
	// Fallback to default method
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
