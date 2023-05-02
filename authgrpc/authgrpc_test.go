package authgrpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/mig-elgt/chatws"
	pb "github.com/mig-elgt/chatws/auth/proto/auth"
	"github.com/mig-elgt/chatws/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthgRPC_Authenticate(t *testing.T) {
	type args struct {
		authenticateFnMock func(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error)
	}
	testCases := map[string]struct {
		args      args
		wantErr   bool
		expectErr error
		want      *chatws.TokenPayload
	}{
		"auth grpc server not available": {
			args: args{
				authenticateFnMock: func(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
					return nil, status.Error(codes.Internal, "service not available")
				},
			},
			wantErr:   true,
			expectErr: errors.New("rpc error: code = Internal desc = service not available"),
		},
		"base case": {
			args: args{
				authenticateFnMock: func(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
					return &pb.AuthenticateResponse{
						ClientID: "foobar",
						Topics: []*pb.Topics{
							{Name: []string{"logs", "panic", "error"}},
							{Name: []string{"sensors", "gps", "battery"}},
						},
					}, nil
				},
			},
			wantErr: false,
			want: &chatws.TokenPayload{
				ClientID: "foobar",
				Topics: map[string][]string{
					"logs":    {"panic", "error"},
					"sensors": {"gps", "battery"},
				},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			svc := authgrpc{
				client: &mocks.AuthServiceClientMock{
					AuthenticateFn: tc.args.authenticateFnMock,
				},
			}
			got, err := svc.Authenticate("token")
			if (err != nil) != tc.wantErr {
				t.Fatalf("got %v; want %v", err, tc.wantErr)
			}
			if tc.wantErr && !reflect.DeepEqual(err.Error(), tc.expectErr.Error()) {
				t.Fatalf("got error %v; want %v", err, tc.expectErr)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got user %v; want %v", got, tc.want)
			}
		})
	}
}
