package gotrue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/bytebufferpool"

	"github.com/ulbqb/gotrue-go/gotrueapi"
)

type APIClient struct {
	baseURL     string
	baseHeaders Headers
	http        *http.Client
}

func NewAPIClient(url string) *APIClient {
	return &APIClient{
		baseURL:     url,
		baseHeaders: Headers{},
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *APIClient) AppendHeaders(headers Headers) {
	for k, v := range headers {
		c.baseHeaders[k] = v
	}
}

func (c *APIClient) do(req *http.Request, err error) func(out interface{}) error {
	return func(out interface{}) error {
		if err != nil {
			return err
		}

		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode >= 400 {
			var apiErr gotrueapi.Error
			err = json.NewDecoder(resp.Body).Decode(&apiErr)
			if err != nil {
				return fmt.Errorf(
					"api: failed to decode error (%d): %w",
					resp.StatusCode,
					err,
				)
			}
			return &apiErr
		}

		if out != nil {
			err = json.NewDecoder(resp.Body).Decode(out)
			if err != nil {
				return errors.Wrap(err, "api: failed to decode response")
			}
		}

		return nil
	}
}

func (c *APIClient) SignUp(params *gotrueapi.SignUpParams) (*gotrueapi.Session, error) {
	var resp struct {
		*gotrueapi.Session
		*gotrueapi.User
	}

	err := c.do(gotrueapi.SignUp(c.baseURL, c.baseHeaders, params))(&resp)
	if err != nil {
		return nil, err
	}

	if resp.Session != nil && len(resp.Session.Token) > 0 {
		return resp.Session, nil
	}

	return &gotrueapi.Session{User: resp.User}, nil
}

func (c *APIClient) IssueTokenWithPassword(params *gotrueapi.TokenWithPasswordGrantParams) (*gotrueapi.Session, error) {
	var resp gotrueapi.Session

	err := c.do(gotrueapi.TokenWithPasswordGrant(c.baseURL, c.baseHeaders, params))(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *APIClient) IssueTokenWithRefreshToken(params *gotrueapi.TokenWithRefreshTokenGrantParams) (*gotrueapi.Session, error) {
	var resp gotrueapi.Session

	err := c.do(gotrueapi.TokenWithRefreshTokenGrant(c.baseURL, c.baseHeaders, params))(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *APIClient) IssueTokenWithIDToken(params *gotrueapi.TokenWithIDTokenGrantParams) (*gotrueapi.Session, error) {
	var resp gotrueapi.Session

	err := c.do(gotrueapi.TokenWithIDTokenGrant(c.baseURL, c.baseHeaders, params))(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *APIClient) SignOut(accessToken string) error {
	return c.do(gotrueapi.Logout(c.baseURL, c.createRequestHeaders(accessToken)))(nil)
}

func (c *APIClient) SendMagicLinkEmail(params *gotrueapi.MagicLinkParams) error {
	return c.do(gotrueapi.MagicLink(c.baseURL, c.baseHeaders, params))(nil)
}

func (c *APIClient) SendMobileOTP(params *gotrueapi.OTPParams) error {
	return c.do(gotrueapi.OTP(c.baseURL, c.baseHeaders, params))(nil)
}

func (c *APIClient) ResetPasswordForEmail(params *gotrueapi.RecoverParams) error {
	return c.do(gotrueapi.Recover(c.baseURL, c.baseHeaders, params))(nil)
}

func (c *APIClient) GetUser(accessToken string) (*gotrueapi.User, error) {
	var resp gotrueapi.User

	err := c.do(gotrueapi.GetUser(c.baseURL, c.createRequestHeaders(accessToken)))(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *APIClient) UpdateUser(accessToken string, params *gotrueapi.PutUserParams) (*gotrueapi.User, error) {
	var resp gotrueapi.User

	err := c.do(gotrueapi.PutUser(c.baseURL, c.createRequestHeaders(accessToken), params))(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *APIClient) GetProviderSignInURL(provider Provider, redirectTo, scopes string) string {
	pathBuf := bytebufferpool.Get()
	defer bytebufferpool.Put(pathBuf)

	pathBuf.B = append(pathBuf.B, c.baseURL...)
	pathBuf.B = append(pathBuf.B, "/authorize?provider="...)
	pathBuf.B = append(pathBuf.B, url.QueryEscape(string(provider))...)
	pathBuf.B = append(pathBuf.B, "&redirect_to="...)
	pathBuf.B = append(pathBuf.B, redirectTo...)
	pathBuf.B = append(pathBuf.B, "&scopes="...)
	pathBuf.B = append(pathBuf.B, scopes...)

	return string(pathBuf.B)
}

func (c *APIClient) createRequestHeaders(accessToken string) Headers {
	headers := c.baseHeaders.Copy()
	headers["Authorization"] = "Bearer " + accessToken
	return headers
}
