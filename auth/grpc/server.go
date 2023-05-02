package grpc

import (
	"context"

	"github.com/mig-elgt/chatws/auth/proto/auth"
	pb "github.com/mig-elgt/chatws/auth/proto/auth"
	"github.com/mig-elgt/jwt"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedAuthServiceServer
	handler
}

func NewAPI(jwtsvc jwt.TokenCreateValidator) *grpc.Server {
	rootServer := grpc.NewServer()
	s := &grpcServer{
		handler: handler{jwtsvc},
	}
	pb.RegisterAuthServiceServer(rootServer, s)
	return rootServer
}

func (g *grpcServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	return g.login(ctx, req)
}

func (g *grpcServer) Authenticate(ctx context.Context, req *auth.AuthenticateRequest) (*auth.AuthenticateResponse, error) {
	return g.authenticate(ctx, req)
}
