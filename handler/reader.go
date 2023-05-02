package handler

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Kind      string `json:"kind"`
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
}

func (h handler) reader(conn *websocket.Conn) {
	var message ChatMessage
	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if err := json.Unmarshal(p, &message); err != nil {
			if err := conn.WriteMessage(msgType, []byte("Invalid chat message format")); err != nil {
				log.Println(err)
			}
		}
	}
}
