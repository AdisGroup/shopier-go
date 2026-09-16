package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	// ErrInvalidSignature indicates the computed HMAC does not match Shopier-Signature.
	ErrInvalidSignature = errors.New("shopier/webhook: invalid signature")

	// ErrMissingSignature indicates the Shopier-Signature header was absent.
	ErrMissingSignature = errors.New("shopier/webhook: missing Shopier-Signature header")

	// ErrMissingEvent indicates the Shopier-Event header was absent.
	ErrMissingEvent = errors.New("shopier/webhook: missing Shopier-Event header")
)

// VerifySignature validates a webhook payload against its signature token using HS256 (HMAC-SHA256).
// It supports both hexadecimal and standard Base64 signature encodings via constant-time comparison.
func VerifySignature(payload []byte, signature, token string) bool {
	if signature == "" || token == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(token))
	mac.Write(payload)
	sum := mac.Sum(nil)

	expectedHex := hex.EncodeToString(sum)
	if hmac.Equal([]byte(strings.ToLower(signature)), []byte(strings.ToLower(expectedHex))) {
		return true
	}

	expectedBase64 := base64.StdEncoding.EncodeToString(sum)
	return hmac.Equal([]byte(signature), []byte(expectedBase64))
}

// Parse extracts, verifies, and returns an Event from an incoming HTTP webhook request.
// It restores the request body after reading so downstream handlers may inspect it if desired.
func Parse(r *http.Request, token string) (*Event, error) {
	sig := r.Header.Get("Shopier-Signature")
	if sig == "" {
		return nil, ErrMissingSignature
	}

	event := r.Header.Get("Shopier-Event")
	if event == "" {
		return nil, ErrMissingEvent
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("shopier/webhook: failed to read request body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if !VerifySignature(bodyBytes, sig, token) {
		return nil, ErrInvalidSignature
	}

	var ts int64
	if tsHeader := r.Header.Get("Shopier-Timestamp"); tsHeader != "" {
		ts, _ = strconv.ParseInt(tsHeader, 10, 64)
	}

	return &Event{
		Header: Header{
			WebhookID:  r.Header.Get("Shopier-Webhook-Id"),
			Event:      event,
			Timestamp:  ts,
			Signature:  sig,
			AccountID:  r.Header.Get("Shopier-Account-Id"),
			APIVersion: r.Header.Get("Shopier-Api-Version"),
		},
		RawPayload: bodyBytes,
	}, nil
}
