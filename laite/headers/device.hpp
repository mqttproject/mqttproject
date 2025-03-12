#ifndef DEVICE_HPP
#define DEVICE_HPP

#include "mqtt_client.hpp"
#include <chrono>

class Device
{
public:
    Device(const std::string& server_uri, const std::string& client_id);
    ~Device();
private:
    mqtt_client client;
    char name[256];
    char ip[256];
    char mac[256];

};

#endif