const aedes = require('aedes')();
const net = require('net');

const PORT = 1883;

const server = net.createServer(aedes.handle);

server.listen(PORT, function () {
    console.log(`MQTT broker started on port ${PORT}`);
});


aedes.on('client', function (client) {
    console.log(`Client connected: ${client.id}`);
});


aedes.on('clientDisconnect', function (client) {
    console.log(`Client disconnected: ${client.id}`);
});


aedes.on('publish', function (packet, client) {
    if (client) {
        console.log(`Message from ${client.id}: ${packet.payload.toString()}`);
    }
});