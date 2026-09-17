# Quick Checkout (Direct Payment Link)

The `QuickCheckoutService` generates a direct Shopier checkout URL (`https://www.shopier.com/s/payment/{account}/{order_id}`) by executing the frontend checkout flow on behalf of the customer.

It can be accessed directly from `client.QuickCheckout` on a configured `shopier.Client`, or initialized standalone via `shopier.NewQuickCheckoutService()` without needing an API token.

---

## Concurrency and Session Isolation

Each invocation of `Create` instantiates a dedicated, isolated cookie jar (`net/http/cookiejar`) and executes an independent HTTP session. This ensures that concurrent requests never leak cookies, sessions, or CSRF tokens across goroutines.

---

## Quickstart Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AdisGroup/shopier-go"
)

func main() {
	qc := shopier.NewQuickCheckoutService(
		shopier.WithQuickCheckoutTimeout(30 * time.Second),
	)

	ctx := context.Background()

	req := &shopier.QuickCheckoutRequest{
		ProductURL: "https://www.shopier.com/adisgroup/50900391",
		Quantity:   1,
		Buyer: shopier.QuickCheckoutBuyer{
			Email:     "customer@example.com",
			FirstName: "Özgün Deniz",
			LastName:  "Küçük",
			Phone:     "+90 532 477 02 10",
			Country:   "Türkiye",
			City:      "İstanbul",
			Address:   "Merkez Mah. No: 1",
			Comment:   "Please deliver during work hours.",
		},
	}

	result, err := qc.Create(ctx, req)
	if err != nil {
		var qErr *shopier.QuickCheckoutError
		if qErr, ok := err.(*shopier.QuickCheckoutError); ok {
			log.Fatalf("Checkout failed at stage [%s]: HTTP %d - %s",
				qErr.Stage, qErr.StatusCode, qErr.Message)
		}
		log.Fatalf("Checkout failed: %v", err)
	}

	fmt.Printf("Order ID   : %s\n", result.OrderID)
	fmt.Printf("Payment URL: %s\n", result.PaymentURL)
}
```

---

## Request Parameters

### `QuickCheckoutRequest`

| Field | Type | Description |
| :--- | :--- | :--- |
| `ProductURL` | `string` | Full product URL (e.g. `https://www.shopier.com/adisgroup/50900391`). `Account` and `ProductID` are automatically parsed if set. |
| `Account` | `string` | Seller account username / shop slug (e.g. `adisgroup`). |
| `ProductID` | `string` | Product identifier (e.g. `50900391`). |
| `Quantity` | `int` | Purchase quantity. Defaults to `1` if omitted. |
| `Options` | `map[string]string` | Optional variation or selection form values. |
| `Buyer` | `QuickCheckoutBuyer` | Buyer contact and shipping information. |

### `QuickCheckoutBuyer`

| Field | Type | Description |
| :--- | :--- | :--- |
| `Email` | `string` | **Required.** Buyer email address. |
| `FirstName` | `string` | **Required.** Buyer first name. |
| `LastName` | `string` | **Required.** Buyer last name. |
| `Phone` | `string` | **Required.** Phone number (e.g. `+90 532 477 02 10`). |
| `PhoneCode` | `string` | Country dial code (defaults to `"TR"`). |
| `Country` | `string` | Destination country (defaults to `"Türkiye"`). |
| `City` | `string` | Delivery city. |
| `State` | `string` | Delivery state/district. |
| `Address` | `string` | Full street address. |
| `ZipCode` | `string` | Postal code. |
| `TCIDNo` | `string` | National identity number (optional). |
| `Comment` | `string` | Order notes or delivery instructions. |

---

## Result Structure

### `QuickCheckoutResult`

```go
type QuickCheckoutResult struct {
	OrderID     string         `json:"order_id"`
	PaymentURL  string         `json:"payment_url"`
	Account     string         `json:"account"`
	ProductID   string         `json:"product_id"`
	RawResponse map[string]any `json:"raw_response,omitempty"`
}
```

---

## Configuration Options

| Option | Description |
| :--- | :--- |
| `WithQuickCheckoutBaseURL(url)` | Sets custom base URL (default: `https://www.shopier.com`). |
| `WithQuickCheckoutUserAgent(ua)` | Custom User-Agent string sent with requests. |
| `WithQuickCheckoutHTTPClient(client)` | Custom `*http.Client` for proxies and transport configurations. |
| `WithQuickCheckoutTimeout(d)` | Request timeout duration (default: `30s`). |

---

## Error Handling

Failures return `*QuickCheckoutError` with the exact execution stage:

```go
type QuickCheckoutStage string

const (
	StageParseInput        QuickCheckoutStage = "parse_input"
	StageFetchProduct      QuickCheckoutStage = "fetch_product"
	StageExtractCSRF       QuickCheckoutStage = "extract_csrf"
	StageCheckPayment      QuickCheckoutStage = "check_payment"
	StageFetchShippingForm QuickCheckoutStage = "fetch_shipping_form"
	StageProcessShipment   QuickCheckoutStage = "process_shipment"
)
```
