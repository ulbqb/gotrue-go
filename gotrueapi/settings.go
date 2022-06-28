package gotrueapi

import (
	"net/http"

	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

func Settings(host string, headers map[string]string) (*http.Request, error) {
	return reqbuilder.New().
		Method("GET").
		Headers(headers).
		Host(host).
		Path("/settings").
		Build()
}
