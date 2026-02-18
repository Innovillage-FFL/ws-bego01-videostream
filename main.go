package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var mutex sync.Mutex

func main() {
	http.HandleFunc("/ws/esp32", handleESP32)
	http.HandleFunc("/ws/viewer", handleViewer)

	log.Println("Server running :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleESP32(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("ESP32 connected")

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("ESP32 disconnected")
			return
		}
		broadcast(msgType, data)
	}
}

func handleViewer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	log.Println("Viewer connected")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			conn.Close()
			log.Println("Viewer disconnected")
			return
		}
	}
}

func broadcast(msgType int, data []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	for c := range clients {
		err := c.WriteMessage(msgType, data)
		if err != nil {
			c.Close()
			delete(clients, c)
		}
	}
}
