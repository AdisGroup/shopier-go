# Test ve Mocking (Testing)

SDK, uygulamanızın birim (unit) ve entegrasyon testlerinde gerçek Shopier sunucularına istek atmadan senaryoları simüle edebilmeniz için dahili `testutil` paketi sunar.

---

## `testutil.NewMockServer` Kullanımı

`MockServer`, yerel bir `httptest.Server` başlatır ve bu sunucuya önceden bağlanmış bir `*shopier.Client` üretir:

```go
package myapp_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/AdisGroup/shopier-go/testutil"
)

func TestMyOrderWorker(t *testing.T) {
	server := testutil.NewMockServer()
	defer server.Close()

	// Beklenen uç nokta yanıtını tanımlayın
	server.HandleFunc("/orders/ord_test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "ord_test",
			"status": "fulfilled",
			"currency": "TRY",
			"totals": {"total": "120.00"},
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

	// Mock sunucuya ayarlanmış istemciyi alın
	client := server.Client()

	order, err := client.Orders.Get(context.Background(), "ord_test")
	if err != nil {
		t.Fatalf("Orders.Get başarısız: %v", err)
	}

	if order.Status != "fulfilled" {
		t.Errorf("beklenen durum: fulfilled, alınan: %s", order.Status)
	}
}
```

---

## Hata ve Retry Senaryolarını Simüle Etme

Rate limit (429), doğrulama hataları (400) veya geçici sunucu kesintilerini (503) kolayca test edebilirsiniz:

```go
server.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(`{"error":"invalid_title","message":"Ürün başlığı zorunludur"}`))
})
```
