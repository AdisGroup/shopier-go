# Shopier Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/AdisGroup/shopier-go.svg)](https://pkg.go.dev/github.com/AdisGroup/shopier-go)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Production-ready, zero-dependency Go SDK and API wrapper for the [Shopier REST API](https://developer.shopier.com). Built for high-concurrency microservices, e-commerce backends, and multi-tenant applications.

---

## Features

- **Zero External Dependencies**: Built entirely with Go standard library packages (`net/http`, `crypto/hmac`, `log/slog`, `context`, `iter`, `time`). Zero third-party vulnerability risk and fast compile times.
- **Full API Coverage**: Complete type safety and models across all 12 Shopier REST API resources.
- **Modern Go Iterators**: Native Go 1.23+ Range-over-func iterators (`client.Orders.All`, `client.Products.All`) alongside header-based pagination metadata (`Shopier-Pagination-*`).
- **Resilient Transport**: Automatic exponential backoff with jitter on `5xx` errors and automatic delay adherence on `429 Too Many Requests` using the `Retry-After` header.
- **OAuth 2.0 Ready**: Built-in support for authorization code exchange, token refresh, and token revocation managing Shopier's port `:8443` requirements.
- **Webhook Verification**: Constant-time `HS256` (HMAC-SHA256) signature verification with standard `http.Handler` and multi-framework adapters (Gin, Chi, Fiber).
- **Testing Utilities**: Built-in `testutil.NewMockServer()` for mocking Shopier API interactions in your own unit test suites.

---

## Installation

```bash
go get github.com/AdisGroup/shopier-go
```

---

## Quickstart

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

	// 1. Check account balances
	balances, err := client.Balance.Get(ctx)
	if err != nil {
		log.Fatalf("Balance inquiry failed: %v", err)
	}
	for _, b := range balances {
		fmt.Printf("Balance: %s %s\n", b.Amount, b.Currency)
	}

	// 2. Stream all unfulfilled orders with automatic pagination
	opts := &shopier.OrderListOptions{
		FulfillmentStatus: "unfulfilled",
	}

	for order, err := range client.Orders.All(ctx, opts) {
		if err != nil {
			log.Fatalf("Order stream error: %v", err)
		}
		fmt.Printf("Order #%s: %s %s (Buyer: %s %s)\n",
			order.ID, order.Totals.Total, order.Currency,
			order.ShippingInfo.FirstName, order.ShippingInfo.LastName)
	}
}
```

---

## API Services Overview

| Resource Service | Client Field | Key Methods |
| :--- | :--- | :--- |
| **Orders** | `client.Orders` | `List`, `All`, `Get`, `Update`, `GetTransaction` |
| **Products** | `client.Products` | `List`, `All`, `Get`, `Create`, `Update`, `Delete` |
| **Categories** | `client.Categories` | `List`, `All`, `Get`, `Create`, `Update`, `Delete` |
| **Variations** | `client.Variations` | `List`, `All`, `Get`, `Create`, `Update`, `Delete` |
| **Selections** | `client.Selections` | `List`, `All`, `Get`, `Create`, `Update`, `Delete` |
| **Discounts** | `client.Discounts` | `ListCodes`, `CreateCode`, `ListAutomatic`, `CreateAutomatic` |
| **Shippings** | `client.Shippings` | `List`, `All`, `Get`, `Create`, `Delete` |
| **Balance** | `client.Balance` | `Get`, `ListTransactions`, `AllTransactions`, `GetTransaction` |
| **Payouts** | `client.Payouts` | `List`, `All`, `Get`, `ListTransactions`, `AllTransactions` |
| **Refunds** | `client.Refunds` | `List`, `All`, `Get`, `Create` |
| **Shop** | `client.Shop` | `GetOwner`, `GetSettings`, `UpdateSettings` |
| **Webhooks** | `client.Webhooks` | `List`, `Create`, `Delete` |

---

## Webhook Signature Verification

Shopier signs event payloads with **HS256 (HMAC-SHA256)** via the `Shopier-Signature` header:

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
		log.Printf("Received Event: %s (ID: %s)", evt.Header.Event, evt.Header.WebhookID)

		switch evt.Header.Event {
		case webhook.EventOrderCreated:
			order, err := evt.Order()
			if err != nil {
				return err
			}
			fmt.Printf("New Order Placed: #%s for %s %s\n", order.ID, order.Totals.Total, order.Currency)

		case webhook.EventRefundRequested:
			refund, err := evt.Refund()
			if err != nil {
				return err
			}
			fmt.Printf("Refund Requested: #%s (Order: %s)\n", refund.ID, refund.OrderID)
		}

		return nil
	})

	http.Handle("/webhooks/shopier", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## OAuth 2.0 Web Application Flow

```go
import "github.com/AdisGroup/shopier-go/oauth"

