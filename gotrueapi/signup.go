package gotrueapi

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

type SignUpParams struct {
	Security Security `json:"gotrue_meta_security,omitempty"`

	Email    string      `json:"email,omitempty"`
	Phone    string      `json:"phone,omitempty"`
	Password string      `json:"password,omitempty"`
	Data     interface{} `json:"data,omitempty"`

	RedirectTo string `json:"-"`
}

func SignUp(host string, headers map[string]string, params *SignUpParams) (*http.Request, error) {
	if len(params.Email) > 0 && len(params.Phone) > 0 {
		return nil, errors.New("api: email and phone were provided at the same time")
	}
	if len(params.Email) == 0 && len(params.Phone) == 0 {
		return nil, errors.New("api: email or phone should be provided")
	}
	if len(params.Password) == 0 {
		return nil, errors.New("api: password is required")
	}

	return reqbuilder.New().
		Method("POST").
		Headers(headers).
		Host(host).
		Path("/signup").
		Queries("redirect_to", params.RedirectTo).
		Body(params).
		Build()
}
