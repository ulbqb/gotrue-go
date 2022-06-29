package gotrue

import (
	"strconv"
	"sync"

	"github.com/pkg/errors"

	"github.com/ulbqb/gotrue-go/gotrueapi"
)

type Client struct {
	sync.RWMutex

	api *APIClient

	currentSession *gotrueapi.Session
	currentUser    *gotrueapi.User

	eventChannel *EventChannel
}

func NewClient(url string) *Client {
	return &Client{
		api:          NewAPIClient(url),
		eventChannel: NewEventChannel(),
	}
}

// SignUpWithEmail creates new account with email address.
func (c *Client) SignUpWithEmail(email, password string, data interface{}) (*gotrueapi.Session, error) {
	return c.signUpWithPassword(&gotrueapi.SignUpParams{
		Email:    email,
		Password: password,
		Data:     data,
	})
}

// SignUpWithPhone creates new account with phone number.
func (c *Client) SignUpWithPhone(phone, password string, data interface{}) (*gotrueapi.Session, error) {
	return c.signUpWithPassword(&gotrueapi.SignUpParams{
		Phone:    phone,
		Password: password,
		Data:     data,
	})
}

// SignInWithEmail issues access token with user email and password.
func (c *Client) SignInWithEmail(email, password string) (*gotrueapi.Session, error) {
	return c.signInWithPasswordGrant(&gotrueapi.TokenWithPasswordGrantParams{
		Email:    email,
		Password: password,
	})
}

// SignInWithPhone issues access token with user's phone number and password.
func (c *Client) SignInWithPhone(phone, password string) (*gotrueapi.Session, error) {
	return c.signInWithPasswordGrant(&gotrueapi.TokenWithPasswordGrantParams{
		Phone:    phone,
		Password: password,
	})
}

// SignInWithMagicLink sends "magic link" email to user.
func (c *Client) SignInWithMagicLink(params *gotrueapi.MagicLinkParams) error {
	return c.api.SendMagicLinkEmail(params)
}

// SignInWithOTP sends OTP to user's phone.
func (c *Client) SignInWithOTP(params *gotrueapi.OTPParams) error {
	return c.api.SendMobileOTP(params)
}

// SignInWithProvider returns sign in url for provider.
func (c *Client) SignInWithProvider(provider Provider, redirectTo, scopes string) string {
	return c.api.GetProviderSignInURL(provider, redirectTo, scopes)
}

// SignOut destroys current session. Note that revoked token is still be valid
// for stateless services.
func (c *Client) SignOut() error {
	c.Lock()
	defer c.Unlock()

	if c.currentSession == nil {
		return nil
	}

	token := c.currentSession.Token
	c.destroySession()
	c.eventChannel.Publish(SignedOutEvent)

	if len(token) == 0 {
		return nil
	}

	return c.api.SignOut(token)
}

// ResetPasswordForEmail sends a recover email to the user.
func (c *Client) ResetPasswordForEmail(params *gotrueapi.RecoverParams) error {
	return c.api.ResetPasswordForEmail(params)
}

// Session returns copy of current session. returned session is only vaild until
// next refresh.
func (c *Client) Session() *gotrueapi.Session {
	c.RLock()
	defer c.RUnlock()
	if c.currentSession == nil {
		return nil
	}
	s := *c.currentSession
	return &s
}

