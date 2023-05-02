package grpc

import (
	"context"

	"github.com/gogo/status"
	"github.com/mig-elgt/chatws/auth/proto/auth"
	"github.com/mig-elgt/jwt"
	"google.golang.org/grpc/codes"
)

type handler struct {
	jwtsvc jwt.TokenCreateValidator
}

var users = map[string]map[string][]string{
	"foo": {
		"logs": {"panic", "errors"},
	},
	"bar": {
		"logs": {"panic"},
	},
	"zaas": {
		"logs": {"info"},
	},
}

func (h handler) login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if _, found := users[req.Username]; !found {
		return nil, status.Error(codes.Unauthenticated, "could not find users")
	}
	type token struct {
		ClientID string              `json:"client_id"`
		Topics   map[string][]string `json:"topics"`
	}
	t := &token{ClientID: req.Username, Topics: users[req.Username]}
	code, err := h.jwtsvc.Create(t)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not create jwt: %v", err)
	}
	return &auth.LoginResponse{JWT: code}, nil
}

func (h handler) authenticate(ctx context.Context, req *auth.AuthenticateRequest) (*auth.AuthenticateResponse, error) {
	panic("not impl")
}
