cmake_minimum_required(VERSION 3.0)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
project(mqtt)

set(CMAKE_PREFIX_PATH "/usr/local/lib/cmake/PahoMqttCpp" ${CMAKE_PREFIX_PATH})
find_package(PahoMqttCpp REQUIRED)


include_directories(headers)
file(GLOB SOURCES "src/*.cpp")
file(COPY resources DESTINATION ${CMAKE_BINARY_DIR})

add_executable(mqtt ${SOURCES})


