package device

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var upgraderStream = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var viewers = make(map[string]map[*websocket.Conn]bool)
var mutex sync.Mutex

func StreamDevice(w http.ResponseWriter, r *http.Request) {

	deviceID := strings.TrimPrefix(r.URL.Path, "/stream/")
	if deviceID == "" {
		http.Error(w, "device id kosong", 400)
		return
	}

	conn, err := upgraderStream.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("Device connected:", deviceID)

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("Device disconnected:", deviceID)
			return
		}

		broadcastToDevice(deviceID, msgType, data)
	}
}

func WatchDevice(w http.ResponseWriter, r *http.Request) {

	deviceID := strings.TrimPrefix(r.URL.Path, "/watch/")
	if deviceID == "" {
		http.Error(w, "device id kosong", 400)
		return
	}

	conn, err := upgraderStream.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	mutex.Lock()

	if viewers[deviceID] == nil {
		viewers[deviceID] = make(map[*websocket.Conn]bool)
	}

	viewers[deviceID][conn] = true
	mutex.Unlock()

	log.Println("Viewer connected to:", deviceID)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(viewers[deviceID], conn)
			mutex.Unlock()

			conn.Close()
			log.Println("Viewer disconnected:", deviceID)
			return
		}
	}
}

func broadcastToDevice(deviceID string, msgType int, data []byte) {

	mutex.Lock()
	defer mutex.Unlock()

	for c := range viewers[deviceID] {
		err := c.WriteMessage(msgType, data)
		if err != nil {
			c.Close()
			delete(viewers[deviceID], c)
		}
	}
}
