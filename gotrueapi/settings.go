package gotrueapi

import (
	"net/http"

	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

func Settings(host string) (*http.Request, error) {
	return reqbuilder.New().
		Host(host).
		Path("/settings").
		Build()
}
