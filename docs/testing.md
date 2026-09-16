# Testing & Mocking

The SDK includes a built-in `testutil` package designed to simulate Shopier API responses in your unit and integration test suites without performing network calls against production servers.

---

## Using `testutil.NewMockServer`

`MockServer` starts a local `httptest.Server` and returns a pre-wired `*shopier.Client`:

```go
package myapp_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/AdisGroup/shopier-go/testutil"
)

func TestMyOrderWorker(t *testing.T) {
	server := testutil.NewMockServer()
	defer server.Close()

	// Register expected endpoint responses
	server.HandleFunc("/orders/ord_test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "ord_test",
			"status": "fulfilled",
			"currency": "TRY",
			"totals": {"total": "120.00"},
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

	// Get a client configured for this mock server
	client := server.Client()

	order, err := client.Orders.Get(context.Background(), "ord_test")
	if err != nil {
		t.Fatalf("Orders.Get failed: %v", err)
	}

	if order.Status != "fulfilled" {
		t.Errorf("expected status fulfilled, got %s", order.Status)
	}
}
```

---

## Simulating Error Responses & Retries

You can easily simulate rate limits, validation errors, or server downtimes:

```go
server.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(`{"error":"invalid_title","message":"Product title is required"}`))
})
```