cfg := oauth.NewConfig(clientID, clientSecret, "https://yourapp.com/callback")

// 1. Generate authorization consent URL
consentURL := cfg.AuthCodeURL("state_nonce", oauth.ScopeOrdersRead, oauth.ScopeProductsWrite)

// 2. Exchange authorization code in callback handler
token, err := cfg.Exchange(ctx, authCode)

// 3. Refresh expired token
newToken, err := cfg.RefreshToken(ctx, token.RefreshToken)

// 4. Revoke access
err := cfg.Revoke(ctx, token.AccessToken)
```

---

## Examples

Runnable example code is located in the [`examples/`](examples) directory:

- [`examples/01_quickstart`](examples/01_quickstart): Basic balance inquiry and order listing with PAT.
- [`examples/02_oauth_flow`](examples/02_oauth_flow): Full OAuth2 consent, code exchange, and store inspection.
- [`examples/03_webhook_receiver`](examples/03_webhook_receiver): Webhook receiver with signature verification.
- [`examples/04_iterators`](examples/04_iterators): Streaming items with Go 1.23+ Range-over-func iterators.

---

## Documentation & API Coverage Checklist

Every endpoint, resource model, webhook event, and authentication method specified in the official Shopier API documentation (`ignore/shopier_docs`) is fully implemented, verified, and tested in this SDK:

### 1. Guides & Onboarding (`00_guides/`)
- [x] `00_guides/01_introduction/01_overview.md` &mdash; SDK Architecture & Base Client ([`shopier.go`](shopier.go))
- [x] `00_guides/02_onboarding_checklist/01_sign-up-for-shopier-account.md` &mdash; Account setup reference ([`docs/`](docs/))
- [x] `00_guides/02_onboarding_checklist/02_enable-developer-mode.md` &mdash; Developer mode configuration ([`docs/`](docs/))
- [x] `00_guides/02_onboarding_checklist/03_create-your-app.md` &mdash; OAuth application configuration ([`oauth/oauth.go`](oauth/oauth.go))
- [x] `00_guides/02_onboarding_checklist/04_get-your-app-credentials.md` &mdash; Client ID & Secret credentials handling ([`oauth/oauth.go`](oauth/oauth.go))
- [x] `00_guides/02_onboarding_checklist/05_request-authorization-from-sellers.md` &mdash; Consent URL generator with scopes ([`oauth/oauth.go`](oauth/oauth.go))
- [x] `00_guides/02_onboarding_checklist/06_bonus-shortcut-to-start-playing-with-the-api.md` &mdash; Personal Access Token (PAT) quickstart ([`examples/01_quickstart/main.go`](examples/01_quickstart/main.go))
- [x] `00_guides/03_authentication/01_pats-and-apps.md` &mdash; PAT and OAuth 2.0 dual authentication support ([`shopier.go`](shopier.go), [`oauth/oauth.go`](oauth/oauth.go))
- [x] `00_guides/03_authentication/02_creating-and-using-pats.md` &mdash; HTTP Bearer Token authentication header injection ([`shopier.go`](shopier.go))
- [x] `00_guides/03_authentication/03_oauth-20-flow.md` &mdash; OAuth 2.0 Web Application Flow implementation ([`oauth/oauth.go`](oauth/oauth.go))
- [x] `00_guides/03_authentication/04_oauth-20-flow.md` &mdash; Authorization code exchange & token lifecycle ([`oauth/oauth.go`](oauth/oauth.go))
- [x] `00_guides/03_authentication/05_oauth-20-implementation.md` &mdash; Dedicated OAuth Token & Revoke endpoints on Port `:8443` ([`oauth/oauth.go`](oauth/oauth.go))
- [x] `00_guides/03_authentication/06_permission-scopes.md` &mdash; Complete typed OAuth permission scopes ([`oauth/scopes.go`](oauth/scopes.go))
- [x] `00_guides/04_app_review/01_introduction-to-app-reviews.md` &mdash; App submission & production readiness guidelines ([`docs/`](docs/))
- [x] `00_guides/04_app_review/02_submitting-your-app-for-review.md` &mdash; Review process documentation ([`docs/`](docs/))
- [x] `00_guides/04_app_review/03_publishing-your-app.md` &mdash; App store publishing guidelines ([`docs/`](docs/))

### 2. Shopier REST API Architecture (`01_shopier_rest_api/`)
- [x] `01_shopier_rest_api/01_introduction.md` &mdash; Standard HTTP transport & Client constructor ([`shopier.go`](shopier.go))
- [x] `01_shopier_rest_api/02_response-codes.md` &mdash; Structured error decoding with `APIError` and status handlers ([`errors.go`](errors.go))
- [x] `01_shopier_rest_api/03_pagination.md` &mdash; Header-based pagination parser (`Shopier-Pagination-*`) & `PageResponse[T]` ([`pagination.go`](pagination.go))
- [x] `01_shopier_rest_api/04_sorting.md` &mdash; URL query parameter sorting & ordering helpers ([`options.go`](options.go))
- [x] `01_shopier_rest_api/05_versioning.md` &mdash; API v1 endpoint versioning support ([`shopier.go`](shopier.go))
- [x] `01_shopier_rest_api/06_rate-limits.md` &mdash; Sliding window rate limiter tracking & `Retry-After` retry backoff ([`shopier.go`](shopier.go), [`errors.go`](errors.go))
- [x] `01_shopier_rest_api/07_webhooks.md` &mdash; HS256 HMAC-SHA256 signature verification & payload parsing ([`webhook/verify.go`](webhook/verify.go), [`webhook/events.go`](webhook/events.go))

### 3. API Endpoints (`02_api_endpoints/`)
- **Balance Service**
  - [x] `02_api_endpoints/01_balance/01_balance.md` &mdash; Balance service initialization ([`balance.go`](balance.go))
  - [x] `02_api_endpoints/01_balance/02_the-balance-model.md` &mdash; `Balance` struct model ([`balance.go`](balance.go))
  - [x] `02_api_endpoints/01_balance/03_get-balance.md` &mdash; `client.Balance.Get(ctx)` ([`balance.go`](balance.go))
  - [x] `02_api_endpoints/01_balance/04_the-transaction-model.md` &mdash; `Transaction` struct model ([`balance.go`](balance.go))
  - [x] `02_api_endpoints/01_balance/05_get-balance-transactions.md` &mdash; `client.Balance.ListTransactions(ctx, opts)` & `AllTransactions` iterator ([`balance.go`](balance.go))
  - [x] `02_api_endpoints/01_balance/06_get-balance-transactions-orderid.md` &mdash; `client.Balance.GetTransaction(ctx, orderID)` ([`balance.go`](balance.go))
- **Categories Service**
  - [x] `02_api_endpoints/02_categories/01_categories.md` &mdash; Category service initialization ([`categories.go`](categories.go))
  - [x] `02_api_endpoints/02_categories/02_the-category-model.md` &mdash; `Category` struct model ([`categories.go`](categories.go))
  - [x] `02_api_endpoints/02_categories/03_get-categories.md` &mdash; `client.Categories.List(ctx, opts)` & `All` iterator ([`categories.go`](categories.go))
  - [x] `02_api_endpoints/02_categories/04_post-categories.md` &mdash; `client.Categories.Create(ctx, req)` ([`categories.go`](categories.go))
  - [x] `02_api_endpoints/02_categories/05_get-categories-id.md` &mdash; `client.Categories.Get(ctx, categoryID)` ([`categories.go`](categories.go))
  - [x] `02_api_endpoints/02_categories/06_delete-categories-id.md` &mdash; `client.Categories.Delete(ctx, categoryID)` ([`categories.go`](categories.go))
  - [x] `02_api_endpoints/02_categories/07_put-categories-id.md` &mdash; `client.Categories.Update(ctx, categoryID, req)` ([`categories.go`](categories.go))
- **Discounts Service**
  - [x] `02_api_endpoints/03_discounts/01_discounts.md` &mdash; Discount service initialization ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/02_the-discountcode-model.md` &mdash; `DiscountCode` struct model ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/03_get-discounts-codes.md` &mdash; `client.Discounts.ListCodes(ctx, opts)` & `AllCodes` iterator ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/04_post-discounts-codes.md` &mdash; `client.Discounts.CreateCode(ctx, req)` ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/05_get-discounts-codes-id.md` &mdash; `client.Discounts.GetCode(ctx, codeID)` ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/06_delete-discounts-codes-id.md` &mdash; `client.Discounts.DeleteCode(ctx, codeID)` ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/07_put-discounts-codes-id.md` &mdash; `client.Discounts.UpdateCode(ctx, codeID, req)` ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/08_the-automaticdiscount-model.md` &mdash; `AutomaticDiscount` struct model ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/09_get-discounts-automatic.md` &mdash; `client.Discounts.ListAutomatic(ctx, opts)` & `AllAutomatic` iterator ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/10_post-discounts-automatic.md` &mdash; `client.Discounts.CreateAutomatic(ctx, req)` ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/11_get-discounts-automatic-id.md` &mdash; `client.Discounts.GetAutomatic(ctx, discountID)` ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/12_put-discounts-automatic-id.md` &mdash; `client.Discounts.UpdateAutomatic(ctx, discountID, req)` ([`discounts.go`](discounts.go))
  - [x] `02_api_endpoints/03_discounts/13_delete-discounts-automatic-id.md` &mdash; `client.Discounts.DeleteAutomatic(ctx, discountID)` ([`discounts.go`](discounts.go))
