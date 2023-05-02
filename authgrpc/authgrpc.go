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
	resp, err := a.client.Authenticate(context.Background(), &pb.AuthenticateRequest{
		JWT: jwt,
	})
	if err != nil {
		return nil, err
	}
	subscribedTopics := map[string][]string{}

	for _, topics := range resp.Topics {
		topic, subTopics := topics.Name[0], topics.Name[1:]
		subscribedTopics[topic] = subTopics
	}

	return &chatws.TokenPayload{ClientID: resp.ClientID, Topics: subscribedTopics}, nil
}
