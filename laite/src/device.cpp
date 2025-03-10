#include "device.h"

device::device(std::string virtual_ip, std::string server_ip, uint32_t port, std::string interface)
    : virtual_ip(virtual_ip), server_ip(server_ip), port(port), interface(interface) {}

std::string device::get_ip() const
{
    return this->virtual_ip;
}

std::string device::get_server_ip() const
{
    return this->server_ip;
}

uint32_t device::get_port() const
{
    return this->port;
}

void device::clean_virtual_ip()
{
    std::string cmd_del_ip = "sudo ip addr del " + virtual_ip + "/24 dev " + interface + " label " + interface + ":0";
    system(cmd_del_ip.c_str());
    std::cout << "Virtual IP removed: " << virtual_ip << "\n";
}

void device::connect_to_network()
{
    std::string cmd_ip = "sudo ip addr add " + virtual_ip + "/24 dev " + interface + " label " + interface + ":0";
    system(cmd_ip.c_str());

    std::cout << "Virtual IP assigned: " << virtual_ip << "\n";

    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd < 0)
    {
        std::cerr << "Socket creation failed.\n";
        clean_virtual_ip();
        return;
    }

    local_addr.sin_family = AF_INET;
    local_addr.sin_port = 0;
    inet_pton(AF_INET, virtual_ip.c_str(), &local_addr.sin_addr);

    if (bind(sockfd, (struct sockaddr *)&local_addr, sizeof(local_addr)) < 0)
    {
        std::cerr << "Binding the socket to the virtual IP failed.\n";
        close(sockfd);
        clean_virtual_ip();
        return;
    }

    std::cout << "Virtual IP bound to local socket " << std::endl;

    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(port);
    inet_pton(AF_INET, server_ip.c_str(), &server_addr.sin_addr);

    if (connect(sockfd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0)
    {
        std::cerr << "Connection to server failed.\n";
        close(sockfd);
        clean_virtual_ip();
        return;
    }

    std::cout << "Connected to server at " << server_ip << " on port " << port << "\n";
    close(sockfd);
    clean_virtual_ip();
}