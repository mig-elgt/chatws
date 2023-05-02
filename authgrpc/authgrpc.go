package authgrpc

import (
	"context"

	"github.com/mig-elgt/chatws"
	pb "github.com/mig-elgt/chatws/auth/proto/auth"
)

type authgrpc struct {
	client pb.AuthServiceClient
}

func (a *authgrpc) Authenticate(jwt string) (*chatws.TokenPayload, error) {
	_, err := a.client.Authenticate(context.Background(), &pb.AuthenticateRequest{
		JWT: jwt,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
