package mocks

import "github.com/mig-elgt/chatws"

type AuthServiceMock struct {
	AuthenticateFn func(jwt string) (*chatws.TokenPayload, error)
}

func (svc *AuthServiceMock) Authenticate(jwt string) (*chatws.TokenPayload, error) {
	return svc.AuthenticateFn(jwt)
}
