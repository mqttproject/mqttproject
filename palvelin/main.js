const aedes = require('aedes')();
const net = require('net');
const { setTimeout } = require('timers');

const PORT = 1883;

const server = net.createServer(aedes.handle);

server.listen(PORT, function () {
    console.log(`MQTT broker started on port ${PORT}`);
});

aedes.on('client', function (client) {
    console.log(`\nClient connected: ${client.id}`);
    const message = "Moi!";
    const topic = `devices/${client.id}/message`;
    sendToClient(client, topic, message);
});



aedes.on('clientDisconnect', function (client) {
    console.log(`Client disconnected: ${client.id}`);
});


aedes.on('publish', function (packet, client) {
    if (client) {
        console.log(`Message from ${client.id}: ${packet.payload.toString()}`);
    }
});


function sendToClient(client, topic, message) {
    const packet = {
        topic: topic,
        payload: message,
        qos: 0, 
        retain: true
    };

    aedes.publish(packet, function (err) {
        if (err) {
            console.log(`Error sending message to client ${client.id}: ${err}`);
        } else {
            console.log(`Message sent to client ${client.id}: ${message}`);
        }
    });
}