package handlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func (h *messageHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Поднимаем соединение
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// По умолчанию Upgrader не разрешает запросы из других Origin
		// Если нужно разрешить – настройки CORS здесь.
		CheckOrigin: func(r *http.Request) bool {
			// Разрешить все источники или написать свою логику
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v\n", err)
		return
	}

	// Регистрируем соединение в нашем хабе
	h.Hub.RegisterConnection(userID, conn)
	defer func() {
		// При выходе (закрытии хендлера) – отрегестрируем
		h.Hub.UnregisterConnection(userID)
		conn.Close()
	}()

	// Запускаем цикл чтения сообщений от клиента (если нужно)
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			// Когда клиент закрывает соединение, err != nil
			log.Printf("Error reading json: %v\n", err)
			break
		}

		// Обрабатываем входящее сообщение
		// (Если нужно реализовать двухсторонний обмен)
		log.Printf("Received message from user %d: %+v\n", userID, msg)

		// Можете проверить тип сообщения (например, "ping", "typing" и т.д.)
		// и отправить нужным пользователям.
	}
}
