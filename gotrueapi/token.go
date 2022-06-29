package gotrueapi

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

type TokenWithPasswordGrantParams struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func TokenWithPasswordGrant(host string, headers map[string]string, params *TokenWithPasswordGrantParams) (*http.Request, error) {
	if len(params.Email) > 0 && len(params.Phone) > 0 {
		return nil, errors.New("api: email and phone were provided at the same time")
	}
	if len(params.Email) == 0 && len(params.Phone) == 0 {
		return nil, errors.New("api: email or phone should be provided")
	}

	return reqbuilder.New().
		Method("POST").
		Headers(headers).
		Host(host).
		Path("/token").
		Queries("grant_type", "password").
		Body(params).
		Build()
}

type TokenWithRefreshTokenGrantParams struct {
	RefreshToken string `json:"refresh_token"`
}

func TokenWithRefreshTokenGrant(host string, headers map[string]string, params *TokenWithRefreshTokenGrantParams) (*http.Request, error) {
	return reqbuilder.New().
		Method("POST").
		Headers(headers).
		Host(host).
		Path("/token").
		Queries("grant_type", "refresh_token").
		Body(params).
		Build()
}

type TokenWithIDTokenGrantParams struct {
	IdToken  string `json:"id_token"`
	Nonce    string `json:"nonce"`
	Provider string `json:"provider"`

	RedirectTo string `json:"-"`
}

func TokenWithIDTokenGrant(host string, headers map[string]string, params *TokenWithIDTokenGrantParams) (*http.Request, error) {
	return reqbuilder.New().
		Method("POST").
		Headers(headers).
		Host(host).
		Path("/token").
		Queries("grant_type", "refresh_token").
		Queries("redirect_to", params.RedirectTo).
		Body(params).
		Build()
}
