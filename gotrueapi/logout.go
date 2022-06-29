package gotrueapi

import (
	"net/http"

	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

func Logout(host string, headers map[string]string) (*http.Request, error) {
	return reqbuilder.New().
		Method("POST").
		Headers(headers).
		Host(host).
		Path("/logout").
		Build()
}
