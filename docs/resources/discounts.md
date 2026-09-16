# Discounts API Reference

The Discounts API manages both promotional coupon codes and automated discount campaigns.

## Methods

### Discount Codes
- `client.Discounts.ListCodes(ctx, opts)` -> `(*PageResponse[DiscountCode], error)`
- `client.Discounts.AllCodes(ctx, opts)` -> `iter.Seq2[DiscountCode, error]`
- `client.Discounts.GetCode(ctx, id)` -> `(*DiscountCode, error)`
- `client.Discounts.CreateCode(ctx, req)` -> `(*DiscountCode, error)`
- `client.Discounts.UpdateCode(ctx, id, req)` -> `(*DiscountCode, error)`
- `client.Discounts.DeleteCode(ctx, id)` -> `error`

### Automatic Discounts
- `client.Discounts.ListAutomatic(ctx, opts)` -> `(*PageResponse[AutomaticDiscount], error)`
- `client.Discounts.AllAutomatic(ctx, opts)` -> `iter.Seq2[AutomaticDiscount, error]`
- `client.Discounts.GetAutomatic(ctx, id)` -> `(*AutomaticDiscount, error)`
- `client.Discounts.CreateAutomatic(ctx, req)` -> `(*AutomaticDiscount, error)`
- `client.Discounts.UpdateAutomatic(ctx, id, req)` -> `(*AutomaticDiscount, error)`
- `client.Discounts.DeleteAutomatic(ctx, id)` -> `error`

---

## 1. Create a Discount Code

```go
code, err := client.Discounts.CreateCode(ctx, &shopier.DiscountCodeCreateRequest{
	Code:          "WELCOME20",
	Type:          "percent", // "percent" or "amount"
	PercentOff:    "20",
	AmountMinimum: "100.00",
	Currency:      "TRY",
	NumAvailable:  500,
	ExpiresAt:     "2026-12-31+0300",
})
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Created Coupon: %s (ID: %s)\n", code.Code, code.ID)
```

## 2. Create an Automatic Discount

```go
autoDisc, err := client.Discounts.CreateAutomatic(ctx, &shopier.AutomaticDiscountCreateRequest{
	Title:         "Summer Sale 15% Off",
	Scope:         "all", // "all", "selectedProducts", "selectedCategories"
	Type:          "percent",
	PercentOff:    "15",
	Requirement:   "amount", // "amount" or "quantity"
	AmountMinimum: "250.00",
	Currency:      "TRY",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Created Campaign ID:", autoDisc.ID)
```
