package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mig-elgt/chatws"
	"github.com/sirupsen/logrus"
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
		logrus.Errorf("could not find JWT")
		http.Error(w, "missing jwt code", http.StatusBadRequest)
		return
	}
	topics, err := h.convertTopics(r.URL.Query().Get("topics"))
	if err != nil {
		logrus.Errorf("could not convert topics: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload, err := h.auth.Authenticate(jwt)
	if err != nil {
		logrus.Errorf("could not authenticate client: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !h.matchTopics(payload.Topics, topics) {
		logrus.Errorf("could not match client topics: %v", topics)
		http.Error(w, "unauthorized to subscribe topics", http.StatusUnauthorized)
		return
	}
	topics["chat"] = []string{payload.ClientID}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("could not upgrade HTTP connection to WebSocket: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Info("client connected: ", payload.ClientID)

	done := make(chan error, 2)
	stop := make(chan struct{})

	go func() {
		done <- h.subscriber(conn, payload.ClientID, topics, stop)
	}()

	go func() {
		done <- h.reader(conn, stop)
	}()

	var stopped bool
	for i := 0; i < 1; i++ {
		if err := <-done; err != nil {
			logrus.Errorf("error: %v", err)
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
	if topicsStr == "" {
		return topics, nil
	}
	kinds := strings.Split(topicsStr, "|")
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

func (h handler) matchTopics(source, subscribe map[string][]string) bool {
	for sub, topics := range subscribe {
		if _, ok := source[sub]; !ok {
			return false
		}
		for _, topic := range topics {
			var found bool
			for _, st := range source[sub] {
				if topic == st {
					found = true
					break
				}
			}
			if !found {
				return false
			}

		}
	}
	return true
}
