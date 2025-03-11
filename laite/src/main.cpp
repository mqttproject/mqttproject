#include <iostream>
#include "mqtt_client.h"
int main()
{
    mqtt_client mqtt_client1("tcp://localhost:1883", "client1");
    mqtt_client mqtt_client2("tcp://localhost:1883", "client2");

    mqtt_client1.connect();
    mqtt_client2.connect();

    mqtt_client1.subscribe("testi/topic", 1);
    mqtt_client2.publish("testi/topic", "juu", 1);

    mqtt_client1.wait_for_messages();
    mqtt_client2.wait_for_messages();

    std::cout << "Not blocking " << std::endl;

    mqtt_client1.disconnect();
    mqtt_client2.disconnect();

    return 0;
}