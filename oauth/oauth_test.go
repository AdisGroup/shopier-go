package oauth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AdisGroup/shopier-go/oauth"
)

func TestOAuth_AuthCodeURL(t *testing.T) {
	cfg := oauth.NewConfig("my_client_id", "my_secret", "https://myapp.com/callback")
	urlStr := cfg.AuthCodeURL("random_state", oauth.ScopeOrdersRead, oauth.ScopeProductsWrite)

	if !strings.Contains(urlStr, "client_id=my_client_id") {
		t.Errorf("url missing client_id: %s", urlStr)
	}
	if !strings.Contains(urlStr, "redirect_uri=https%3A%2F%2Fmyapp.com%2Fcallback") {
		t.Errorf("url missing redirect_uri: %s", urlStr)
	}
	if !strings.Contains(urlStr, "state=random_state") {
		t.Errorf("url missing state: %s", urlStr)
	}
	if !strings.Contains(urlStr, "orders%3Aread") {
		t.Errorf("url missing scope orders:read: %s", urlStr)
	}
}

func TestOAuth_ExchangeAndRefresh(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		grantType := r.Form.Get("grant_type")

		if grantType == "authorization_code" {
			if r.Form.Get("code") != "auth_code_123" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"access_token": "acc_tok_abc",
				"token_type": "bearer",
				"expires_in": 259200,
				"refresh_token": "ref_tok_xyz"
			}`))
			return
		}

		if grantType == "refresh_token" {
			if r.Form.Get("refresh_token") != "ref_tok_xyz" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"access_token": "acc_tok_new",
				"token_type": "bearer",
				"expires_in": 259200,
				"refresh_token": "ref_tok_new"
			}`))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	cfg := oauth.NewConfig("client_id", "secret", "https://myapp.com/cb")
	cfg.TokenURL = server.URL
	cfg.HTTPClient = server.Client()

	tok, err := cfg.Exchange(context.Background(), "auth_code_123")
	if err != nil {
		t.Fatalf("Exchange failed: %v", err)
	}
	if tok.AccessToken != "acc_tok_abc" || tok.RefreshToken != "ref_tok_xyz" {
		t.Errorf("unexpected token payload: %+v", tok)
	}
	if tok.Expired() {
		t.Error("fresh token should not be expired")
	}

	refreshed, err := cfg.RefreshToken(context.Background(), tok.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}
	if refreshed.AccessToken != "acc_tok_new" {
		t.Errorf("unexpected refreshed token: %+v", refreshed)
	}
}

func TestOAuth_Revoke(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("token") == "valid_token" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	cfg := oauth.NewConfig("client_id", "secret", "https://cb")
	cfg.RevokeURL = server.URL
	cfg.HTTPClient = server.Client()

	if err := cfg.Revoke(context.Background(), "valid_token"); err != nil {
		t.Fatalf("Revoke failed: %v", err)
	}
}
