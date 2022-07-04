package gotrueapi

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("(%d) %s", e.Status, e.Message)
}

type Identity struct {
	ID           string                 `json:"id"`
	UserID       uuid.UUID              `json:"user_id"`
	IdentityData map[string]interface{} `json:"identity_data,omitempty"`
	Provider     string                 `json:"provider"`
	LastSignInAt *time.Time             `json:"last_sign_in_at,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type User struct {
	ID uuid.UUID `json:"id"`

	Aud              string     `json:"aud"`
	Role             string     `json:"role"`
	Email            string     `json:"email"`
	EmailConfirmedAt *time.Time `json:"email_confirmed_at,omitempty"`
	InvitedAt        *time.Time `json:"invited_at,omitempty"`

	Phone            string     `json:"phone"`
	PhoneConfirmedAt *time.Time `json:"phone_confirmed_at,omitempty"`

	ConfirmationSentAt *time.Time `json:"confirmation_sent_at,omitempty"`

	RecoverySentAt *time.Time `json:"recovery_sent_at,omitempty"`

	EmailChange       string     `json:"new_email,omitempty"`
	EmailChangeSentAt *time.Time `json:"email_change_sent_at,omitempty"`

	PhoneChange       string     `json:"new_phone,omitempty"`
	PhoneChangeSentAt *time.Time `json:"phone_change_sent_at,omitempty"`

	LastSignInAt *time.Time `json:"last_sign_in_at,omitempty"`

	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetaData map[string]interface{} `json:"user_metadata"`

	Identities []Identity `json:"identities" has_many:"identities"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	BannedUntil *time.Time `json:"banned_until,omitempty"`
}

type Session struct {
	Token        string `json:"access_token"`
	TokenType    string `json:"token_type"` // Bearer
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}
