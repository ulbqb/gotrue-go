package gotrueapi

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ulbqb/gotrue-go/internal/reqbuilder"
)

type UpdateUserByIdParams struct {
	AppMetadata map[string]interface{} `json:"app_metadata,omitempty"`
}

func UpdateUserById(host string, headers map[string]string, uid uuid.UUID, params *UpdateUserByIdParams) (*http.Request, error) {
	return reqbuilder.New().
		Method("PUT").
		Host(host).
		Path("/admin/users/" + uid.String()).
		Headers(headers).
		Body(params).
		Build()
}
