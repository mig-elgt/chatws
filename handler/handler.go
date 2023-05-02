package handler

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

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
	topicsToSub, err := h.convertTopics(topics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload, err := h.auth.Authenticate(jwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !reflect.DeepEqual(topicsToSub, payload.Topics) {
		http.Error(w, "unauthorized to subscribe topics", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	done := make(chan error, 2)
	stop := make(chan struct{})

	go func() {
		done <- h.subscriber(conn, payload.ClientID, topicsToSub, stop)
	}()

	go func() {
		done <- h.reader(conn, stop)
	}()

	var stopped bool
	for i := 0; i < 1; i++ {
		if err := <-done; err != nil {
			fmt.Printf("error: %v\n", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}

// logs:foo,bar|sensors:a,b
func (h handler) convertTopics(topicsStr string) (map[string][]string, error) {
	topics := map[string][]string{}
	kinds := strings.Split(topicsStr, "|")
	if len(kinds) == 0 {
		return nil, chatws.ErrClientTopicsBadFormat
	}
	for _, kind := range kinds {
		t := strings.Split(kind, ":")
		if len(t) != 2 {
			return nil, chatws.ErrClientTopicsBadFormat
		}
		topic := t[0]
		subTopics := strings.Split(t[1], ",")
		if len(subTopics) == 0 {
			return nil, chatws.ErrClientTopicsBadFormat
		}
		topics[topic] = subTopics
	}
	return topics, nil
}
