package shopier

import (
	"context"
	"fmt"
	"net/http"
)

// WebhookSubscription represents an event callback registration with Shopier.
type WebhookSubscription struct {
	ID    string `json:"id"`
	Event string `json:"event"`
	URL   string `json:"url"`
	// Token is the secret signing key used to compute HMAC-SHA256 signatures.
	// Shopier returns this value only on the initial subscription creation response.
	Token string `json:"token,omitempty"`
}

// WebhookCreateRequest defines parameters for POST /webhooks.
type WebhookCreateRequest struct {
	Event string `json:"event"`
	URL   string `json:"url"`
}

// WebhookService exposes webhook registration and subscription lifecycle endpoints.
type WebhookService struct {
	client *Client
}

// List retrieves all active webhook event subscriptions.
func (s *WebhookService) List(ctx context.Context) ([]WebhookSubscription, error) {
	var subscriptions []WebhookSubscription
	_, err := s.client.execute(ctx, http.MethodGet, "/webhooks", nil, nil, &subscriptions)
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// Create registers a new notification URL for a specific webhook event.
// Save the returned Token securely, as Shopier will not expose it in subsequent List queries.
func (s *WebhookService) Create(ctx context.Context, req *WebhookCreateRequest) (*WebhookSubscription, error) {
	if req == nil || req.Event == "" || req.URL == "" {
		return nil, fmt.Errorf("shopier: event and url are required")
	}

	var sub WebhookSubscription
	_, err := s.client.execute(ctx, http.MethodPost, "/webhooks", nil, req, &sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// Delete permanently removes a webhook subscription by its ID.
func (s *WebhookService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("shopier: webhook subscription id is required")
	}

	_, err := s.client.execute(ctx, http.MethodDelete, "/webhooks/"+id, nil, nil, nil)
	return err
}
