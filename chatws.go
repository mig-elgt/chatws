package chatws

import (
	"errors"
	"io"
)

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
	Subscribe(topics map[string][]string, clientID string, stop chan struct{}, callback func(msg io.Reader) error) error
	Publish(topic, subTopic string, msg io.Reader) error
	Close() error
}

var (
	ErrAuthServiceNotAvilable = errors.New("auth service is not available")
	ErrInvalidJWT             = errors.New("JWT is invalid or has expired")
	ErrClientTopicsBadFormat  = errors.New("client topics bad format")
)
