# Shippings API Reference

The Shippings API manages Shopier Contracted Shipping labels, barcode tracking codes, and carrier dispatches.

## Methods

- `client.Shippings.List(ctx, opts)` -> `(*PageResponse[Shipping], error)`
- `client.Shippings.All(ctx, opts)` -> `iter.Seq2[Shipping, error]`
- `client.Shippings.Get(ctx, code)` -> `(*Shipping, error)`
- `client.Shippings.Create(ctx, req)` -> `(*Shipping, error)`
- `client.Shippings.Delete(ctx, code)` -> `error`

---

## 1. Generate Contracted Shipping Code

```go
shipping, err := client.Shippings.Create(ctx, &shopier.ShippingCreateRequest{
	OrderID: "ord_123456",
	Company: "yurtici", // yurtici, mng, ptt, aras, surat, ups, etc.
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Shipping Code: %s (Company: %s)\n", shipping.Code, shipping.Company)
```

## 2. Get Shipping Status

```go
ship, err := client.Shippings.Get(ctx, "SHP-987654")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Status: %s | Tracking Number: %s (URL: %s)\n",
	ship.Status, ship.TrackingNumber, ship.TrackingURL)
```
