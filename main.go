package main

import (
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

var packets chan []byte
//var connections []*websocket.Conn
var connections sync.Map

// configure handles configuration from env variables
func configure() {
	viper.SetDefault("http.listen", ":9090")
	viper.SetDefault("http.path", "/ws")
	viper.SetDefault("udp.listen", ":9000")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func main() {
	log.Println("Starting UDP->WebSocket forwarder...")
	configure()

	packets = make(chan []byte)
	go transfer()
	go listenUDP()

	log.Printf("Listening for websockets on %s%s", viper.GetString("http.listen"), viper.GetString("http.path"))
	http.HandleFunc(viper.GetString("http.path"), websocketHandler)

	log.Fatal(http.ListenAndServe(viper.GetString("http.listen"), nil))
}

// transfer reads from the incoming packets channel and sends
// the data to all registered websockets.
func transfer() {
	for {
		p := <-packets
		connections.Range(func(key, value interface{}) bool {
			ws := key.(*websocket.Conn)
			if err := ws.WriteMessage(websocket.BinaryMessage, p); err != nil {
				if err == websocket.ErrCloseSent {
					log.Println("Closed websocket")
				} else {
					log.Println("Error sending message to client:", err)
				}
			}
			return true
		})
	}
}

// websocketReader is a blocking call the reads from the websocket
// indefinitely and discards all data
func websocketReader(ws *websocket.Conn) {
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
	}
}

// websocketHandler is the HTTP handler for incoming connections.
// The connection is upgraded to a persistent websocket connection.
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 3000,
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading socket connection:", err)
		return
	}
	connections.Store(ws, true)
	defer connections.Delete(ws)
//	connections = append(connections, ws)
	websocketReader(ws)
}

// listenUDP is a blocking call that opens a UDP port and listens for
// incoming data packages. The incoming packages are put on the packets channel.
func listenUDP() {
	log.Println("Listening for UDP on", viper.GetString("udp.listen"))
	pc, err := net.ListenPacket("udp", viper.GetString("udp.listen"))
	if err != nil {
		log.Fatal("Could not listen for UDP:", err)
	}

	defer pc.Close()

	for {
		buf := make([]byte, 4096)
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		if n < 1 {
			continue
		}
		buf = buf[:n]
		packets <- buf
	}
}
