# Refunds API Reference

The Refunds API allows you to issue full or partial reimbursements for paid orders.

## Methods

- `client.Refunds.List(ctx, opts)` -> `(*PageResponse[Refund], error)`
- `client.Refunds.All(ctx, opts)` -> `iter.Seq2[Refund, error]`
- `client.Refunds.Get(ctx, id)` -> `(*Refund, error)`
- `client.Refunds.Create(ctx, req)` -> `(*Refund, error)`

---

## 1. Issue a Full Refund

```go
refund, err := client.Refunds.Create(ctx, &shopier.RefundCreateRequest{
	OrderID: "ord_123456",
	Note:    "Customer returned merchandise within 14 days",
})
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Refund ID: %s | Status: %s\n", refund.ID, refund.Status)
```

## 2. Issue a Partial Refund

```go
partialRefund, err := client.Refunds.Create(ctx, &shopier.RefundCreateRequest{
	OrderID: "ord_123456",
	Amount:  "50.00",
	Note:    "Partial refund for damaged outer box",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Partial Refund ID:", partialRefund.ID)
```

---

## Data Models & Types

### `Refund` {#refund-model}

Represents a monetary reimbursement issued on an order.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Unique refund identifier |
| `OrderID` | `string` | `orderId` | Associated Shopier order ID |
| `Type` | `string` | `type` | `"full"` or `"partial"` |
| `Status` | `string` | `status` | `"pending"`, `"succeeded"`, or `"failed"` |
| `Currency` | `string` | `currency` | 3-letter currency code (`"TRY"`, `"USD"`, etc.) |
| `Total` | `string` | `total` | Total refunded amount |
| `Note` | `string` | `note` | Reason or explanatory note |
| `DateCreated` | `string` | `dateCreated` | ISO-8601 creation timestamp |
| `DateRefunded` | `string` | `dateRefunded` | Timestamp when reimbursement settled |

---

### `RefundCreateRequest`

| Field | Type | JSON Tag | Required | Description |
| :--- | :--- | :--- | :---: | :--- |
| `OrderID` | `string` | `orderId` | **Yes** | Shopier order ID to refund |
| `Amount` | `string` | `amount` | No | Amount for partial refund (omit for 100% full refund) |
| `Note` | `string` | `note` | No | Internal or customer-visible refund reason |

---

### `RefundListOptions`

| Field | Type | URL Param | Description |
| :--- | :--- | :--- | :--- |
| `OrderID` | `string` | `orderId` | Filter refunds by specific order ID |
| `Page` | `int` | `page` | Pagination page index (1-based) |
| `Limit` | `int` | `limit` | Records per page (default: 20) |

