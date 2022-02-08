package gotrueapi

import (
	"net/http"

	"go.lair.cx/gotrue-go/internal/reqbuilder"
)

func Settings(host string) (*http.Request, error) {
	return reqbuilder.New().
		Host(host).
		Path("/settings").
		Build()
}