- **Orders Service**
  - [x] `02_api_endpoints/04_orders/01_orders.md` &mdash; Order service initialization ([`orders.go`](orders.go))
  - [x] `02_api_endpoints/04_orders/02_the-order-model.md` &mdash; `Order`, `BillingInfo`, `ShippingInfo`, `Totals` models ([`orders.go`](orders.go))
  - [x] `02_api_endpoints/04_orders/03_get-orders.md` &mdash; `client.Orders.List(ctx, opts)` & `All` iterator ([`orders.go`](orders.go))
  - [x] `02_api_endpoints/04_orders/04_get-orders-id.md` &mdash; `client.Orders.Get(ctx, orderID)` ([`orders.go`](orders.go))
  - [x] `02_api_endpoints/04_orders/05_put-orders-id.md` &mdash; `client.Orders.Update(ctx, orderID, req)` ([`orders.go`](orders.go))
  - [x] `02_api_endpoints/04_orders/06_the-transaction-model-3.md` &mdash; `OrderTransaction` struct model ([`orders.go`](orders.go))
  - [x] `02_api_endpoints/04_orders/07_get-orders-transactions-orderid.md` &mdash; `client.Orders.GetTransaction(ctx, orderID)` ([`orders.go`](orders.go))
- **Payouts Service**
  - [x] `02_api_endpoints/05_payouts/01_payouts.md` &mdash; Payout service initialization ([`payouts.go`](payouts.go))
  - [x] `02_api_endpoints/05_payouts/02_the-payout-model.md` &mdash; `Payout` struct model ([`payouts.go`](payouts.go))
  - [x] `02_api_endpoints/05_payouts/03_get-payouts.md` &mdash; `client.Payouts.List(ctx, opts)` & `All` iterator ([`payouts.go`](payouts.go))
  - [x] `02_api_endpoints/05_payouts/04_get-payouts-id.md` &mdash; `client.Payouts.Get(ctx, payoutID)` ([`payouts.go`](payouts.go))
  - [x] `02_api_endpoints/05_payouts/05_the-transaction-model-2.md` &mdash; `PayoutTransaction` struct model ([`payouts.go`](payouts.go))
  - [x] `02_api_endpoints/05_payouts/06_get-payouts-transactions-id.md` &mdash; `client.Payouts.ListTransactions(ctx, payoutID, opts)` & `AllTransactions` iterator ([`payouts.go`](payouts.go))
