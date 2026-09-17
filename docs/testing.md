# Testing & Diagnostics

The SDK provides testing and diagnostic capabilities for both offline unit testing (with the built-in `testutil` package) and live API smoke testing against the production Shopier platform.

---

## 1. Offline Unit Testing with `testutil.MockServer`

The SDK includes a zero-dependency `testutil` package that runs a local `httptest.Server` and returns an initialized `*shopier.Client`:

```go
package myapp_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/AdisGroup/shopier-go/testutil"
)

func TestOrderRetrieval(t *testing.T) {
	server := testutil.NewMockServer()
	defer server.Close()

	// Register expected endpoint response
	server.HandleFunc("/orders/ord_123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "ord_123",
			"status": "fulfilled",
			"currency": "TRY",
			"totals": {"total": "250.00"},
			"shippingInfo": {
				"firstName": "Ahmet",
				"lastName": "Yilmaz",
				"email": "ahmet@example.com",
				"phone": "+905551234567",
				"address": "Ataturk Cad No:1",
				"city": "Istanbul",
				"country": "Turkiye"
			},
			"lineItems": []
		}`))
	})

	client := server.Client()

	order, err := client.Orders.Get(context.Background(), "ord_123")
	if err != nil {
		t.Fatalf("Orders.Get failed: %v", err)
	}

	if order.ID != "ord_123" || order.Status != "fulfilled" {
		t.Errorf("unexpected order data: %+v", order)
	}
}
```

### Simulating Errors & Retry Scenarios

You can test error resilience for rate limits, validation errors, and server downtimes:

```go
// Simulate Rate Limit (429)
server.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", "1")
	w.WriteHeader(http.StatusTooManyRequests)
	_, _ = w.Write([]byte(`{"error":"rate_limit_exceeded","message":"Too many requests"}`))
})

// Simulate Permission Error (403)
server.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":"forbidden","message":"Access to this resource on the server is denied"}`))
})
```

---

## 2. Live API Diagnostics & Smoke Testing

To verify your Personal Access Token (PAT), store credentials, and active permission scopes against the live Shopier API (`api.shopier.com`), you can run a diagnostic health check script:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/AdisGroup/shopier-go"
)

func main() {
	token := os.Getenv("SHOPIER_API_KEY")
	if token == "" {
		fmt.Println("Please set SHOPIER_API_KEY environment variable")
		return
	}

	client, err := shopier.NewClient(token, shopier.WithTimeout(10*time.Second))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	fmt.Println("🔍 Running Shopier API Diagnostics...")

	// 1. Check Store Identity
	if owner, err := client.Shop.GetOwner(ctx); err != nil {
		fmt.Printf("❌ Shop.GetOwner: %v\n", err)
	} else {
		fmt.Printf("✅ Shop Owner: %s %s (ID: %s)\n", owner.FirstName, owner.LastName, owner.ID)
	}

	// 2. Check Store Settings
	if settings, err := client.Shop.GetSettings(ctx); err != nil {
		fmt.Printf("❌ Shop.GetSettings: %v\n", err)
	} else {
		fmt.Printf("✅ Store Name: %s | URL: %s\n", settings.Name, settings.URL)
	}

	// 3. Check Orders Endpoint
	if orders, err := client.Orders.List(ctx, nil); err != nil {
		fmt.Printf("❌ Orders: %v\n", err)
	} else {
		fmt.Printf("✅ Orders: %d orders found\n", len(orders.Items))
	}

	// 4. Check Balance
	if balances, err := client.Balance.Get(ctx); err != nil {
		fmt.Printf("❌ Balance: %v\n", err)
	} else {
		fmt.Printf("✅ Balance: %s %s\n", balances[0].Amount, balances[0].Currency)
	}

	// 5. Check Product Catalog Scope
	if _, err := client.Products.List(ctx, nil); err != nil {
		fmt.Printf("⚠️ Products: %v (If 403 Forbidden, contact hello@shopier.com)\n", err)
	} else {
		fmt.Println("✅ Products: Access granted")
	}
}
```

---

## 3. Automated CI/CD Testing

To run unit tests across your Go codebase in GitHub Actions or any CI system:

```bash
go test -v -race -cover ./...
```

