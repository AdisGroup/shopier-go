package webhook_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdisGroup/shopier-go/webhook"
)

func TestVerifySignature(t *testing.T) {
	token := "whsec_supersecretkey"
	payload := []byte(`{"id":"ord_test","status":"fulfilled"}`)

	mac := hmac.New(sha256.New, []byte(token))
	mac.Write(payload)
	validHex := hex.EncodeToString(mac.Sum(nil))

	if !webhook.VerifySignature(payload, validHex, token) {
		t.Fatal("expected valid hex signature to pass verification")
	}

	if webhook.VerifySignature(payload, "invalid_sig", token) {
		t.Fatal("expected invalid signature to fail verification")
	}

	if webhook.VerifySignature(payload, validHex, "wrong_token") {
		t.Fatal("expected signature with wrong token to fail verification")
	}
}

func TestWebhook_Handler(t *testing.T) {
	token := "whsec_sample123"
	payload := []byte(`{"id":"ord_webhook_1","status":"unfulfilled","currency":"TRY","totals":{"total":"199.99"},"shippingInfo":{"firstName":"Ali","lastName":"Kaya","email":"ali@kaya.com","phone":"555123","address":"Kadikoy","city":"Istanbul","country":"Turkiye"},"lineItems":[]}`)

	mac := hmac.New(sha256.New, []byte(token))
	mac.Write(payload)
	validSig := hex.EncodeToString(mac.Sum(nil))

	var receivedEvent string
	var receivedOrderID string

	handler := webhook.NewHandler(token, func(ctx context.Context, evt *webhook.Event) error {
		receivedEvent = evt.Header.Event
		ord, err := evt.Order()
		if err != nil {
			return err
		}
		receivedOrderID = ord.ID
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("Shopier-Signature", validSig)
	req.Header.Set("Shopier-Event", webhook.EventOrderCreated)
	req.Header.Set("Shopier-Webhook-Id", "wh_evt_999")
	req.Header.Set("Shopier-Timestamp", "1658402691")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got status %d: %s", rec.Code, rec.Body.String())
	}
	if receivedEvent != webhook.EventOrderCreated {
		t.Errorf("expected event %s, got %s", webhook.EventOrderCreated, receivedEvent)
	}
	if receivedOrderID != "ord_webhook_1" {
		t.Errorf("expected order ID ord_webhook_1, got %s", receivedOrderID)
	}
}
