cmake_minimum_required(VERSION 3.0)

project(mqtt)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include_directories(headers)
file(GLOB SOURCES "src/*.cpp")
file(COPY resources DESTINATION ${CMAKE_BINARY_DIR})

add_executable(mqtt ${SOURCES})
target_link_libraries(mqtt /usr/local/lib/libpaho-mqttpp3.so /usr/local/lib/libpaho-mqtt3c.so)


