package device

import (
	"log"
	"time"

	"begows01/utils"
)

func StartDeviceOfflineJob() {

	db := utils.ConnectDB()

	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	log.Println("Device offline job started (interval 3 menit)")

	for {
		<-ticker.C

		log.Println("Running offline job...")

		_, err := db.Exec(`
			INSERT INTO alat.status_perangkat (status_perangkat, id_product)
			SELECT 0, id_product
			FROM alat.produk
		`)

		if err != nil {
			log.Println("Offline job error:", err)
			continue
		}

		log.Println("Semua device diset offline")
	}
}
