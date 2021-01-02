# UDP Websocket

This is a utility program to send UDP datagrams to websockets as binary packets.

## Configuration

The program is configured using enviremental variables as the following table.

| Name | Default | Description |
|------|---------|------|
| `HTTP_LISTEN` | `:9090` | the HTTP Listen port for incoming connections |
| `HTTP_PATH` | `/ws` | The path for incomming connections |
| `UDP_LISTEN` | `:9000` | The UDP port to listen for incoming connections |

## Docker
A docker container can be built using the dockerfile provided in the repository.

Simply run the following:
```bash
docker build -t udp-websocket .
```