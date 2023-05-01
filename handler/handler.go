package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mig-elgt/chatws"
)

type handler struct {
	auth   chatws.AuthService
	broker chatws.MessageBroker
}

func New(auth chatws.AuthService, broker chatws.MessageBroker) http.Handler {
	r := mux.NewRouter()
	h := handler{auth, broker}
	r.HandleFunc("/ws", h.wsHandler)
	return r
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (h handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	conn.WriteMessage(1, []byte("hi client!"))
	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			// Handle close error
			fmt.Printf("%T, %v", err, err)
			return
		}
		fmt.Printf("got message: %v\n", string(p))
		if err := conn.WriteMessage(msgType, p); err != nil {
			log.Println(err)
		}
	}
}
