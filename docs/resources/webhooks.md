# Webhook Subscriptions API Reference

The Webhook Subscriptions API allows programmatic registration and management of webhook event notification URLs.

> [!NOTE]
> For verifying incoming webhook payloads, see the [Webhooks & Signature Verification Guide](/webhooks).

## Methods

- `client.Webhooks.List(ctx)` -> `([]WebhookSubscription, error)`
- `client.Webhooks.Create(ctx, req)` -> `(*WebhookSubscription, error)`
- `client.Webhooks.Delete(ctx, id)` -> `error`

---

## 1. Register a Webhook Subscription

```go
sub, err := client.Webhooks.Create(ctx, &shopier.WebhookCreateRequest{
	Event: "order.created",
	URL:   "https://api.mybrand.com/webhooks/shopier",
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Subscription Created! ID: %s\n", sub.ID)
// IMPORTANT: Save sub.Token securely. It is only returned upon initial creation.
fmt.Printf("Signing Secret Token: %s\n", sub.Token)
```

## 2. List Active Subscriptions

```go
subs, err := client.Webhooks.List(ctx)
if err != nil {
	log.Fatal(err)
}

for _, s := range subs {
	fmt.Printf("[%s] Event: %s -> %s\n", s.ID, s.Event, s.URL)
}
```

## 3. Delete a Webhook Subscription

```go
err := client.Webhooks.Delete(ctx, "wh_12345")
if err != nil {
	log.Fatal(err)
}
```

---

## Data Models & Types

### `WebhookSubscription` {#webhooksubscription-model}

Represents a configured callback notification endpoint.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Unique webhook subscription ID |
| `Event` | `string` | `event` | Subscribed event name (e.g. `"order.created"`) |
| `URL` | `string` | `url` | Destination HTTPS callback URL |
| `Token` | `string` | `token` | HMAC secret signing key (*only returned on initial creation*) |

---

### `WebhookCreateRequest`

| Field | Type | JSON Tag | Required | Description |
| :--- | :--- | :--- | :---: | :--- |
| `Event` | `string` | `event` | **Yes** | Event name to subscribe to |
| `URL` | `string` | `url` | **Yes** | HTTPS target destination URL |

