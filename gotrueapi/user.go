package gotrueapi

import (
	"net/http"

	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

func GetUser(host, accessToken string) (*http.Request, error) {
	return reqbuilder.New().
		Method("GET").
		Headers("Authorization", "Bearer "+accessToken).
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

func PutUser(host, accessToken string, params *PutUserParams) (*http.Request, error) {
	return reqbuilder.New().
		Method("PUT").
		Host(host).
		Path("/user").
		Headers("Authorization", "Bearer "+accessToken).
		Body(params).
		Build()
}