- **Products Service**
  - [x] `02_api_endpoints/06_products/01_products.md` &mdash; Product service initialization ([`products.go`](products.go))
  - [x] `02_api_endpoints/06_products/02_the-product-model.md` &mdash; `Product`, `ProductMedia` struct models ([`products.go`](products.go))
  - [x] `02_api_endpoints/06_products/03_get-products.md` &mdash; `client.Products.List(ctx, opts)` & `All` iterator ([`products.go`](products.go))
  - [x] `02_api_endpoints/06_products/04_post-products.md` &mdash; `client.Products.Create(ctx, req)` ([`products.go`](products.go))
  - [x] `02_api_endpoints/06_products/05_get-products-id.md` &mdash; `client.Products.Get(ctx, productID)` ([`products.go`](products.go))
  - [x] `02_api_endpoints/06_products/06_delete-products-id.md` &mdash; `client.Products.Delete(ctx, productID)` ([`products.go`](products.go))
  - [x] `02_api_endpoints/06_products/07_put-products-id.md` &mdash; `client.Products.Update(ctx, productID, req)` ([`products.go`](products.go))
- **Refunds Service**
  - [x] `02_api_endpoints/07_refunds/01_refunds.md` &mdash; Refund service initialization ([`refunds.go`](refunds.go))
  - [x] `02_api_endpoints/07_refunds/02_the-refund-model.md` &mdash; `Refund` struct model ([`refunds.go`](refunds.go))
  - [x] `02_api_endpoints/07_refunds/03_get-refunds.md` &mdash; `client.Refunds.List(ctx, opts)` & `All` iterator ([`refunds.go`](refunds.go))
  - [x] `02_api_endpoints/07_refunds/04_post-refunds.md` &mdash; `client.Refunds.Create(ctx, req)` ([`refunds.go`](refunds.go))
  - [x] `02_api_endpoints/07_refunds/05_get-refunds-id.md` &mdash; `client.Refunds.Get(ctx, refundID)` ([`refunds.go`](refunds.go))
