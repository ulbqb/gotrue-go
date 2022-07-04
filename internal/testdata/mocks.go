package testdata

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

type GoTrueClaims struct {
	jwt.RegisteredClaims
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetaData map[string]interface{} `json:"user_metadata"`
	Role         string                 `json:"role"`
}

func MockAccessToken() string {
	return mockAccessToken("anon_key")
}

func mockAccessToken(role string) string {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &GoTrueClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1234567890",
		},
		Role: role,
	}).SignedString(GotrueJWTSecret)
	if err != nil {
		panic(err)
	}
	return token
}

func MockUserEmail() string {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return strconv.FormatUint(binary.BigEndian.Uint64(buf), 36) + "@example.com"
}

func MockUserPhone() string {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	n := binary.BigEndian.Uint64(buf)
	phone := "8210" + strconv.FormatUint(n|0x8000000000000000, 10)
	phone = phone[:12]
	return phone
}

func MockUserPassword() string {
	buf := make([]byte, 10)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(buf)
}

func MockUserMetadata() map[string]interface{} {
	name := make([]byte, 5)
	_, err := rand.Read(name)
	if err != nil {
		panic(err)
	}
	return map[string]interface{}{
		"profile_image": fmt.Sprintf("https://example.com/avatars/%s.png",
			base64.RawURLEncoding.EncodeToString(name)),
	}
}

func MockAppMetadata() map[string]interface{} {
	return map[string]interface{}{
		"roles": []string{"editor", "publisher"},
	}
}