// GetSessionFromURL parses url and returns generated session.
func (c *Client) GetSessionFromURL(url string, storeSession bool) (*gotrueapi.Session, error) {
	values, err := getParametersFromURI(url)
	if err != nil {
		return nil, err
	}

	if v := values.Get("error_description"); len(v) > 0 {
		return nil, errors.New(v)
	}

	accessToken := values.Get("access_token")
	if len(accessToken) == 0 {
		return nil, errors.New("api: no access_token was provided")
	}

	expiresIn, err := strconv.ParseInt(values.Get("expires_in"), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "api: expires_in is invalid")
	}

	refreshToken := values.Get("refresh_token")
	if len(refreshToken) == 0 {
		return nil, errors.New("api: no refresh_token was provided")
	}

	tokenType := values.Get("token_type")
	if len(tokenType) == 0 {
		return nil, errors.New("api: no token_type was provided")
	}

	user, err := c.api.GetUser(accessToken)
	if err != nil {
		return nil, err
	}

	session := &gotrueapi.Session{
		Token:        accessToken,
		TokenType:    tokenType,
		ExpiresIn:    int(expiresIn),
		RefreshToken: refreshToken,
		User:         user,
	}

	if storeSession {
		c.saveSession(session)
		c.eventChannel.Publish(SignedInEvent)
		if values.Get("type") == "recovery" {
			c.eventChannel.Publish(PasswordRecoveryEvent)
		}
	}

	return session, nil
}

func (c *Client) User() *gotrueapi.User {
	c.RLock()
	defer c.RUnlock()
	if c.currentUser == nil {
		return nil
	}
	u := *c.currentUser
	return &u
}

// UpdateUser updates current user with provided params and returns updated user.
func (c *Client) UpdateUser(params *gotrueapi.PutUserParams) (*gotrueapi.User, error) {
	c.Lock()
	defer c.Unlock()

	if c.currentSession == nil || len(c.currentSession.Token) == 0 {
		return nil, errors.New("not signed in")
	}

	user, err := c.api.UpdateUser(c.currentSession.Token, params)
	if err != nil {
		return nil, err
	}

	c.currentSession.User = user
	c.eventChannel.Publish(UserUpdatedEvent)

	return user, nil
}

func (c *Client) Subscribe(event AuthChangeEvent, fn func()) func() {
	return c.eventChannel.Subscribe(event, fn)
}

func (c *Client) signUpWithPassword(params *gotrueapi.SignUpParams) (*gotrueapi.Session, error) {
	c.Lock()
	defer c.Unlock()

	c.destroySession()

	session, err := c.api.SignUp(params)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("an error occurred on sign up")
	}

	// Handles when auto confirm is set.
	if len(session.Token) > 0 {
		c.saveSession(session)
		c.eventChannel.Publish(SignedInEvent)
	}

	return session, nil
}

func (c *Client) signInWithPasswordGrant(params *gotrueapi.TokenWithPasswordGrantParams) (*gotrueapi.Session, error) {
	c.Lock()
	defer c.Unlock()

	c.destroySession()

	session, err := c.api.IssueTokenWithPassword(params)
	if err != nil {
		return nil, err
	}

	if session.User != nil && session.User.EmailConfirmedAt != nil {
		c.saveSession(session)
		c.eventChannel.Publish(SignedInEvent)
	}

	return session, err
}

// refreshSession refreshes session.
// refreshSession is not thread safe.
func (c *Client) refreshSession() (*gotrueapi.Session, error) {
	if c.currentSession == nil {
		return nil, errors.New("no session")
	}
	refreshToken := c.currentSession.RefreshToken
	return c.refreshSessionWithToken(refreshToken)
}

// refreshSessionWithToken refreshes session with provided refresh token.
// refreshSessionWithToken is not thread safe.
func (c *Client) refreshSessionWithToken(refreshToken string) (*gotrueapi.Session, error) {
	session, err := c.api.IssueTokenWithRefreshToken(&gotrueapi.TokenWithRefreshTokenGrantParams{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}

	c.saveSession(session)
	c.eventChannel.Publish(TokenRefreshedEvent)
	c.eventChannel.Publish(SignedInEvent)

	return session, nil
}

// saveSession saves the token. saveSession does not fire any events and
// not thread safe.
func (c *Client) saveSession(session *gotrueapi.Session) {
	c.currentSession = session
	c.currentUser = session.User
}

// destroySession destroys the session. destroySession is not thread safe.
func (c *Client) destroySession() {
	c.currentSession = nil
	c.currentUser = nil
}
