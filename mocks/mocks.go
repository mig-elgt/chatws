package mocks

import (
	"io"

	"github.com/mig-elgt/chatws"
)

type AuthServiceMock struct {
	AuthenticateFn func(jwt string) (*chatws.TokenPayload, error)
}

func (svc *AuthServiceMock) Authenticate(jwt string) (*chatws.TokenPayload, error) {
	return svc.AuthenticateFn(jwt)
}

type MessageBrokerMock struct {
	SubscribeFn func(topics map[string][]string, clientID string, stop chan struct{}, callback func(msg io.Reader) error) error
	PublishFn   func(topic string, subTopic string, msg io.Reader) error
}

func (m *MessageBrokerMock) Subscribe(topics map[string][]string, clientID string, stop chan struct{}, callback func(msg io.Reader) error) error {
	return m.SubscribeFn(topics, clientID, stop, callback)
}

func (m *MessageBrokerMock) Publish(topic string, subTopic string, msg io.Reader) error {
	return m.PublishFn(topic, subTopic, msg)
}

func (m *MessageBrokerMock) Close() error {
	return nil
}
