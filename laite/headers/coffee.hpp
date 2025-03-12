#ifndef COFFEE_HPP
#define COFFEE_HPP

#include "device.hpp"


class Coffee : public Device
{
public:
    Coffee(const std::string& server_uri, const std::string& client_id);
    void turn_on();
    void turn_off();
private:
    void start_brewing();
};

#endif