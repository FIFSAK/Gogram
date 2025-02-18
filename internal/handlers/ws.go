package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func (h *messageHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v\n", err)
		return
	}

	h.Hub.RegisterConnection(userID, conn)
	defer func() {
		h.Hub.UnregisterConnection(userID)
		conn.Close()
	}()

	for {
		var msg interface{}
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading json: %v\n", err)
			break
		}
		fmt.Println(msg)

	}
}
