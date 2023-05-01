package handler

import (
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
	// TODO: Validate client connection
	// TODO: Validate client topics
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	// TODO: Run Subscriber
	// TODO: Run Reader
	// TODO: Error Handler for Subs and Reader to close stop channel
	h.reader(conn)
}
