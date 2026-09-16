package shopier

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the standard Shopier REST API v1 endpoint.
	DefaultBaseURL = "https://api.shopier.com/v1"

	// DefaultOAuthBaseURL is the Shopier OAuth 2.0 endpoint running on port 8443.
	DefaultOAuthBaseURL = "https://api.shopier.com:8443/v1"

	// Version is the current release version of this SDK.
	Version = "1.0.0"

	// DefaultUserAgent sent with outbound API requests.
	DefaultUserAgent = "shopier-go/" + Version
)

// Client coordinates API requests across all Shopier resource services.
type Client struct {
	token             string
	baseURL           string
	oauthBaseURL      string
	httpClient        *http.Client
	maxRetries        int
	initialRetryDelay time.Duration
	maxRetryDelay     time.Duration
	logger            *slog.Logger
	timeout           time.Duration
	userAgent         string

	// Resource services
	Balance    *BalanceService
	Categories *CategoryService
	Discounts  *DiscountService
	Orders     *OrderService
	Payouts    *PayoutService
	Products   *ProductService
	Refunds    *RefundService
	Selections *SelectionService
	Shippings  *ShippingService
	Shop       *ShopService
	Variations *VariationService
	Webhooks   *WebhookService
}

// NewClient returns a configured Shopier API client.
// An access token (Personal Access Token or OAuth2 bearer token) is required.
func NewClient(token string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(token) == "" {
		return nil, ErrMissingToken
	}

	c := &Client{
		token:             strings.TrimSpace(token),
		baseURL:           DefaultBaseURL,
		oauthBaseURL:      DefaultOAuthBaseURL,
		httpClient:        &http.Client{Timeout: 30 * time.Second},
		maxRetries:        3,
		initialRetryDelay: 500 * time.Millisecond,
		maxRetryDelay:     10 * time.Second,
		logger:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		timeout:           30 * time.Second,
		userAgent:         DefaultUserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	// Canonicalize base URLs (strip trailing slashes)
	c.baseURL = strings.TrimRight(c.baseURL, "/")
	c.oauthBaseURL = strings.TrimRight(c.oauthBaseURL, "/")

	// Wire service instances
	c.Balance = &BalanceService{client: c}
	c.Categories = &CategoryService{client: c}
	c.Discounts = &DiscountService{client: c}
	c.Orders = &OrderService{client: c}
	c.Payouts = &PayoutService{client: c}
	c.Products = &ProductService{client: c}
	c.Refunds = &RefundService{client: c}
	c.Selections = &SelectionService{client: c}
	c.Shippings = &ShippingService{client: c}
	c.Shop = &ShopService{client: c}
	c.Variations = &VariationService{client: c}
	c.Webhooks = &WebhookService{client: c}

	return c, nil
}

// execute dispatches an HTTP request with exponential backoff and rate limit handling.
func (c *Client) execute(ctx context.Context, method, endpoint string, query url.Values, body any, dest any) (http.Header, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, hasDeadline := ctx.Deadline(); !hasDeadline && c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	fullURL := c.baseURL + "/" + strings.TrimLeft(endpoint, "/")
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var reqBytes []byte
	if body != nil {
		var err error
		reqBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("shopier: failed to marshal request body: %w", err)
		}
	}

	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.InfoContext(ctx, "retrying Shopier API request",
				slog.String("method", method),
				slog.String("url", fullURL),
				slog.Int("attempt", attempt),
			)
		}

		var bodyReader io.Reader
		if reqBytes != nil {
			bodyReader = bytes.NewReader(reqBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("shopier: failed to create HTTP request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)
		if reqBytes != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		start := time.Now()
		resp, err = c.httpClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			// Do not retry on context cancellation or deadline expiration
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			lastErr = err
			c.logger.WarnContext(ctx, "HTTP request failed",
				slog.String("method", method),
				slog.String("url", fullURL),
				slog.Int("attempt", attempt),
				slog.Duration("duration", duration),
				slog.String("error", err.Error()),
			)
			if attempt < c.maxRetries {
				if sleepErr := c.backoffSleep(ctx, attempt, 0); sleepErr != nil {
					return nil, sleepErr
				}
				continue
			}
			return nil, fmt.Errorf("shopier: network request failed: %w", lastErr)
		}

		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("shopier: failed to read response body: %w", readErr)
		}

		// Handle 2xx success responses
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			c.logger.DebugContext(ctx, "API request successful",
				slog.String("method", method),
				slog.String("url", fullURL),
				slog.Int("status", resp.StatusCode),
				slog.Duration("duration", duration),
			)
			if dest != nil && len(respBody) > 0 {
				if err := json.Unmarshal(respBody, dest); err != nil {
					return resp.Header, fmt.Errorf("shopier: failed to decode response JSON: %w", err)
				}
			}
			return resp.Header, nil
		}

		apiErr := parseAPIError(resp, respBody)

		// 429 Rate Limit: Respect Retry-After header
		if resp.StatusCode == http.StatusTooManyRequests && attempt < c.maxRetries {
			var retryAfter time.Duration
			if rlErr, ok := apiErr.(*RateLimitError); ok {
				retryAfter = rlErr.RetryAfter
			}
			c.logger.WarnContext(ctx, "rate limited by Shopier API",
				slog.Duration("retry_after", retryAfter),
				slog.Int("attempt", attempt),
			)
			if sleepErr := c.backoffSleep(ctx, attempt, retryAfter); sleepErr != nil {
				return nil, sleepErr
			}
			lastErr = apiErr
			continue
		}

		// 5xx Server Errors: Retry with exponential backoff
		if resp.StatusCode >= 500 && attempt < c.maxRetries {
			c.logger.WarnContext(ctx, "Shopier server error, scheduling retry",
				slog.Int("status", resp.StatusCode),
				slog.Int("attempt", attempt),
			)
			if sleepErr := c.backoffSleep(ctx, attempt, 0); sleepErr != nil {
				return nil, sleepErr
			}
			lastErr = apiErr
			continue
		}

		return resp.Header, apiErr
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("shopier: request failed after %d retries", c.maxRetries)
}

// backoffSleep calculates exponential backoff delay with full jitter.
func (c *Client) backoffSleep(ctx context.Context, attempt int, explicitDelay time.Duration) error {
	var sleepDuration time.Duration

	if explicitDelay > 0 {
		sleepDuration = explicitDelay
	} else {
		// Backoff formula: initialDelay * 2^attempt capped at maxRetryDelay
		multiplier := 1 << attempt
		calculated := c.initialRetryDelay * time.Duration(multiplier)
		if calculated > c.maxRetryDelay {
			calculated = c.maxRetryDelay
		}

		// Add jitter between 0% and 50% using zero external dependencies
		jitterMax := int64(calculated / 2)
		if jitterMax > 0 {
			var b [8]byte
			_, _ = rand.Read(b[:])
			randVal := int64(binary.LittleEndian.Uint64(b[:])) & 0x7FFFFFFFFFFFFFFF
			jitter := time.Duration(randVal % jitterMax)
			calculated += jitter
		}
		sleepDuration = calculated
	}

	timer := time.NewTimer(sleepDuration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
