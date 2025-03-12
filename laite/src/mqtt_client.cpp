#include "mqtt_client.hpp"

mqtt_client::mqtt_client(std::string server_uri, std::string client_id)
    : client(server_uri, client_id), connected(false) {}

mqtt_client::~mqtt_client()
{
    if (connected)
    {
        client.disconnect()->wait();
    }
}
bool mqtt_client::connect()
{
    if (connected)
    {
        std::cout << "Already connected to MQTT server.\n";
        return true;
    }

    try
    {
        client.connect()->wait();
        connected = true;
        std::cout << "Connected to MQTT server.\n";
        return true;
    }
    catch (const mqtt::exception &e)
    {
        std::cerr << "Error while connecting: " << e.what() << std::endl;
        return false;
    }
}

bool mqtt_client::publish(std::string topic, std::string payload, int qos)
{
    if (!connected)
    {
        std::cerr << "Not connected. Please connect first.\n";
        return false;
    }

    try
    {
        std::cout << "Publishing message: " << payload << " to topic: " << topic << "\n";
        mqtt::message_ptr message = mqtt::make_message(topic, payload);
        message->set_qos(qos);
        client.publish(message);
        return true;
    }
    catch (const mqtt::exception &e)
    {
        std::cerr << "Error while publishing: " << e.what() << std::endl;
        return false;
    }
}

bool mqtt_client::subscribe(std::string topic, int qos)
{
    if (!connected)
    {
        std::cerr << "Not connected. Please connect first.\n";
        return false;
    }

    try
    {
        std::cout << "Subscribing to topic: " << topic << "\n";
        client.subscribe(topic, qos);
        return true;
    }
    catch (const mqtt::exception &e)
    {
        std::cerr << "Error while subscribing: " << e.what() << std::endl;
        return false;
    }
}

void mqtt_client::disconnect()
{
    if (connected)
    {
        client.disconnect()->wait();
        connected = false;
        std::cout << "Disconnected from MQTT server.\n";
    }
}

bool mqtt_client::is_connected()
{
    return connected;
}
void mqtt_client::wait_for_messages()
{
    try
    {
        client.start_consuming();
    }
    catch (const mqtt::exception &e)
    {
        std::cerr << "Error while waiting for messages: " << e.what() << std::endl;
    }
}
