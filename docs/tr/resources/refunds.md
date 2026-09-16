# İadeler (Refunds API)

İade servisi, ödenmiş siparişler için tam veya kısmi tutarda iade işlemi başlatmanızı ve iade süreçlerini takip etmenizi sağlar.

## Metotlar

- `client.Refunds.List(ctx, opts)` -> `(*PageResponse[Refund], error)`
- `client.Refunds.All(ctx, opts)` -> `iter.Seq2[Refund, error]`
- `client.Refunds.Get(ctx, id)` -> `(*Refund, error)`
- `client.Refunds.Create(ctx, req)` -> `(*Refund, error)`

---

## 1. Tam İade Oluşturma

```go
refund, err := client.Refunds.Create(ctx, &shopier.RefundCreateRequest{
	OrderID: "ord_123456",
	Note:    "Müşteri cayma hakkı kapsamında iade edildi",
})
if err != nil {
	log.Fatal(err)
}
fmt.Printf("İade ID: %s | Durum: %s\n", refund.ID, refund.Status)
```

## 2. Kısmi İade Oluşturma

```go
partialRefund, err := client.Refunds.Create(ctx, &shopier.RefundCreateRequest{
	OrderID: "ord_123456",
	Amount:  "50.00",
	Note:    "Kargo hasarı sebebiyle kısmi iade",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Kısmi İade ID:", partialRefund.ID)
```
