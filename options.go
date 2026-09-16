package shopier

import (
	"log/slog"
	"net/http"
	"time"
)

// Option represents a functional configuration option for Client.
type Option func(*Client)

// WithBaseURL overrides the default REST API base URL (default: "https://api.shopier.com/v1").
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithOAuthBaseURL overrides the default OAuth base URL (default: "https://api.shopier.com:8443/v1").
// Shopier routes OAuth token exchanges and revocations via port 8443.
func WithOAuthBaseURL(url string) Option {
	return func(c *Client) {
		c.oauthBaseURL = url
	}
}

// WithHTTPClient sets a custom http.Client for network transport.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithMaxRetries sets the maximum retry attempts for transient server errors (5xx)
// and rate limit throttling (429). Setting 0 disables automatic retries.
func WithMaxRetries(retries int) Option {
	return func(c *Client) {
		if retries >= 0 {
			c.maxRetries = retries
		}
	}
}

// WithRetryDelay configures initial delay and maximum delay cap for exponential backoff.
func WithRetryDelay(initial, max time.Duration) Option {
	return func(c *Client) {
		if initial > 0 {
			c.initialRetryDelay = initial
		}
		if max >= initial {
			c.maxRetryDelay = max
		}
	}
}

// WithLogger configures a structured logger using Go's standard log/slog.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithTimeout configures a default client-level timeout applied when a request
// context does not specify its own deadline.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout > 0 {
			c.timeout = timeout
		}
	}
}

// WithUserAgent customizes the User-Agent request header string.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		if userAgent != "" {
			c.userAgent = userAgent
		}
	}
}
