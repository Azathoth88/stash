package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/session"
)

const (
	oidcLoginEndpoint    = "/oidc/login"
	oidcCallbackEndpoint = "/oidc/callback"

	// oidcHTTPTimeout bounds requests made to the OIDC provider (discovery,
	// token exchange and JWKS retrieval).
	oidcHTTPTimeout = 30 * time.Second

	// defaultGroupsClaim is used when allowed groups are configured but no
	// explicit groups claim is set.
	defaultGroupsClaim = "groups"

	// loginErrorParam carries a user-facing error message to the login page.
	loginErrorParam = "error"
)

// oidcProvider lazily initialises and caches the OIDC provider and ID token
// verifier for the configured issuer. It re-initialises automatically if the
// issuer or client id changes.
type oidcProvider struct {
	mu       sync.Mutex
	key      string
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

// oidcHTTPClient returns an HTTP client bound to a context, used for all
// requests to the OIDC provider. Binding it to the context ensures that the
// verifier's key set uses a client with sensible timeouts for later JWKS
// refreshes.
func oidcClientContext(ctx context.Context) context.Context {
	return oidc.ClientContext(ctx, &http.Client{Timeout: oidcHTTPTimeout})
}

func (p *oidcProvider) get(oc config.OIDCConfig) (*oidc.Provider, *oidc.IDTokenVerifier, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := oc.Issuer + "\x00" + oc.ClientID
	if p.provider != nil && p.key == key {
		return p.provider, p.verifier, nil
	}

	// Use a background-derived context so that the verifier's key set is not
	// tied to a single request's lifetime.
	ctx := oidcClientContext(context.Background())

	provider, err := oidc.NewProvider(ctx, oc.Issuer)
	if err != nil {
		return nil, nil, fmt.Errorf("initialising OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: oc.ClientID})

	p.provider = provider
	p.verifier = verifier
	p.key = key

	return provider, verifier, nil
}

// oidcOAuthConfig builds the oauth2 config for the current request.
func oidcOAuthConfig(oc config.OIDCConfig, provider *oidc.Provider, redirectURL string) oauth2.Config {
	scopes := []string{oidc.ScopeOpenID}
	if len(oc.Scopes) > 0 {
		for _, s := range oc.Scopes {
			if s != oidc.ScopeOpenID {
				scopes = append(scopes, s)
			}
		}
	} else {
		scopes = append(scopes, "profile", "email")
	}

	return oauth2.Config{
		ClientID:     oc.ClientID,
		ClientSecret: oc.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}
}

// oidcRedirectURL returns the callback URL to send to the provider. It uses the
// configured value if present, otherwise derives it from the request.
func oidcRedirectURL(r *http.Request, oc config.OIDCConfig) string {
	if oc.RedirectURL != "" {
		return oc.RedirectURL
	}

	baseURL, _ := r.Context().Value(BaseURLCtxKey).(string)
	return baseURL + oidcCallbackEndpoint
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// sanitizeReturnURL only permits local (relative) return URLs to prevent open
// redirects through the login flow.
func sanitizeReturnURL(raw string) string {
	if raw == "" {
		return ""
	}

	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	// reject absolute URLs and protocol-relative URLs
	if u.IsAbs() || u.Host != "" {
		return ""
	}

	if !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") {
		return ""
	}

	return u.String()
}

// redirectToLoginWithError sends the user back to the login page with a
// user-facing error message. The message is rendered through html/template and
// is therefore escaped.
func redirectToLoginWithError(w http.ResponseWriter, r *http.Request, msg string) {
	prefix := getProxyPrefix(r)
	q := make(url.Values)
	q.Set(loginErrorParam, msg)
	u := url.URL{
		Path:     prefix + loginEndpoint,
		RawQuery: q.Encode(),
	}
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func handleOIDCLogin(prov *oidcProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := config.GetInstance()
		oc := c.GetOIDCConfig()
		if !oc.IsValid() {
			http.Error(w, "OIDC is not configured", http.StatusNotFound)
			return
		}

		provider, _, err := prov.get(oc)
		if err != nil {
			logger.Errorf("OIDC: %v", err)
			http.Error(w, "OIDC provider is unavailable", http.StatusServiceUnavailable)
			return
		}

		oauthConfig := oidcOAuthConfig(oc, provider, oidcRedirectURL(r, oc))

		state, err := randomString(32)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		nonce, err := randomString(32)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		verifier := oauth2.GenerateVerifier()

		returnTo := sanitizeReturnURL(r.URL.Query().Get(returnURLParam))

		if err := manager.GetInstance().SessionStore.SaveOIDCState(w, r, session.OIDCState{
			State:    state,
			Nonce:    nonce,
			Verifier: verifier,
			ReturnTo: returnTo,
		}); err != nil {
			logger.Errorf("OIDC: error saving login state: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		authURL := oauthConfig.AuthCodeURL(state, oidc.Nonce(nonce), oauth2.S256ChallengeOption(verifier))
		http.Redirect(w, r, authURL, http.StatusFound)
	}
}

func handleOIDCCallback(prov *oidcProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := config.GetInstance()
		oc := c.GetOIDCConfig()
		if !oc.IsValid() {
			http.Error(w, "OIDC is not configured", http.StatusNotFound)
			return
		}

		ss := manager.GetInstance().SessionStore

		st, err := ss.GetOIDCState(r)
		// the login state is single-use; always clear it
		_ = ss.ClearOIDCState(w, r)

		if err != nil || st.State == "" {
			redirectToLoginWithError(w, r, "Login session expired. Please try again.")
			return
		}

		q := r.URL.Query()

		if errParam := q.Get("error"); errParam != "" {
			logger.Warnf("OIDC callback returned error %q: %s", errParam, q.Get("error_description"))
			redirectToLoginWithError(w, r, "Authentication was denied by the provider.")
			return
		}

		if q.Get("state") != st.State {
			redirectToLoginWithError(w, r, "Invalid login state. Please try again.")
			return
		}

		code := q.Get("code")
		if code == "" {
			redirectToLoginWithError(w, r, "Missing authorization code.")
			return
		}

		provider, verifier, err := prov.get(oc)
		if err != nil {
			logger.Errorf("OIDC: %v", err)
			http.Error(w, "OIDC provider is unavailable", http.StatusServiceUnavailable)
			return
		}

		oauthConfig := oidcOAuthConfig(oc, provider, oidcRedirectURL(r, oc))
		ctx := oidcClientContext(r.Context())

		token, err := oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(st.Verifier))
		if err != nil {
			logger.Errorf("OIDC: token exchange failed: %v", err)
			redirectToLoginWithError(w, r, "Failed to complete authentication.")
			return
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok || rawIDToken == "" {
			logger.Errorf("OIDC: no id_token in token response")
			redirectToLoginWithError(w, r, "Provider did not return an ID token.")
			return
		}

		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			logger.Errorf("OIDC: ID token verification failed: %v", err)
			redirectToLoginWithError(w, r, "Failed to verify identity token.")
			return
		}

		if idToken.Nonce != st.Nonce {
			logger.Errorf("OIDC: nonce mismatch")
			redirectToLoginWithError(w, r, "Invalid login state. Please try again.")
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			logger.Errorf("OIDC: failed to parse claims: %v", err)
			redirectToLoginWithError(w, r, "Failed to read identity claims.")
			return
		}

		username := oidcUsername(idToken, claims, oc)
		if username == "" {
			logger.Errorf("OIDC: could not determine username from claims")
			redirectToLoginWithError(w, r, "Could not determine username.")
			return
		}

		if !oidcGroupsAllowed(claims, oc) {
			logger.Warnf("OIDC: user %q denied - not a member of an allowed group", username)
			redirectToLoginWithError(w, r, "You are not authorised to access this instance.")
			return
		}

		// Stash uses a single-user model. When password credentials are also
		// configured, bind the session to the configured username so that
		// features tied to that identity (such as signed media URLs) keep
		// working. Otherwise use the identity provided by the OIDC provider.
		userID := username
		if c.HasCredentials() {
			userID = c.GetUsername()
		}

		if err := ss.LoginUserID(w, r, userID); err != nil {
			logger.Errorf("OIDC: error creating session: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		logger.Info("User logged in via OIDC")

		returnTo := st.ReturnTo
		if returnTo == "" {
			returnTo = getProxyPrefix(r) + "/"
		}
		http.Redirect(w, r, returnTo, http.StatusFound)
	}
}

// oidcUsername resolves the stash username from the configured claim, falling
// back to the token subject.
func oidcUsername(idToken *oidc.IDToken, claims map[string]interface{}, oc config.OIDCConfig) string {
	if v, ok := claims[oc.GetUsernameClaim()].(string); ok && v != "" {
		return v
	}
	return idToken.Subject
}

// oidcGroupsAllowed returns true if the user is permitted to log in based on the
// configured allowed groups. When no groups are configured, all authenticated
// users are permitted.
func oidcGroupsAllowed(claims map[string]interface{}, oc config.OIDCConfig) bool {
	if len(oc.AllowedGroups) == 0 {
		return true
	}

	groupsClaim := oc.GroupsClaim
	if groupsClaim == "" {
		groupsClaim = defaultGroupsClaim
	}

	userGroups := extractGroups(claims[groupsClaim])

	allowed := make(map[string]struct{}, len(oc.AllowedGroups))
	for _, g := range oc.AllowedGroups {
		allowed[g] = struct{}{}
	}

	for _, g := range userGroups {
		if _, ok := allowed[g]; ok {
			return true
		}
	}

	return false
}

// extractGroups normalises the various shapes a groups claim can take into a
// slice of strings.
func extractGroups(v interface{}) []string {
	switch t := v.(type) {
	case []string:
		return t
	case []interface{}:
		out := make([]string, 0, len(t))
		for _, e := range t {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case string:
		return []string{t}
	default:
		return nil
	}
}
