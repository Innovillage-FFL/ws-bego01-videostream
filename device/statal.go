package device

import (
	"encoding/json"
	"log"
	"net/http"

	"begows01/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type StatusRequest struct {
	IDProduct       int64 `json:"id_product"`
	StatusPerangkat int   `json:"status_perangkat"`
}

type WSResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func UpdateStatusWS(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	db := utils.ConnectDB()

	for {

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var req StatusRequest
		err = json.Unmarshal(msg, &req)
		if err != nil {
			conn.WriteJSON(WSResponse{
				Success: false,
				Message: "format JSON salah",
			})
			continue
		}

		_, err = db.Exec(
			`INSERT INTO alat.status_perangkat (status_perangkat, id_product)
			 VALUES ($1,$2)`,
			req.StatusPerangkat,
			req.IDProduct,
		)

		if err != nil {
			log.Println("DB error:", err)
			conn.WriteJSON(WSResponse{
				Success: false,
				Message: "gagal update status",
			})
			continue
		}

		conn.WriteJSON(WSResponse{
			Success: true,
			Message: "status perangkat diperbarui",
		})
	}
}
