package handler

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func (h handler) reader(conn *websocket.Conn) {
	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			// Handle close error
			fmt.Printf("%T, %v", err, err)
			return
		}
		// parse p message
		// publish message to message broker
		// handler error
		fmt.Printf("got message: %v\n", string(p))
		if err := conn.WriteMessage(msgType, p); err != nil {
			log.Println(err)
		}
	}
}
