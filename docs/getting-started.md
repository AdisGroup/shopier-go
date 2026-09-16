# Overview & Installation

`shopier-go` is an idiomatic Go SDK and API wrapper for the Shopier REST API. It is designed for high-concurrency production environments, offering zero external dependencies, robust retry logic, context propagation, and type safety across all Shopier resources.

## Requirements

- Go 1.23 or higher

## Installation

Install the package using the standard `go get` command:

```bash
go get github.com/AdisGroup/shopier-go
```

## Quickstart

Initialize the client with either a Personal Access Token (PAT) or an OAuth2 access token:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AdisGroup/shopier-go"
)

func main() {
	client, err := shopier.NewClient(os.Getenv("SHOPIER_TOKEN"))
	if err != nil {
		log.Fatalf("Client init failed: %v", err)
	}

	ctx := context.Background()

	// Fetch current account balances
	balances, err := client.Balance.Get(ctx)
	if err != nil {
		log.Fatalf("Balance inquiry failed: %v", err)
	}

	for _, b := range balances {
		fmt.Printf("Balance: %s %s\n", b.Amount, b.Currency)
	}
}
```

## Client Configuration

The client constructor accepts functional options to customize transport, logging, timeouts, and retry policies:

```go
import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AdisGroup/shopier-go"
)

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

client, err := shopier.NewClient(os.Getenv("SHOPIER_TOKEN"),
	shopier.WithTimeout(15*time.Second),
	shopier.WithMaxRetries(3),
	shopier.WithRetryDelay(500*time.Millisecond, 10*time.Second),
	shopier.WithLogger(logger),
	shopier.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	shopier.WithUserAgent("my-custom-service/1.0"),
)
```

### Configuration Options

| Option | Default | Description |
| :--- | :--- | :--- |
| `WithBaseURL(url)` | `https://api.shopier.com/v1` | Custom REST API base URL |
| `WithOAuthBaseURL(url)` | `https://api.shopier.com:8443/v1` | Custom OAuth2 base URL (port 8443) |
| `WithTimeout(d)` | `30s` | Fallback request timeout when context has no deadline |
| `WithMaxRetries(n)` | `3` | Maximum retry attempts on 429 and 5xx responses |
| `WithRetryDelay(initial, max)` | `500ms`, `10s` | Exponential backoff delay parameters |
| `WithLogger(logger)` | Discard logger | Structured `*slog.Logger` instance |
| `WithHTTPClient(client)` | Default client | Custom `*http.Client` network transport |
| `WithUserAgent(ua)` | `shopier-go/1.0.0` | Custom User-Agent header value |
