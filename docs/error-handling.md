# Error Handling & Retries

The SDK translates HTTP error codes into strongly typed Go errors and includes built-in retry resilience for transient network and rate limit conditions.

---

## Error Hierarchy

All API error responses implement the Go `error` interface and can be inspected using `errors.As`:

```go
import (
	"errors"
	"github.com/AdisGroup/shopier-go"
)

res, err := client.Orders.Get(ctx, "invalid_id")
if err != nil {
	var rateLimitErr *shopier.RateLimitError
	var authErr *shopier.AuthenticationError
	var apiErr *shopier.APIError

	switch {
	case errors.As(err, &rateLimitErr):
		// 429 Too Many Requests
		fmt.Printf("Rate limit reached. Retry after: %s\n", rateLimitErr.RetryAfter)

	case errors.As(err, &authErr):
		// 401 Unauthorized
		fmt.Println("Invalid or expired API token. Please refresh credentials.")

	case errors.As(err, &apiErr):
		// 400 Bad Request, 404 Not Found, 403 Forbidden, 5xx etc.
		fmt.Printf("API Error [%d]: %s (Detail: %s)\n",
			apiErr.StatusCode,
			apiErr.ErrorCode,
			apiErr.Message,
		)
		// Access raw response payload if needed
		fmt.Println("Raw JSON:", string(apiErr.RawBody))

	default:
		// Context timeout, connection drop, DNS failure
		fmt.Printf("Network/Transport Error: %v\n", err)
	}
}
```

---

## Rate Limit Policy

Shopier limits API throughput on a **60-second sliding window**:
- Standard quota: **200 requests / minute** per application-user pair.
- When the threshold is exceeded, Shopier returns `HTTP 429 Too Many Requests` with a `Retry-After: <seconds>` header.

### Automated Rate Limit Handling

The client automatically detects `429` responses, extracts the `Retry-After` header value, pauses the goroutine until the window resets, and retries the request transparently up to `WithMaxRetries(n)`.

```go
client, err := shopier.NewClient("token",
	shopier.WithMaxRetries(3), // Will retry up to 3 times on 429 and 5xx
)
```

---

## Exponential Backoff & Jitter for 5xx Errors

When Shopier returns transient server errors (`500 Internal Server Error`, `502 Bad Gateway`, `503 Service Unavailable`, `504 Gateway Timeout`), the client calculates an exponential backoff with randomized jitter:

$$\text{delay} = \min(\text{initialDelay} \times 2^{\text{attempt}}, \text{maxDelay}) + \text{jitter}$$

This prevents thundering herd issues on high-volume workers.
