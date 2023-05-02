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

// ws://localhost:8080/ws?jwt=header.payload.signature&topics=logs:foo,bar;
func (h handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	jwt := r.URL.Query().Get("jwt")
	if jwt == "" {
		http.Error(w, "missing jwt code", http.StatusBadRequest)
		return
	}
	topics := r.URL.Query().Get("topics")
	if topics == "" {
		http.Error(w, "missing topics", http.StatusBadRequest)
		return
	}
}
