# CONFIGURATION

Configuration is done through laite/devices.toml.

The structure of the toml is expected to be the following.

```
#Header "general" will hold the general configuration
[general]
#id will be generated during runtime. Dont touch it.
id = "a7474cba-4403-4c97-8ba7-d561c9c0f983" 

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

The program will try to create multiple macvlan interfaces under the physical interface given through the interface configuration.


## Example configuration

```
[general]
id = "a7474cba-4403-4c97-8ba7-d561c9c0f983"
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
Existing configuration can be fetched with GET-request while new configuration can be swapped with POST-request. The new configuration will not get evaluated in runtime.

* POST (example curl)
```
curl -X POST http://localhost:8080/configuration \
     -H "Content-Type: application/json" \
     -d '{
  "general": {
      "id":"a7474cba-4403-4c97-8ba7-d561c9c0f983",
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
## /device/:id
 
Path /device/:id takes a device id as a parameter. Existing device with matching id can be fetched through a GET-request while a new device with unique id can be created with POST-request. This route can be used to add new devices during runtime.

* POST (Example curl)

``` 
curl -X POST http://localhost:8080/device/coffee2/off \
 -H "Content-Type: application/json" \
 -d '{
  "action":"coffeeAction","broker":"tcp://0.0.0.0:1883"
}'
``` 
Notice that new device id is not placed to the json body but to the path parameters.

* GET (example curl)
  
```
curl -X GET http://localhost:8080/device/coffee2
```


## /device/:id/on

Path /device/:id/on takes a device id as a parameter. Existing device with a matching id can be turned on with a post request to this path. This route can be used to turn on devices during runtime.

* POST (example curl)
```
curl -X POST http://localhost:8080/device/coffee2/on
```

## /device/:id/off

Path /device/:id/off takes a device id as a parameter. Existing device with a matching id can be turned off with a post request to this path. This route can be used to turn off devices during runtime.

* POST (example curl)
```
curl -X POST http://localhost:8080/device/coffee2/off
```

## /reboot

Path /reboot attempts to reboot the simulator.
This will cause all the devices on that particular simulator to lose connection to the broker momentarily.

* POST (example curl)
```
curl -X POST http://localhost:8080/reboot
```

## /devices 

POST-request to this path accepts a json of multiple devices. It will add every device on the received json to the devices list during runtime.

GET-request to this path will return a list of all the devices currently on the devices list. 

* POST (example curl)

``` 
curl -X POST http://localhost:8080/devices \
  -H "Content-Type: application/json" \
  -d '{
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
curl localhost:8080/devices
```



## /device/:id/delete 

This path can be used to delete a device from the runtime list of devices. The deleted device will disconnect from broker if its connected.

* POST (example curl)

```
curl -X POST http://localhost:8080/device/coffee2/delete
```

