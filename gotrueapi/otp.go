package gotrueapi

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

type OTPParams struct {
	Security Security `json:"gotrue_meta_security,omitempty"`

	Email      string `json:"email"`
	Phone      string `json:"phone"`
	CreateUser bool   `json:"create_user"`

	RedirectTo string `json:"-"`
}

func OTP(host string, params *OTPParams) (*http.Request, error) {
	if len(params.Email) > 0 && len(params.Phone) > 0 {
		return nil, errors.New("api: email and phone were provided at the same time")
	}
	if len(params.Email) == 0 && len(params.Phone) == 0 {
		return nil, errors.New("api: email or phone should be provided")
	}

	return reqbuilder.New().
		Method("POST").
		Host(host).
		Path("/otp").
		Queries("redirect_to", params.RedirectTo).
		Body(params).
		Build()
}
