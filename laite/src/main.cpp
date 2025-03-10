#include <iostream>
#include "device.h"
int main(void)
{
    std::string server = "0.0.0.0";

    device device1("192.168.1.105", server, 8080, "enp14s0");
    device1.connect_to_network();

    device device2("192.168.1.106", server, 8080, "enp14s0");
    device2.connect_to_network();
    return 0;
}