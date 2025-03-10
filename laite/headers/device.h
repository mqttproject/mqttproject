#include <iostream>
#include <stdint.h>
#include <string>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <net/if.h>
#include <cstring>
class device
{
public:
    device(std::string virtual_ip, std::string server_ip, uint32_t port, std::string interface);
    std::string get_ip() const;
    std::string get_server_ip() const;
    uint32_t get_port() const;
    void connect_to_network();
    void clean_virtual_ip();

private:
    std::string virtual_ip;
    uint32_t port;
    std::string interface;
    std::string server_ip;

    // Structures for creating packets to send through tcp
    struct sockaddr_in local_addr;
    struct sockaddr_in server_addr;

    // Socket file descriptor
    int sockfd;
};