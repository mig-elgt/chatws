package handler

import (
	"io"
	"io/ioutil"

	"github.com/gorilla/websocket"
)

func (h handler) subscriber(conn *websocket.Conn, clientID string, topics map[string][]string, stop chan struct{}) error {
	return h.broker.Subscribe(topics, clientID, stop, func(msg io.Reader) error {
		data, err := ioutil.ReadAll(msg)
		if err != nil {
			return err
		}
		if err := conn.WriteMessage(1, data); err != nil {
			return err
		}
		return nil
	})
}
