package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleChatConnection(conn *websocket.Conn) {
	for {
		// Чтение сообщения от клиента
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Connection closed by client")
			}
			break
		}

		// Логирование полученного сообщения
		log.Printf("Received: %s", message)
		t := time.Now()

		messageText := "[" + t.Format("2006-01-02 15:04:05") + "] This is test message"
		messageType, message = 1, []byte(messageText)

		log.Printf("Sending: %s", message)

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Error while writing message: %v", err)
			// Check for closed connection error
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Connection closed by client")
			}
			break
		}
		time.Sleep(2 * time.Second)

	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// обновление соединения до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error updateding connection: %v", err)
		return
	}
	defer conn.Close()

	go HandleChatConnection(conn)
}
