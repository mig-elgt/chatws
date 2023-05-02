package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/mig-elgt/chatws"
	"github.com/mig-elgt/chatws/mocks"
)

func TestWebSocketHandler(t *testing.T) {
	type args struct {
		authenticateFnMock func(jwt string) (*chatws.TokenPayload, error)
		clientQuery        string
	}
	testCases := map[string]struct {
		args           args
		wantStatusCode int
	}{
		"missing jwt param": {
			args: args{
				clientQuery: "",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		"missing topics": {
			args: args{
				clientQuery: "jwt=header.payload.signature",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		"topics bad format": {
			args: args{
				clientQuery: "jwt=header.payload.signature&topics=logs",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		"auth service not available": {
			args: args{
				authenticateFnMock: func(jwt string) (*chatws.TokenPayload, error) {
					return nil, chatws.ErrAuthServiceNotAvilable
				},
				clientQuery: "jwt=header.payload.signature&topics=logs:foo,bar",
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		"client topics not authorization": {
			args: args{
				authenticateFnMock: func(jwt string) (*chatws.TokenPayload, error) {
					return &chatws.TokenPayload{
						Topics: map[string][]string{
							"logs": {"panic"},
						},
					}, nil
				},
				clientQuery: "jwt=header.payload.signature&topics=logs:foo,bar|sensors:gps",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := handler{
				auth: &mocks.AuthServiceMock{
					AuthenticateFn: tc.args.authenticateFnMock,
				},
			}
			svr := httptest.NewServer(http.HandlerFunc(h.wsHandler))
			u := "ws" + strings.TrimPrefix(svr.URL, "http")
			ws, resp, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%v?%v", u, tc.args.clientQuery), nil)
			if err != nil {
				if got, want := resp.StatusCode, tc.wantStatusCode; got != want {
					t.Fatalf("got %v; want %v status code", got, want)
				}
			}
			if err == nil {
				ws.Close()
				svr.Close()
			}
		})
	}
}
