# Kargo İşlemleri (Shippings API)

Kargo servisi, Shopier Anlaşmalı Kargo entegrasyonu üzerinden indirimli kargo kodları üretmenizi ve gönderi durumlarını sorgulamanızı sağlar.

## Metotlar

- `client.Shippings.List(ctx, opts)` -> `(*PageResponse[Shipping], error)`
- `client.Shippings.All(ctx, opts)` -> `iter.Seq2[Shipping, error]`
- `client.Shippings.Get(ctx, code)` -> `(*Shipping, error)`
- `client.Shippings.Create(ctx, req)` -> `(*Shipping, error)`
- `client.Shippings.Delete(ctx, code)` -> `error`

---

## 1. Anlaşmalı Kargo Kodu Üretme

```go
shipping, err := client.Shippings.Create(ctx, &shopier.ShippingCreateRequest{
	OrderID: "ord_123456",
	Company: "yurtici", // yurtici, mng, ptt, aras, surat, ups vb.
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Kargo Kodu: %s (Firma: %s)\n", shipping.Code, shipping.Company)
```

## 2. Kargo Durumu Sorgulama

```go
ship, err := client.Shippings.Get(ctx, "SHP-987654")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Durum: %s | Takip No: %s (Takip Linki: %s)\n",
	ship.Status, ship.TrackingNumber, ship.TrackingURL)
```
