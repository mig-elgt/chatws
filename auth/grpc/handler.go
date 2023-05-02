package grpc

import (
	"context"

	"github.com/mig-elgt/chatws/auth/proto/auth"
	"github.com/mig-elgt/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	jwtsvc jwt.TokenCreateValidator
}

var users = map[string]map[string][]string{
	"foo": {
		"logs":    {"panic", "errors"},
		"devices": {"battery", "gps"},
	},
	"bar": {
		"logs": {"panic"},
	},
	"zaas": {
		"logs": {"info"},
	},
}

type token struct {
	ClientID string              `json:"client_id"`
	Topics   map[string][]string `json:"topics"`
}

func (h handler) login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if _, found := users[req.Username]; !found {
		return nil, status.Error(codes.Unauthenticated, "could not find users")
	}
	t := &token{ClientID: req.Username, Topics: users[req.Username]}
	code, err := h.jwtsvc.Create(t)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create jwt: %v", err)
	}
	return &auth.LoginResponse{JWT: code}, nil
}

func (h handler) authenticate(ctx context.Context, req *auth.AuthenticateRequest) (*auth.AuthenticateResponse, error) {
	data, err := h.jwtsvc.Validate(req.JWT)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "could not validate jwt: %v", err)
	}
	payload := data.(map[string]interface{})
	clientID, topics := payload["client_id"].(string), payload["topics"].(map[string]interface{})

	var clientTopics []*auth.Topics

	for topic, subTopics := range topics {
		t := []string{topic}
		for _, sub := range subTopics.([]interface{}) {
			t = append(t, sub.(string))
		}
		clientTopics = append(clientTopics, &auth.Topics{Name: t})
	}
	return &auth.AuthenticateResponse{ClientID: clientID, Topics: clientTopics}, nil
}
