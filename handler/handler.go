package handler

import (
	"net/http"

	"github.com/gorilla/mux"
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

func (h handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
