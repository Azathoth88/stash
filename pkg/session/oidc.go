package session

import (
	"net/http"
)

const (
	// oidcStateCookieName is the name of the short-lived cookie that stores
	// the transient state of an in-progress OIDC login flow.
	oidcStateCookieName = "oidc_state"

	oidcStateValueKey    = "state"
	oidcNonceValueKey    = "nonce"
	oidcVerifierValueKey = "verifier"
	oidcReturnToValueKey = "returnTo"

	// oidcStateMaxAge is the lifetime, in seconds, of the transient OIDC login
	// state cookie. The login flow is expected to complete within this window.
	oidcStateMaxAge = 60 * 10 // 10 minutes
)

// OIDCState holds the transient values needed to complete an OIDC login flow.
// These are stored client-side in a signed, short-lived cookie between the
// login redirect and the provider callback.
type OIDCState struct {
	// State is the opaque CSRF token echoed back by the provider.
	State string
	// Nonce binds the ID token to this login attempt.
	Nonce string
	// Verifier is the PKCE code verifier.
	Verifier string
	// ReturnTo is the local URL to redirect to after a successful login.
	ReturnTo string
}

// SaveOIDCState persists the transient OIDC login state in a short-lived cookie.
func (s *Store) SaveOIDCState(w http.ResponseWriter, r *http.Request, st OIDCState) error {
	// ignore error - we always want a fresh state cookie
	newSession, _ := s.sessionStore.New(r, oidcStateCookieName)

	newSession.Values[oidcStateValueKey] = st.State
	newSession.Values[oidcNonceValueKey] = st.Nonce
	newSession.Values[oidcVerifierValueKey] = st.Verifier
	newSession.Values[oidcReturnToValueKey] = st.ReturnTo

	newSession.Options.MaxAge = oidcStateMaxAge
	newSession.Options.SameSite = http.SameSiteLaxMode

	return newSession.Save(r, w)
}

// GetOIDCState reads the transient OIDC login state from the request cookie.
func (s *Store) GetOIDCState(r *http.Request) (OIDCState, error) {
	session, err := s.sessionStore.Get(r, oidcStateCookieName)
	if err != nil {
		return OIDCState{}, err
	}

	st := OIDCState{}
	st.State, _ = session.Values[oidcStateValueKey].(string)
	st.Nonce, _ = session.Values[oidcNonceValueKey].(string)
	st.Verifier, _ = session.Values[oidcVerifierValueKey].(string)
	st.ReturnTo, _ = session.Values[oidcReturnToValueKey].(string)

	return st, nil
}

// ClearOIDCState removes the transient OIDC login state cookie.
func (s *Store) ClearOIDCState(w http.ResponseWriter, r *http.Request) error {
	session, err := s.sessionStore.Get(r, oidcStateCookieName)
	if err != nil {
		// nothing to clear
		return nil
	}

	session.Options.MaxAge = -1
	return session.Save(r, w)
}

// LoginUserID creates an authenticated session for the given user id, without
// validating a password. It is used by external authentication flows such as
// OIDC, where the credentials have already been verified by the provider.
func (s *Store) LoginUserID(w http.ResponseWriter, r *http.Request, userID string) error {
	// ignore error - we want a new session regardless
	newSession, _ := s.sessionStore.Get(r, cookieName)

	newSession.Values[userIDKey] = userID

	return newSession.Save(r, w)
}
