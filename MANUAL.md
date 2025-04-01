# CONFIGURATION

Configuration is done through laite/devices.toml.

The structure of the toml is expected to be the following.

```
#Header "general" will hold the general configuration
[general]

#Interface that you want to use for the creation of the network
interface = ""  

#Header "devices" will hold all the devices 
[devices]

#Each device will have their own subheader.
[devices.mydevice]

#Unique string marking the name of the device. (Will be used as an id by the mqtt-protocol)
id = "" 

#Entry point function that the device will run when its launched.
action = ""

#If the action needs to connect to a broker, this will be the ip that it will try to connect to
broker = ""
```

The program will try to create a virtual lan under the physical interface given through the interface configuration.


## Example configuration

```
[general]
interface = "enp14s0" 

[devices]
[devices.coffee]
id = "coffee"
action = "coffeeAction"
broker = "tcp://192.168.100.11:1883"

[devices.coffee2]
id = "coffee2"
action = "coffeeAction"
broker = "tcp://192.168.100.11:1883"

[devices.coffee3]
id = "coffee3"
action = "coffeeAction"
broker = "tcp://192.168.100.11:1883"
```


# API-ROUTES

## /configuration

Path /configuration can be used to modify the devices.toml configuration.
Existing configuration can be fetched with GET-request while new configuration can be swapped with POST-request.

* POST (example curl)
```
curl -X POST http://localhost:8080/configuration \
     -H "Content-Type: application/json" \
     -d '{
  "general": {
     "interface": "wlan0"
  },
  "devices": {
    "coffee4": {
      "id": "coffee4",
      "action": "coffeeAction",
      "broker": "tcp://192.168.100.12:1883"
    },
    "coffee5": {
      "id": "coffee5",
      "action": "coffeeAction",
      "broker": "tcp://192.168.100.12:1883"
    },
    "coffee6": {
      "id": "coffee6",
      "action": "coffeeAction",
      "broker": "tcp://192.168.100.14:1883"
    }
  }
}'



```
* GET (example curl)

``` 
curl localhost:8080/configuration
```

