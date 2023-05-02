package handler

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Kind      string `json:"kind"`
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
}

func (h handler) reader(conn *websocket.Conn, stop chan struct{}) error {
	go func() {
		<-stop
		conn.Close()
	}()
	var message ChatMessage
	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if err := json.Unmarshal(p, &message); err != nil {
			if err := conn.WriteMessage(msgType, []byte("Invalid chat message format")); err != nil {
				log.Println(err)
			}
			continue
		}
		if message.Kind != "chat" {
			if err := conn.WriteMessage(msgType, []byte("Invalid chat topic")); err != nil {
				log.Println(err)
			}
			continue
		}
		if err := h.broker.Publish(
			message.Kind,
			message.Recipient,
			strings.NewReader(message.Message)); err != nil {
			if err := conn.WriteMessage(msgType, []byte("Could not send message")); err != nil {
				log.Println(err)
			}
			return err
		}
	}
}
