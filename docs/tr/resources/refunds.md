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

---

## Veri Modelleri ve Tipler

### `Refund` (İade Modeli) {#iade-modeli-refund}

Bir sipariş için gerçekleştirilen geri ödeme işlemini temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Benzersiz iade işlem numarası |
| `OrderID` | `string` | `orderId` | İade yapılan Shopier sipariş numarası |
| `Type` | `string` | `type` | `"full"` (tam iade) veya `"partial"` (kısmi iade) |
| `Status` | `string` | `status` | İade durumu (`"pending"`, `"succeeded"`, `"failed"`) |
| `Currency` | `string` | `currency` | Para birimi kodu (`"TRY"`, `"USD"` vb.) |
| `Total` | `string` | `total` | İade edilen toplam tutar |
| `Note` | `string` | `note` | İade gerekçesi veya açıklama notu |
| `DateCreated` | `string` | `dateCreated` | İade talebinin oluşturulma tarihi |
| `DateRefunded` | `string` | `dateRefunded` | İadenin banka tarafından tamamlanma tarihi |

---

### `RefundCreateRequest`

| Alan (Field) | Tip | JSON Etiketi | Zorunlu | Açıklama |
| :--- | :--- | :--- | :---: | :--- |
| `OrderID` | `string` | `orderId` | **Evet** | İade yapılacak Shopier sipariş ID |
| `Amount` | `string` | `amount` | Hayır | Kısmi iade tutarı (belirtilmezse siparişin tamamı iade edilir) |
| `Note` | `string` | `note` | Hayır | İade işlemine ait açıklama notu |

---

### `RefundListOptions`

| Alan (Field) | Tip | URL Parametresi | Açıklama |
| :--- | :--- | :--- | :--- |
| `OrderID` | `string` | `orderId` | Belirli bir siparişe ait iadeleri filtreler |
| `Page` | `int` | `page` | Sayfa indeksi (1 tabanlı) |
| `Limit` | `int` | `limit` | Sayfa başına kayıt sayısı (varsayılan: 20) |

