FROM golang:1.15.3

ENV GOPRIVATE="*"
RUN apt update && apt -y install apt-transport-https && \
    apt -y install sqlite libsqlite-dev build-essential libssl-dev openssl && \
    mkdir /src

COPY main.go go.mod go.sum /src/
WORKDIR /src

RUN go mod verify && go build -o app && \
    mv app /usr/bin/udp-websocket
CMD ["/usr/bin/udp-websocket"]