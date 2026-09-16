# Webhooks & Signature Verification

Shopier dispatches event notifications via HTTP POST to your configured Notification URL. To guarantee authenticity and prevent tampering, every webhook payload is signed with **HS256 (HMAC-SHA256)** using your app's webhook secret token.

---

## Webhook Headers

| Header | Description |
| :--- | :--- |
| `Shopier-Signature` | HMAC-SHA256 signature string (hex or base64) |
| `Shopier-Event` | Event type identifier (e.g., `order.created`) |
| `Shopier-Webhook-Id` | Unique ID of the notification event |
| `Shopier-Timestamp` | Unix timestamp in seconds (UTC) |
| `Shopier-Account-Id` | The shop account ID associated with the event |
| `Shopier-Api-Version` | API version string |

---

## Supported Events

| Event Constant | Event Name | Payload Model | Trigger Scenario |
| :--- | :--- | :--- | :--- |
| `webhook.EventOrderCreated` | `order.created` | `*shopier.Order` | A new order is paid and placed |
| `webhook.EventOrderAddressUpdated` | `order.addressUpdated` | `*shopier.Order` | Buyer shipping address is updated |
| `webhook.EventOrderFulfilled` | `order.fulfilled` | `*shopier.Order` | Order is closed and marked shipped |
| `webhook.EventProductCreated` | `product.created` | `*shopier.Product` | A new product listing is published |
| `webhook.EventProductUpdated` | `product.updated` | `*shopier.Product` | An existing product is updated |
| `webhook.EventRefundRequested` | `refund.requested` | `*shopier.Refund` | A seller requests an order refund |
| `webhook.EventRefundUpdated` | `refund.updated` | `*shopier.Refund` | Refund succeeds or fails |

---

## 1. Using the Built-in `http.Handler` Adapter

The fastest way to mount a production-ready webhook endpoint in Go:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AdisGroup/shopier-go/webhook"
)

func main() {
	secret := os.Getenv("SHOPIER_WEBHOOK_SECRET")

	handler := webhook.NewHandler(secret, func(ctx context.Context, evt *webhook.Event) error {
		log.Printf("Received event: %s (ID: %s)", evt.Header.Event, evt.Header.WebhookID)

		switch evt.Header.Event {
		case webhook.EventOrderCreated:
			order, err := evt.Order()
			if err != nil {
				return err
			}
			fmt.Printf("New Order #%s for %s %s\n", order.ID, order.Totals.Total, order.Currency)

		case webhook.EventRefundRequested:
			refund, err := evt.Refund()
			if err != nil {
				return err
			}
			fmt.Printf("Refund requested: #%s (Order: %s)\n", refund.ID, refund.OrderID)
		}

		return nil
	})

	http.Handle("/webhooks/shopier", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 2. Manual Verification in Web Frameworks (Gin, Chi, Fiber)

### Standard `webhook.Parse`

You can verify and parse any `*http.Request` directly:

```go
event, err := webhook.Parse(r, webhookSecret)
if err != nil {
	http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
	return
}

order, _ := event.Order()
```

### Direct Byte Verification (`webhook.VerifySignature`)

For frameworks like Fiber that do not expose standard `*http.Request`:

```go
// Fiber example
app.Post("/webhooks/shopier", func(c *fiber.Ctx) error {
	sig := c.Get("Shopier-Signature")
	body := c.Body()

	if !webhook.VerifySignature(body, sig, webhookSecret) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid signature")
	}

	// Process payload...
	return c.SendStatus(fiber.StatusOK)
})
```
