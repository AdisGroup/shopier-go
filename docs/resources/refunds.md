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
