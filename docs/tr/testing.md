# Test ve Tanılama (Testing & Diagnostics)

SDK, hem yerel birim (unit) testlerinde mock sunucu ile ağ çağrısı yapmadan çalışabilmeniz için `testutil` paketini, hem de canlı Shopier ortamında API anahtarlarınızı ve yetkilerinizi doğrulayabilmeniz için tanı araçlarını sunar.

---

## 1. Çevrimdışı Birim Testleri (`testutil.MockServer`)

SDK, yerel bir `httptest.Server` başlatan ve bu sunucuya bağlanmış bir `*shopier.Client` sağlayan sıfır bağımlılıklı `testutil` paketi içerir:

```go
package myapp_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/AdisGroup/shopier-go/testutil"
)

func TestOrderRetrieval(t *testing.T) {
	server := testutil.NewMockServer()
	defer server.Close()

	// Beklenen uç nokta yanıtını mock sunucuya tanımlayın
	server.HandleFunc("/orders/ord_123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "ord_123",
			"status": "fulfilled",
			"currency": "TRY",
			"totals": {"total": "250.00"},
			"shippingInfo": {
				"firstName": "Ahmet",
				"lastName": "Yilmaz",
				"email": "ahmet@example.com",
				"phone": "+905551234567",
				"address": "Ataturk Cad No:1",
				"city": "Istanbul",
				"country": "Turkiye"
			},
			"lineItems": []
		}`))
	})

	client := server.Client()

	order, err := client.Orders.Get(context.Background(), "ord_123")
	if err != nil {
		t.Fatalf("Orders.Get başarısız: %v", err)
	}

	if order.ID != "ord_123" || order.Status != "fulfilled" {
		t.Errorf("beklenmeyen sipariş verisi: %+v", order)
	}
}
```

### Hata ve Yeniden Deneme (Retry) Senaryoları

Rate limit (429), yetki hatası (403) veya sunucu kesintilerini simüle edebilirsiniz:

```go
// Rate Limit (429) Simülasyonu
server.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", "1")
	w.WriteHeader(http.StatusTooManyRequests)
	_, _ = w.Write([]byte(`{"error":"rate_limit_exceeded","message":"İstek limiti aşıldı"}`))
})

// Yetki Hatası (403) Simülasyonu
server.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":"forbidden","message":"Access to this resource on the server is denied"}`))
})
```

---

## 2. Canlı API Tanı ve Duman Testi (Smoke Test)

Kişisel Erişim Jetonunuzun (PAT), mağaza kimliğinizin ve aktif API yetkilerinizin çalıştığını doğrulamak için aşağıdaki tanı betiğini çalıştırabilirsiniz:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/AdisGroup/shopier-go"
)

func main() {
	token := os.Getenv("SHOPIER_API_KEY")
	if token == "" {
		fmt.Println("Lütfen SHOPIER_API_KEY ortam değişkenini tanımlayın.")
		return
	}

	client, err := shopier.NewClient(token, shopier.WithTimeout(10*time.Second))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	fmt.Println("🔍 Shopier API Tanı Testi Başlatılıyor...")

	// 1. Mağaza Sahibi Doğrulama
	if owner, err := client.Shop.GetOwner(ctx); err != nil {
		fmt.Printf("❌ Shop.GetOwner: %v\n", err)
	} else {
		fmt.Printf("✅ Mağaza Sahibi: %s %s (ID: %s)\n", owner.FirstName, owner.LastName, owner.ID)
	}

	// 2. Mağaza Ayarları Doğrulama
	if settings, err := client.Shop.GetSettings(ctx); err != nil {
		fmt.Printf("❌ Shop.GetSettings: %v\n", err)
	} else {
		fmt.Printf("✅ Mağaza Adı: %s | URL: %s\n", settings.Name, settings.URL)
	}

	// 3. Siparişler Servisi
	if orders, err := client.Orders.List(ctx, nil); err != nil {
		fmt.Printf("❌ Siparişler: %v\n", err)
	} else {
		fmt.Printf("✅ Siparişler: %d adet bulundu\n", len(orders.Items))
	}

	// 4. Bakiye Servisi
	if balances, err := client.Balance.Get(ctx); err != nil {
		fmt.Printf("❌ Bakiye: %v\n", err)
	} else {
		fmt.Printf("✅ Bakiye: %s %s\n", balances[0].Amount, balances[0].Currency)
	}

	// 5. Ürünler Servisi Yetki Kontrolü
	if _, err := client.Products.List(ctx, nil); err != nil {
		fmt.Printf("⚠️ Ürünler: %v (403 Forbidden durumunda hello@shopier.com ile iletişime geçin)\n", err)
	} else {
		fmt.Println("✅ Ürünler: Erişim başarılı")
	}
}
```

---

## 3. CI/CD Otomasyonu

Birim testlerinizi GitHub Actions veya CI sistemlerinizde çalıştırmak için:

```bash
go test -v -race -cover ./...
```

