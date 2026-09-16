package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AdisGroup/shopier-go"
	"github.com/AdisGroup/shopier-go/oauth"
)

func main() {
	clientID := os.Getenv("SHOPIER_CLIENT_ID")
	clientSecret := os.Getenv("SHOPIER_CLIENT_SECRET")
	redirectURI := os.Getenv("SHOPIER_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/callback"
	}

	oauthCfg := oauth.NewConfig(clientID, clientSecret, redirectURI)

	// Step 1: Redirect merchant to consent page
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		state := "random_csrf_token_string"
		consentURL := oauthCfg.AuthCodeURL(
			state,
			oauth.ScopeOrdersRead,
			oauth.ScopeProductsRead,
			oauth.ScopeProductsWrite,
		)
		http.Redirect(w, r, consentURL, http.StatusTemporaryRedirect)
	})

	// Step 2: Handle OAuth2 callback and code exchange
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code parameter", http.StatusBadRequest)
			return
		}

		token, err := oauthCfg.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, fmt.Sprintf("Token exchange failed: %v", err), http.StatusInternalServerError)
			return
		}

		// Initialize client using the exchanged access token
		client, err := shopier.NewClient(token.AccessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		shopOwner, err := client.Shop.GetOwner(context.Background())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to query shop owner: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Authentication successful! Welcome, %s %s (Account ID: %s)",
			shopOwner.FirstName,
			shopOwner.LastName,
			shopOwner.ID,
		)
	})

	log.Println("OAuth server listening on :8080. Visit http://localhost:8080/login to initiate authorization.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
