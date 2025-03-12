#include "device.hpp"

Device::Device(const std::string &server_uri, const std::string &client_id)
    : client(server_uri, client_id)
{
    client.connect();
    client.subscribe("testi/topic", 1);
    client.wait_for_messages();
}

Device::~Device()
{
    client.disconnect();
}
