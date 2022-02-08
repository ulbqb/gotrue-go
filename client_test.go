package gotrue

import (
	"strconv"
	"testing"
	"time"

	"go.lair.cx/gotrue-go/gotrueapi"
	"go.lair.cx/gotrue-go/internal/testdata"
)

var (
	authClient                = NewClient(testdata.GotrueURLSignUpEnabledAutoConfirmDisabled)
	authClientWithAutoConfirm = NewClient(testdata.GotrueURLSignUpEnabledAutoConfirmEnabled)
	authClientWithoutSignUp   = NewClient(testdata.GotrueURLSignUpDisabledAutoConfirmDisabled)
)

func TestClient_signUpWithPassword(t *testing.T) {
	t.Run("without auto confirm", func(t *testing.T) {
		t.Run("sign up with email", func(t *testing.T) {
			email := testdata.MockUserEmail()
			session, err := authClient.signUpWithPassword(&gotrueapi.SignUpParams{
				Email:    email,
				Password: testdata.MockUserPassword(),
			})
			if err != nil {
				t.Errorf("signUpWithPassword() error = %v", err)
				return
			}
			if session.User == nil || session.User.Email != email {
				t.Errorf("signUpWithPassword() failed; email = %s, got = %s",
					strconv.Quote(email), strconv.Quote(session.User.Email))
				return
			}
		})
	})

	t.Run("with auto confirm", func(t *testing.T) {
		t.Run("sign up with email", func(t *testing.T) {
			session, err := authClientWithAutoConfirm.signUpWithPassword(&gotrueapi.SignUpParams{
				Email:    testdata.MockUserEmail(),
				Password: testdata.MockUserPassword(),
			})
			if err != nil {
				t.Errorf("signUpWithPassword() error = %v", err)
				return
			}
			if len(session.Token) == 0 {
				t.Errorf("signUpWithPassword() failed; got = %v", session)
				return
			}
		})

		t.Run("sign up with phone", func(t *testing.T) {
			session, err := authClientWithAutoConfirm.signUpWithPassword(&gotrueapi.SignUpParams{
				Phone:    testdata.MockUserPhone(),
				Password: testdata.MockUserPassword(),
			})
			if err != nil {
				t.Errorf("signUpWithPassword() error = %v", err)
				return
			}
			if len(session.Token) == 0 {
				t.Errorf("signUpWithPassword() failed; got = %v", session)
				return
			}
		})

		t.Run("sign up should fire sign in event", func(t *testing.T) {
			var ch = make(chan struct{}, 1)

			unsubscribe := authClientWithAutoConfirm.Subscribe(SignedInEvent, func() {
				ch <- struct{}{}
			})
			defer unsubscribe()

			_, err := authClientWithAutoConfirm.signUpWithPassword(&gotrueapi.SignUpParams{
				Email:    testdata.MockUserEmail(),
				Password: testdata.MockUserPassword(),
			})
			if err != nil {
				t.Errorf("signUpWithPassword() error = %v", err)
				return
			}

			select {
			case <-ch:
				return

			case <-time.After(10 * time.Second):
				t.Errorf("signUpWithPassword() does not call signed in event")
				return
			}
		})
	})

	t.Run("signup disabled", func(t *testing.T) {
		_, err := authClientWithoutSignUp.signUpWithPassword(&gotrueapi.SignUpParams{
			Email:    testdata.MockUserEmail(),
			Password: testdata.MockUserPassword(),
		})
		if err == nil {
			t.Errorf("signUpWithPassword() returns no error")
			return
		}
	})
}

func TestClient_signInWithPasswordGrant(t *testing.T) {
	t.Run("sign in with email", func(t *testing.T) {
		var (
			email    = testdata.MockUserEmail()
			password = testdata.MockUserPassword()
		)
		_, err := authClientWithAutoConfirm.SignUpWithEmail(
			email,
			password,
			nil,
		)
		if err != nil {
			t.Errorf("SignUpWithEmail() error = %v", err)
			return
		}

		sess, err := authClientWithAutoConfirm.SignInWithEmail(email, password)
		if err != nil {
			t.Errorf("SignInWithEmail() error = %v", err)
			return
		}
		if sess.User.Email != email {
			t.Errorf(
				"Sign in failed; email = %s, want = %s",
				strconv.Quote(sess.User.Email),
				strconv.Quote(email),
			)
		}
	})

	t.Run("sign in with phone", func(t *testing.T) {
		var (
			phone    = testdata.MockUserPhone()
			password = testdata.MockUserPassword()
		)
		_, err := authClientWithAutoConfirm.SignUpWithPhone(
			phone,
			password,
			nil,
		)
		if err != nil {
			t.Errorf("SignUpWithEmail() error = %v", err)
			return
		}

		sess, err := authClientWithAutoConfirm.SignInWithPhone(phone, password)
		if err != nil {
			t.Errorf("SignInWithEmail() error = %v", err)
			return
		}
		if sess.User.Phone != phone {
			t.Errorf(
				"Sign in failed; phone = %s, want = %s",
				strconv.Quote(sess.User.Phone),
				strconv.Quote(phone),
			)
		}
	})
}

func TestClient_UpdateUser(t *testing.T) {
	t.Run("sign in with email", func(t *testing.T) {
		_, err := authClientWithAutoConfirm.SignUpWithEmail(
			testdata.MockUserEmail(),
			testdata.MockUserPassword(),
			nil,
		)
		if err != nil {
			t.Errorf("SignUpWithEmail() error = %v", err)
			return
		}

		updated, err := authClientWithAutoConfirm.UpdateUser(&gotrueapi.PutUserParams{
			Data: map[string]interface{}{
				"Hello": "World",
			},
		})
		if err != nil {
			t.Errorf("UpdateUser() error = %v", err)
			return
		}

		if updated.UserMetaData["Hello"] != "World" {
			t.Errorf(
				"User update faild; want = World, got = %v",
				updated.UserMetaData["Hello"],
			)
			return
		}
	})
}

func TestClient_SignOut(t *testing.T) {
	t.Run("get user after sign out", func(t *testing.T) {
		_, err := authClientWithAutoConfirm.SignUpWithEmail(
			testdata.MockUserEmail(),
			testdata.MockUserPassword(),
			nil,
		)
		if err != nil {
			t.Errorf("SignUpWithEmail() error = %v", err)
			return
		}
		err = authClientWithAutoConfirm.SignOut()
		if err != nil {
			t.Errorf("SignOut() error = %v", err)
			return
		}

		u := authClientWithAutoConfirm.User()
		if u != nil {
			t.Errorf("GetUser() returns user data after signed out")
		}
	})
}
