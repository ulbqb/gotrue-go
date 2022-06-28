package gotrueapi

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

type RecoverParams struct {
	Security Security `json:"gotrue_meta_security,omitempty"`

	Email string `json:"email"`

	RedirectTo string `json:"-"`
}

func Recover(host string, headers map[string]string, params *RecoverParams) (*http.Request, error) {
	if len(params.Email) == 0 {
		return nil, errors.New("api: email should be provided")
	}

	return reqbuilder.New().
		Method("POST").
		Headers(headers).
		Host(host).
		Path("/recover").
		Body(params).
		Queries("redirect_to", params.RedirectTo).
		Build()
}
