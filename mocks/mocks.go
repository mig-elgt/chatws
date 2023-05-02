package mocks

import (
	"context"
	"io"

	"github.com/mig-elgt/chatws"
	pb "github.com/mig-elgt/chatws/auth/proto/auth"
	"google.golang.org/grpc"
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

type AuthServiceClientMock struct {
	AuthenticateFn func(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error)
}

func (a *AuthServiceClientMock) Login(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (a *AuthServiceClientMock) Authenticate(ctx context.Context, in *pb.AuthenticateRequest, opts ...grpc.CallOption) (*pb.AuthenticateResponse, error) {
	return a.AuthenticateFn(ctx, in)
}
