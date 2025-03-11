project('mqtt','cpp')

inc = include_directories('headers')
src = files(
    'src/main.cpp',
    'src/mqtt_client.cpp',
)
install_emptydir('resources')

paho_mqttpp = meson.get_compiler('cpp').find_library('paho-mqttpp3', dirs: ['/usr/local/lib'], required: true)
paho_mqttc = meson.get_compiler('cpp').find_library('paho-mqtt3c', dirs: ['/usr/local/lib'], required: true)

executable('mqtt', src,
           include_directories: [inc],
           dependencies: [paho_mqttpp, paho_mqttc]
)