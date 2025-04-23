# mqttproject

mqttproject is a testing platform for networks using simulated IoT devices. These devices are implemented with a Go client and communicate with an MQTT broker built with Node.js.

## Project Structure
```
mqttproject/
├── laite/              # IoT device simulator written in Go
│   ├── main.go
│   └── ...             
├── palvelin/           # MQTT broker built with Node.js
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
cd laite 
go build
./laite
```

### Starting the MQTT Broker
Next, set up and start the Node.js MQTT broker:
```bash
cd palvelin
npm ci
npm start
```

### API Documentation
The simulator includes an HTTP REST API that allows control via a separate web interface. For full details on the available endpoints and usage, please refer to the [manual](MANUAL.md).

This API is intended to be used with the companion control interface available at  repository.

The [adminDashboard](https://github.com/mqttproject/adminDashboard) repository provides a web interface for controlling the simulator through this API.