# mqttproject

mqttproject is a testing platform for networks using simulated IoT devices. These devices are implemented with a Go client and communicate with an MQTT broker built with Node.js. It is designed to test broker performance, network behavior, and device concurrency at scale.

## Project Structure
```
mqttproject/
├── simulator/              # IoT device simulator written in Go
│   ├── main.go
│   └── ...             
├── server/           # MQTT broker built with Node.js
│   ├── main.js
│   ├── package.json
│   └── ...             
```

## Getting Started

### Prerequisites
- Go version 1.23 or above.
- Node.js version 18 or above.

### Running the IoT Device Simulator
First, compile and run the simulator:
```bash
cd simulator 
go build
./simulator
```

### Starting the MQTT Broker
Next, set up and start the Node.js MQTT broker:
```bash
cd server
npm ci
npm start
```

## How It Works
The simulator spins up a configurable number of virtual IoT devices. Each device:
- Connects to the MQTT broker.
- Periodically publishes messages to predefined topics.
- Can subscribe to topics to receive commands or responses.

Concurrency is handled using Go’s lightweight goroutines, making it efficient even with thousands of simulated devices.
In addition to MQTT simulation, the simulator also hosts an HTTP REST API using the gin-gonic/gin framework.

This API allows external tools to:
- Start or stop devices.
- Dynamically update configuration at runtime.
- Restart the simulator.
- Update the actions performed by devices.

### Stress testing
The project can be used to stress devices in between the simulator and the broker with simulated network traffic.
![image](https://github.com/user-attachments/assets/9500a1d5-578f-4fe4-a90d-2dde8eb284dc)


## API Documentation
The simulator includes an HTTP REST API that allows control via a separate web interface. For full details on the available endpoints and usage, please refer to the [manual](MANUAL.md).

This API is intended to be used with the companion control interface available at  repository.

The [adminDashboard](https://github.com/mqttproject/adminDashboard) repository provides a web interface for controlling the simulator through this API.

