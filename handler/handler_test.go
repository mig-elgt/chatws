package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/mig-elgt/chatws"
	"github.com/mig-elgt/chatws/mocks"
)

func TestWebSocketHandler(t *testing.T) {
	type args struct {
		authenticateFnMock func(jwt string) (*chatws.TokenPayload, error)
		publishFnMock      func(topic, subTopic string, msg io.Reader) error
		clientQuery        string
		clientMessage      string
	}
	testCases := map[string]struct {
		args           args
		wantStatusCode int
		wantErrMessage []byte
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
				publishFnMock: func(topic, subTopic string, msg io.Reader) error {
					return nil
				},
				clientQuery: "jwt=header.payload.signature&topics=logs:foo,bar|sensors:gps",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		"send bad message format": {
			args: args{
				authenticateFnMock: func(jwt string) (*chatws.TokenPayload, error) {
					return &chatws.TokenPayload{
						Topics: map[string][]string{
							"logs": {"panic"},
						},
					}, nil
				},
				publishFnMock: func(topic, subTopic string, msg io.Reader) error {
					return nil
				},
				clientQuery:   "jwt=header.payload.signature&topics=logs:panic",
				clientMessage: "foobar",
			},
			wantErrMessage: []byte("Invalid chat message format"),
		},
		"invalid chat topic": {
			args: args{
				authenticateFnMock: func(jwt string) (*chatws.TokenPayload, error) {
					return &chatws.TokenPayload{
						Topics: map[string][]string{
							"logs": {"panic"},
						},
					}, nil
				},
				publishFnMock: func(topic, subTopic string, msg io.Reader) error {
					return nil
				},
				clientQuery:   "jwt=header.payload.signature&topics=logs:panic",
				clientMessage: `{"kind":"foo","message":"foobar","recipient":"clientB"}`,
			},
			wantErrMessage: []byte("Invalid chat topic"),
		},
		"failed to send a message": {
			args: args{
				authenticateFnMock: func(jwt string) (*chatws.TokenPayload, error) {
					return &chatws.TokenPayload{
						Topics: map[string][]string{
							"logs": {"panic"},
						},
					}, nil
				},
				publishFnMock: func(topic, subTopic string, msg io.Reader) error {
					return errors.New("could not publish message")
				},
				clientQuery:   "jwt=header.payload.signature&topics=logs:panic",
				clientMessage: `{"kind":"chat","message":"foobar","recipient":"clientB"}`,
			},
			wantErrMessage: []byte("Could not send message"),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			h := handler{
				auth: &mocks.AuthServiceMock{
					AuthenticateFn: tc.args.authenticateFnMock,
				},
				broker: &mocks.MessageBrokerMock{
					PublishFn: tc.args.publishFnMock,
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
				if err := ws.WriteMessage(websocket.TextMessage, []byte(tc.args.clientMessage)); err != nil {
					t.Fatalf("%v", err)
				}
				_, p, err := ws.ReadMessage()
				if err != nil {
					t.Fatalf("%v", err)
				}
				if got, want := p, tc.wantErrMessage; !reflect.DeepEqual(got, want) {
					t.Fatalf("got %v; want %v", string(got), string(want))
				}
				ws.Close()
				svr.Close()
			}
		})
	}
}
