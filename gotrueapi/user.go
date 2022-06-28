package gotrueapi

import (
	"net/http"

	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

func GetUser(host string, headers map[string]string) (*http.Request, error) {
	return reqbuilder.New().
		Method("GET").
		Headers(headers).
		Host(host).
		Path("/user").
		Build()
}

type PutUserParams struct {
	Email    string                 `json:"email"`
	Password *string                `json:"password"`
	Data     map[string]interface{} `json:"data"`
	AppData  map[string]interface{} `json:"app_metadata,omitempty"`
	Phone    string                 `json:"phone"`
}

func PutUser(host string, headers map[string]string, params *PutUserParams) (*http.Request, error) {
	return reqbuilder.New().
		Method("PUT").
		Host(host).
		Path("/user").
		Headers(headers).
		Body(params).
		Build()
}