- **Selections Service**
  - [x] `02_api_endpoints/08_selections/01_selections.md` &mdash; Selection service initialization ([`selections.go`](selections.go))
  - [x] `02_api_endpoints/08_selections/02_the-selection-model.md` &mdash; `Selection`, `SelectionOption` struct models ([`selections.go`](selections.go))
  - [x] `02_api_endpoints/08_selections/03_get-selections.md` &mdash; `client.Selections.List(ctx, opts)` & `All` iterator ([`selections.go`](selections.go))
  - [x] `02_api_endpoints/08_selections/04_post-selections.md` &mdash; `client.Selections.Create(ctx, req)` ([`selections.go`](selections.go))
  - [x] `02_api_endpoints/08_selections/05_get-selections-id.md` &mdash; `client.Selections.Get(ctx, selectionID)` ([`selections.go`](selections.go))
  - [x] `02_api_endpoints/08_selections/06_put-selections-id.md` &mdash; `client.Selections.Update(ctx, selectionID, req)` ([`selections.go`](selections.go))
  - [x] `02_api_endpoints/08_selections/07_delete-selections-id.md` &mdash; `client.Selections.Delete(ctx, selectionID)` ([`selections.go`](selections.go))
- **Shippings Service**
  - [x] `02_api_endpoints/09_shippings/01_shippings.md` &mdash; Shipping service initialization ([`shippings.go`](shippings.go))
  - [x] `02_api_endpoints/09_shippings/02_the-shipping-model.md` &mdash; `Shipping`, `ShippingCountry` struct models ([`shippings.go`](shippings.go))
  - [x] `02_api_endpoints/09_shippings/03_get-shippings.md` &mdash; `client.Shippings.List(ctx, opts)` & `All` iterator ([`shippings.go`](shippings.go))
  - [x] `02_api_endpoints/09_shippings/04_post-shippings.md` &mdash; `client.Shippings.Create(ctx, req)` ([`shippings.go`](shippings.go))
  - [x] `02_api_endpoints/09_shippings/05_get-shippings-code.md` &mdash; `client.Shippings.Get(ctx, shippingCode)` ([`shippings.go`](shippings.go))
  - [x] `02_api_endpoints/09_shippings/06_delete-shippings-code.md` &mdash; `client.Shippings.Delete(ctx, shippingCode)` ([`shippings.go`](shippings.go))
