# Sayfalama ve İteratörler (Pagination)

Shopier REST API tüm liste uç noktalarında sayfalama kullanır. Sayfalama meta verileri HTTP yanıt başlıklarında (headers) döner:
- `Shopier-Pagination-Page`: Mevcut sayfa numarası
- `Shopier-Pagination-Limit`: Sayfa başına kayıt sayısı
- `Shopier-Pagination-Total-Pages`: Toplam sayfa adedi
- `Shopier-Pagination-Total-Items`: Toplam kayıt adedi

SDK, sayfalanmış verilerle çalışmak için iki farklı yöntem sunar:
1. **Go 1.23+ Range-Over-Func İteratörleri (`All`)**: Native `for ... range` döngüsü ile sayfaları otomatik çeken modern yaklaşım.
2. **Manuel Sayfa İstekleri (`List`)**: Sayfa numarası ve limit üzerinde tam kontrol sağlayan klasik yaklaşım.

---

## 1. Modern Go 1.23+ Range İteratörleri (`All`)

`All` metodu `iter.Seq2[T, error]` döndürür. Döngü mevcut sayfanın sonuna geldikçe sonraki sayfayı arka planda otomatik olarak ister:

```go
ctx := context.Background()

// Tüm siparişleri sayfalama karmaşası olmadan tek bir döngüde dolaşın
for order, err := range client.Orders.All(ctx, nil) {
	if err != nil {
		log.Fatalf("İterasyon hatası: %v", err)
	}

	fmt.Printf("Sipariş #%s: %s %s\n", order.ID, order.Totals.Total, order.Currency)
}
```

Filtre parametrelerini `All` metoduna da iletebilirsiniz:

```go
opts := &shopier.OrderListOptions{
	FulfillmentStatus: "unfulfilled",
}

for order, err := range client.Orders.All(ctx, opts) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Kargolanacak sipariş:", order.ID)
}
```

---

## 2. Manuel Sayfa İstekleri (`List`)

Ön yüzde sayfalama (pagination) düğmeleri oluştururken veya belirli sayfaları çekerken:

```go
opts := &shopier.OrderListOptions{
	ListOptions: shopier.ListOptions{
		Page:  1,
		Limit: 20, // Min: 1, Max: 50, Varsayılan: 10
		Sort:  "desc",
	},
	FulfillmentStatus: "unfulfilled",
}

res, err := client.Orders.List(ctx, opts)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Mevcut Sayfa: %d / %d (Toplam Kayıt: %d)\n",
	res.Pagination.Page,
	res.Pagination.TotalPages,
	res.Pagination.TotalItems,
)

for _, order := range res.Items {
	fmt.Printf("- Sipariş: %s\n", order.ID)
}

// Sonraki sayfa kontrolü
if res.Pagination.HasNextPage() {
	opts.Page = res.Pagination.NextPage()
	// Sonraki sayfayı isteyin...
}
```
