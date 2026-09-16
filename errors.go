package shopier

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	// ErrMissingToken indicates that an API request was attempted without an access token.
	ErrMissingToken = errors.New("shopier: access token is required")

	// ErrInvalidBaseURL indicates that the provided API base URL is malformed.
	ErrInvalidBaseURL = errors.New("shopier: invalid base URL")
)

// APIError represents a non-2xx HTTP response from the Shopier REST API.
type APIError struct {
	StatusCode int         `json:"status_code"`
	ErrorCode  string      `json:"error,omitempty"`
	Message    string      `json:"message,omitempty"`
	Detail     string      `json:"error_description,omitempty"`
	RawBody    []byte      `json:"-"`
	Header     http.Header `json:"-"`
}

func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = e.Detail
	}
	if msg == "" && e.ErrorCode != "" {
		msg = e.ErrorCode
	}
	if msg == "" {
		msg = http.StatusText(e.StatusCode)
	}
	if e.ErrorCode != "" && e.ErrorCode != msg {
		return fmt.Sprintf("shopier: status %d (%s): %s", e.StatusCode, e.ErrorCode, msg)
	}
	return fmt.Sprintf("shopier: status %d: %s", e.StatusCode, msg)
}

// RateLimitError represents a 429 Too Many Requests response.
// Shopier enforces rate limits on a 60-second sliding window (default: 200 req/min).
type RateLimitError struct {
	APIError
	// RetryAfter indicates the duration requested by the server before retrying.
	// Parsed from the 'Retry-After' response header (in seconds).
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("%s (retry after %s)", e.APIError.Error(), e.RetryAfter)
	}
	return e.APIError.Error()
}

// AuthenticationError represents a 401 Unauthorized response, typically caused
// by an expired, invalid, or revoked Personal Access Token or OAuth2 access token.
type AuthenticationError struct {
	APIError
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("shopier: authentication failed: %s", e.APIError.Error())
}

// parseAPIError unmarshals an HTTP error response into the appropriate error type.
func parseAPIError(resp *http.Response, body []byte) error {
	apiErr := APIError{
		StatusCode: resp.StatusCode,
		RawBody:    body,
		Header:     resp.Header.Clone(),
	}

	// Attempt to extract structured error details from JSON body.
	if len(body) > 0 {
		_ = json.Unmarshal(body, &apiErr)
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		var retryAfter time.Duration
		if headerVal := resp.Header.Get("Retry-After"); headerVal != "" {
			var seconds int
			if _, err := fmt.Sscanf(headerVal, "%d", &seconds); err == nil && seconds > 0 {
				retryAfter = time.Duration(seconds) * time.Second
			}
		}
		return &RateLimitError{
			APIError:   apiErr,
			RetryAfter: retryAfter,
		}

	case http.StatusUnauthorized:
		return &AuthenticationError{
			APIError: apiErr,
		}

	default:
		return &apiErr
	}
}