- **Shop Service**
  - [x] `02_api_endpoints/10_shop/01_shop.md` &mdash; Shop service initialization ([`shop.go`](shop.go))
  - [x] `02_api_endpoints/10_shop/02_the-shopowner-model.md` &mdash; `ShopOwner` struct model ([`shop.go`](shop.go))
  - [x] `02_api_endpoints/10_shop/03_get-shop-owner.md` &mdash; `client.Shop.GetOwner(ctx)` ([`shop.go`](shop.go))
  - [x] `02_api_endpoints/10_shop/04_the-shopsetting-model.md` &mdash; `ShopSetting` struct model ([`shop.go`](shop.go))
  - [x] `02_api_endpoints/10_shop/05_get-shop-settings.md` &mdash; `client.Shop.GetSettings(ctx)` ([`shop.go`](shop.go))
  - [x] `02_api_endpoints/10_shop/06_put-shop-settings.md` &mdash; `client.Shop.UpdateSettings(ctx, req)` ([`shop.go`](shop.go))
- **Variations Service**
  - [x] `02_api_endpoints/11_variations/01_variations.md` &mdash; Variation service initialization ([`variations.go`](variations.go))
  - [x] `02_api_endpoints/11_variations/02_the-variation-model.md` &mdash; `Variation`, `VariationOption` struct models ([`variations.go`](variations.go))
  - [x] `02_api_endpoints/11_variations/03_get-variations.md` &mdash; `client.Variations.List(ctx, opts)` & `All` iterator ([`variations.go`](variations.go))
  - [x] `02_api_endpoints/11_variations/04_post-variations.md` &mdash; `client.Variations.Create(ctx, req)` ([`variations.go`](variations.go))
  - [x] `02_api_endpoints/11_variations/05_get-variations-id.md` &mdash; `client.Variations.Get(ctx, variationID)` ([`variations.go`](variations.go))
  - [x] `02_api_endpoints/11_variations/06_delete-variations-id.md` &mdash; `client.Variations.Delete(ctx, variationID)` ([`variations.go`](variations.go))
  - [x] `02_api_endpoints/11_variations/07_put-variations-id.md` &mdash; `client.Variations.Update(ctx, variationID, req)` ([`variations.go`](variations.go))
- **Webhooks Service**
  - [x] `02_api_endpoints/12_webhooks/01_the-webhook-model.md` &mdash; `Webhook` struct model ([`webhooks.go`](webhooks.go))
  - [x] `02_api_endpoints/12_webhooks/02_get-webhooks.md` &mdash; `client.Webhooks.List(ctx, opts)` & `All` iterator ([`webhooks.go`](webhooks.go))
  - [x] `02_api_endpoints/12_webhooks/03_post-webhooks.md` &mdash; `client.Webhooks.Create(ctx, req)` ([`webhooks.go`](webhooks.go))
  - [x] `02_api_endpoints/12_webhooks/04_delete-webhooks-id.md` &mdash; `client.Webhooks.Delete(ctx, webhookID)` ([`webhooks.go`](webhooks.go))

### 4. Webhooks Subsystem (`03_webhooks/`)
- [x] `03_webhooks/01_webhook-configuration.md` &mdash; Webhook subscription management & `http.Handler` integration ([`webhook/handler.go`](webhook/handler.go), [`webhooks.go`](webhooks.go))
- [x] `03_webhooks/02_events-headers-payloads.md` &mdash; Header extraction, HS256 HMAC verification, and Event unmarshaling ([`webhook/verify.go`](webhook/verify.go), [`webhook/events.go`](webhook/events.go))

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

