package main

import (
	"log"
	"net/http"

	"begows01/device"
)

func main() {
	go device.StartDeviceOfflineJob()

	http.HandleFunc("/stream/", device.StreamDevice)
	http.HandleFunc("/watch/", device.WatchDevice)

	http.HandleFunc("/ws/status", device.UpdateStatusWS)

	log.Println("Server running :8081")
	http.ListenAndServe(":8081", nil)
}
