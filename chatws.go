package chatws

import "io"

// TokenPayload describes a JTW payload data for a client
type TokenPayload struct {
	ClientID string
	Topics   map[string][]string
}

// AuthService describes the behaviour to perform clients authentication
type AuthService interface {
	Authenticate(jwt string) (*TokenPayload, error)
}

// MessageBroker describes the behavior to perform Pub/Sub operations.
type MessageBroker interface {
	Subscribe(topics map[string][]string, clientID string, stop chan struct{}, callback func(msg io.Reader) error)
	Publish(topic, subTopic string, msg io.Reader) error
	Close() error
}
