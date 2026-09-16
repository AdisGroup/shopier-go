package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultAuthorizeURL is the Shopier seller consent authorization page.
	DefaultAuthorizeURL = "https://developer.shopier.com/v1/oauth2/authorize"

	// DefaultTokenURL routes token exchange and refresh through port 8443.
	DefaultTokenURL = "https://api.shopier.com:8443/v1/oauth2/token"

	// DefaultRevokeURL routes token revocation through port 8443.
	DefaultRevokeURL = "https://api.shopier.com:8443/v1/oauth2/revoke"
)

// Token represents credentials issued upon successful authorization code exchange or token refresh.
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"-"`
}

// Expiry computes the exact time when this access token expires.
func (t *Token) Expiry() time.Time {
	if t.CreatedAt.IsZero() {
		return time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
	}
	return t.CreatedAt.Add(time.Duration(t.ExpiresIn) * time.Second)
}

// Expired reports whether the access token has expired (with a 60-second grace window).
func (t *Token) Expired() bool {
	return time.Now().After(t.Expiry().Add(-60 * time.Second))
}

// Config describes a registered Shopier OAuth 2.0 app configuration.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AuthorizeURL string
	TokenURL     string
	RevokeURL    string
	HTTPClient   *http.Client
}

// NewConfig initializes an OAuth configuration populated with Shopier's default endpoints.
func NewConfig(clientID, clientSecret, redirectURI string) *Config {
	return &Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		AuthorizeURL: DefaultAuthorizeURL,
		TokenURL:     DefaultTokenURL,
		RevokeURL:    DefaultRevokeURL,
		HTTPClient:   &http.Client{Timeout: 15 * time.Second},
	}
}

// AuthCodeURL constructs the browser redirect URL where sellers grant permission.
func (c *Config) AuthCodeURL(state string, scopes ...string) string {
	u, err := url.Parse(c.AuthorizeURL)
	if err != nil {
		u = &url.URL{Path: c.AuthorizeURL}
	}

	q := u.Query()
	q.Set("client_id", c.ClientID)
	q.Set("redirect_uri", c.RedirectURI)
	q.Set("response_type", "code")
	if len(scopes) > 0 {
		q.Set("scope", strings.Join(scopes, " "))
	}
	if state != "" {
		q.Set("state", state)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// Exchange swaps a temporary authorization code for an initial access and refresh token.
func (c *Config) Exchange(ctx context.Context, code string) (*Token, error) {
	if code == "" {
		return nil, fmt.Errorf("shopier/oauth: authorization code is required")
	}

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"code":          {code},
	}

	return c.postToken(ctx, c.TokenURL, form)
}

// RefreshToken obtains a fresh access token using an unexpired refresh token.
func (c *Config) RefreshToken(ctx context.Context, refreshToken string) (*Token, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("shopier/oauth: refresh token is required")
	}

	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"refresh_token": {refreshToken},
	}

	return c.postToken(ctx, c.TokenURL, form)
}

// Revoke invalidates an issued token (access token or refresh token).
func (c *Config) Revoke(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("shopier/oauth: token is required")
	}

	form := url.Values{
		"token":         {token},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.RevokeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("shopier/oauth: failed to construct revoke request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("shopier/oauth: revoke request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("shopier/oauth: revoke failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Config) postToken(ctx context.Context, endpoint string, form url.Values) (*Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("shopier/oauth: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("shopier/oauth: token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("shopier/oauth: failed to read token response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("shopier/oauth: token request returned status %d: %s", resp.StatusCode, string(body))
	}

	var tok Token
	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, fmt.Errorf("shopier/oauth: failed to decode token JSON: %w", err)
	}
	tok.CreatedAt = time.Now()

	return &tok, nil
}
