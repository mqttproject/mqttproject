#include <iostream>
#include <string>
#include "mqtt/async_client.h"

class mqtt_client
{
public:
    mqtt_client(std::string server_uri, std::string client_id);
    ~mqtt_client();
    bool connect();
    bool publish(std::string topic, std::string payload, int qos = 1);
    bool subscribe(std::string topic, int qos);
    void disconnect();
    bool is_connected();
    void wait_for_messages();

private:
    mqtt::async_client client;
    bool connected;
    std::string client_id;
};
