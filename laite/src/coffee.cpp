#include "coffee.hpp"


Coffee::Coffee(const std::string &server_uri, const std::string &client_id) 
: Device(server_uri, client_id)
{
    start_brewing();
}


void Coffee::start_brewing()
{
    printf("Aloitetaan kahvin keitto\n");
    std::this_thread::sleep_for(std::chrono::minutes(1));
    printf("kahvi keitetty\n");
}