// Package workos adds WorkOS SSO/SAML (enterprise single sign-on) to togo auth.
// It registers /api/auth/workos + callback; on success it find-or-creates the
// user via auth and issues a togo session. Depends on the auth plugin.
//
// Install: `togo install togo-framework/auth-workos`.
package workos

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/togo-framework/auth"
	"github.com/togo-framework/togo"
)

const (
	authorizeURL = "https://api.workos.com/sso/authorize"
	exchangeURL  = "https://api.workos.com/sso/token"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

func init() {
	togo.RegisterProviderFunc("auth-workos", togo.PriorityLate+20, func(k *togo.Kernel) error {
		svc, ok := auth.FromKernel(k)
		if !ok {
			if k.Log != nil {
				k.Log.Warn("auth-workos: auth plugin not installed; skipping")
			}
			return nil
		}
		k.Router.Get("/api/auth/workos", redirectHandler())
		k.Router.Get("/api/auth/workos/callback", callbackHandler(svc))
		return nil
	})
}

func redirectURI() string {
	return strings.TrimRight(os.Getenv("APP_URL"), "/") + "/api/auth/workos/callback"
}

func redirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := os.Getenv("WORKOS_CLIENT_ID")
		if clientID == "" {
			http.Error(w, "workos not configured", http.StatusServiceUnavailable)
			return
		}
		state := randHex()
		http.SetCookie(w, &http.Cookie{Name: "workos_state", Value: state, Path: "/", HttpOnly: true, MaxAge: 600, SameSite: http.SameSiteLaxMode}) //#nosec G124 -- short-lived CSRF state cookie (HttpOnly+SameSite); Secure via TLS/proxy in prod
		q := url.Values{
			"response_type": {"code"},
			"client_id":     {clientID},
			"redirect_uri":  {redirectURI()},
			"state":         {state},
		}
		// One selector is required by WorkOS.
		if c := os.Getenv("WORKOS_CONNECTION"); c != "" {
			q.Set("connection", c)
		} else if o := os.Getenv("WORKOS_ORGANIZATION"); o != "" {
			q.Set("organization", o)
		} else if p := os.Getenv("WORKOS_PROVIDER"); p != "" {
			q.Set("provider", p)
		}
		http.Redirect(w, r, authorizeURL+"?"+q.Encode(), http.StatusFound)
	}
}

func callbackHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("workos_state")
		if err != nil || c.Value == "" || c.Value != r.URL.Query().Get("state") {
			http.Error(w, "invalid sso state", http.StatusBadRequest)
			return
		}
		email, err := exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil || email == "" {
			http.Error(w, "sso exchange failed", http.StatusBadGateway)
			return
		}
		id, err := svc.FindOrCreateByEmail(r.Context(), email)
		if err != nil {
			http.Error(w, "login failed", http.StatusInternalServerError)
			return
		}
		if _, err := svc.IssueSession(w, *id); err != nil {
			http.Error(w, "session failed", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}
}

func exchange(ctx context.Context, code string) (string, error) {
	clientID := os.Getenv("WORKOS_CLIENT_ID")
	apiKey := os.Getenv("WORKOS_API_KEY")
	if clientID == "" || apiKey == "" {
		return "", errors.New("WORKOS_CLIENT_ID/WORKOS_API_KEY not set")
	}
	form := url.Values{
		"client_id":     {clientID},
		"client_secret": {apiKey},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, exchangeURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out struct {
		Profile struct {
			Email string `json:"email"`
		} `json:"profile"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Profile.Email, nil
}

func randHex() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
