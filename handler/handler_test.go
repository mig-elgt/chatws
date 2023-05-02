package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/mig-elgt/chatws"
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
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := handler{}
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
