# Orders API Reference

The Orders API allows you to retrieve, filter, and update orders placed through Shopier storefronts.

## Methods

- `client.Orders.List(ctx, opts)` -> `(*PageResponse[Order], error)`
- `client.Orders.All(ctx, opts)` -> `iter.Seq2[Order, error]`
- `client.Orders.Get(ctx, id)` -> `(*Order, error)`
- `client.Orders.Update(ctx, id, req)` -> `(*Order, error)`
- `client.Orders.GetTransaction(ctx, orderID)` -> `(*OrderTransaction, error)`

---

## 1. List Orders

```go
opts := &shopier.OrderListOptions{
	ListOptions: shopier.ListOptions{
		Page:  1,
		Limit: 20,
	},
	FulfillmentStatus: "unfulfilled", // "unfulfilled" or "fulfilled"
	RefundType:        "none",        // "none", "partial", "full"
	DateStart:         "2026-01-01T00:00:00Z",
}

res, err := client.Orders.List(ctx, opts)
if err != nil {
	log.Fatal(err)
}

for _, ord := range res.Items {
	fmt.Printf("Order #%s: Total %s %s\n", ord.ID, ord.Totals.Total, ord.Currency)
}
```

## 2. Get Order

```go
order, err := client.Orders.Get(ctx, "123456789")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Buyer: %s %s (%s)\n", order.ShippingInfo.FirstName, order.ShippingInfo.LastName, order.ShippingInfo.Email)
```

## 3. Update Order (Fulfill & Track)

```go
updated, err := client.Orders.Update(ctx, "123456789", &shopier.OrderUpdateRequest{
	Fulfillments: &shopier.OrderFulfillments{
		ShippingCompany: "yurtici", // yurtici, mng, ptt, aras, surat, ups, dhl, etc.
		TrackingNumber:  "9876543210",
	},
})
if err != nil {
	log.Fatal(err)
}

fmt.Println("New status:", updated.Status) // "fulfilled"
```

## 4. Get Order Transaction

```go
tx, err := client.Orders.GetTransaction(ctx, "123456789")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Gross: %s %s | Fee: %s %s | Net: %s %s\n",
	tx.Gross.Amount, tx.Gross.Currency,
	tx.Fee.Amount, tx.Fee.Currency,
	tx.Net.Amount, tx.Net.Currency,
)
```
