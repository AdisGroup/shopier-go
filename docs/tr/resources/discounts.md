# İndirimler (Discounts API)

İndirim servisi, hem sepet aşamasında girilen kupon kodlarını hem de sepet tutarına/adedine göre kendiliğinden devreye giren otomatik indirim kampanyalarını yönetir.

## Metotlar

### İndirim Kodları (Discount Codes)
- `client.Discounts.ListCodes(ctx, opts)` -> `(*PageResponse[DiscountCode], error)`
- `client.Discounts.AllCodes(ctx, opts)` -> `iter.Seq2[DiscountCode, error]`
- `client.Discounts.GetCode(ctx, id)` -> `(*DiscountCode, error)`
- `client.Discounts.CreateCode(ctx, req)` -> `(*DiscountCode, error)`
- `client.Discounts.UpdateCode(ctx, id, req)` -> `(*DiscountCode, error)`
- `client.Discounts.DeleteCode(ctx, id)` -> `error`

### Otomatik İndirimler (Automatic Discounts)
- `client.Discounts.ListAutomatic(ctx, opts)` -> `(*PageResponse[AutomaticDiscount], error)`
- `client.Discounts.AllAutomatic(ctx, opts)` -> `iter.Seq2[AutomaticDiscount, error]`
- `client.Discounts.GetAutomatic(ctx, id)` -> `(*AutomaticDiscount, error)`
- `client.Discounts.CreateAutomatic(ctx, req)` -> `(*AutomaticDiscount, error)`
- `client.Discounts.UpdateAutomatic(ctx, id, req)` -> `(*AutomaticDiscount, error)`
- `client.Discounts.DeleteAutomatic(ctx, id)` -> `error`

---

## 1. İndirim Kodu Oluşturma

```go
code, err := client.Discounts.CreateCode(ctx, &shopier.DiscountCodeCreateRequest{
	Code:          "HOSGELDIN20",
	Type:          "percent", // "percent" veya "amount"
	PercentOff:    "20",
	AmountMinimum: "100.00",
	Currency:      "TRY",
	NumAvailable:  500,
	ExpiresAt:     "2026-12-31+0300",
})
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Kupon Oluşturuldu: %s (ID: %s)\n", code.Code, code.ID)
```

## 2. Otomatik İndirim Kampanyası Tanımlama

```go
autoDisc, err := client.Discounts.CreateAutomatic(ctx, &shopier.AutomaticDiscountCreateRequest{
	Title:         "250 TL Üzeri %15 İndirim",
	Scope:         "all", // "all", "selectedProducts", "selectedCategories"
	Type:          "percent",
	PercentOff:    "15",
	Requirement:   "amount", // "amount" veya "quantity"
	AmountMinimum: "250.00",
	Currency:      "TRY",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Kampanya ID:", autoDisc.ID)
```
