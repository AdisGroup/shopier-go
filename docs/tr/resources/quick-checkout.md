# Hızlı Ödeme Linki (Quick Checkout)

`QuickCheckoutService`, alıcı adına Shopier sepet/teslimat adımlarını uçtan uca işleterek doğrudan ödeme sayfasına yönlendiren bir Shopier ödeme linki (`https://www.shopier.com/s/payment/{account}/{order_id}`) üretir.

Bu servis, `shopier.Client` üzerindeki `client.QuickCheckout` alanı üzerinden doğrudan veya API token'ına ihtiyaç duyulmadan `shopier.NewQuickCheckoutService()` kurucusu ile bağımsız olarak kullanılabilir.

---

## Eşzamanlılık ve Oturum İzolasyonu

Her `Create` çağrısı, kendi izole çerez havuzunu (`net/http/cookiejar`) ve HTTP oturumunu oluşturur. Bu sayede eşzamanlı (concurrent) çalışan goroutine'ler arasında oturum, çerez veya CSRF token çakışması kesinlikle yaşanmaz.

---

## Hızlı Başlangıç Örneği

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AdisGroup/shopier-go"
)

func main() {
	qc := shopier.NewQuickCheckoutService(
		shopier.WithQuickCheckoutTimeout(30 * time.Second),
	)

	ctx := context.Background()

	req := &shopier.QuickCheckoutRequest{
		ProductURL: "https://www.shopier.com/adisgroup/50900391",
		Quantity:   1,
		Buyer: shopier.QuickCheckoutBuyer{
			Email:     "alici@example.com",
			FirstName: "Özgün Deniz",
			LastName:  "Küçük",
			Phone:     "+90 532 477 02 10",
			Country:   "Türkiye",
			City:      "İstanbul",
			Address:   "Merkez Mah. No: 1",
			Comment:   "Lütfen mesai saatlerinde teslim ediniz.",
		},
	}

	result, err := qc.Create(ctx, req)
	if err != nil {
		var qErr *shopier.QuickCheckoutError
		if qErr, ok := err.(*shopier.QuickCheckoutError); ok {
			log.Fatalf("Hızlı ödeme oluşturulamadı [%s]: HTTP %d - %s",
				qErr.Stage, qErr.StatusCode, qErr.Message)
		}
		log.Fatalf("Hata: %v", err)
	}

	fmt.Printf("Sipariş ID : %s\n", result.OrderID)
	fmt.Printf("Ödeme Linki: %s\n", result.PaymentURL)
}
```

---

## İstek Parametreleri

### `QuickCheckoutRequest`

| Alan | Tip | Açıklama |
| :--- | :--- | :--- |
| `ProductURL` | `string` | Ürün linki (ör. `https://www.shopier.com/adisgroup/50900391`). `Account` ve `ProductID` otomatik ayrıştırılır. |
| `Account` | `string` | Mağaza kullanıcı adı / slug (ör. `adisgroup`). |
| `ProductID` | `string` | Ürün kimlik numarası (ör. `50900391`). |
| `Quantity` | `int` | Satın alınacak adet. Belirtilmezse varsayılan `1` alınır. |
| `Options` | `map[string]string` | Varyasyon veya seçenek form verileri. |
| `Buyer` | `QuickCheckoutBuyer` | Alıcı iletişim ve teslimat bilgileri. |

### `QuickCheckoutBuyer`

| Alan | Tip | Açıklama |
| :--- | :--- | :--- |
| `Email` | `string` | **Zorunlu.** Alıcının e-posta adresi. |
| `FirstName` | `string` | **Zorunlu.** Alıcının adı. |
| `LastName` | `string` | **Zorunlu.** Alıcının soyadı. |
| `Phone` | `string` | **Zorunlu.** Telefon numarası (ör. `+90 532 477 02 10`). |
| `PhoneCode` | `string` | Ülke arama kodu (varsayılan: `"TR"`). |
| `Country` | `string` | Teslimat ülkesi (varsayılan: `"Türkiye"`). |
| `City` | `string` | Teslimat ili. |
| `State` | `string` | Teslimat ilçesi. |
| `Address` | `string` | Açık adres. |
| `ZipCode` | `string` | Posta kodu. |
| `TCIDNo` | `string` | T.C. Kimlik No (opsiyonel). |
| `Comment` | `string` | Sipariş/teslimat notu. |

---

## Yanıt Yapısı

### `QuickCheckoutResult`

```go
type QuickCheckoutResult struct {
	OrderID     string         `json:"order_id"`
	PaymentURL  string         `json:"payment_url"`
	Account     string         `json:"account"`
	ProductID   string         `json:"product_id"`
	RawResponse map[string]any `json:"raw_response,omitempty"`
}
```

---

## Yapılandırma Seçenekleri

| Fonksiyon | Açıklama |
| :--- | :--- |
| `WithQuickCheckoutBaseURL(url)` | Özel taban URL tanımlar (varsayılan: `https://www.shopier.com`). |
| `WithQuickCheckoutUserAgent(ua)` | İsteklerle gönderilecek özel `User-Agent` başlığı. |
| `WithQuickCheckoutHTTPClient(client)` | Özel proxy veya transport kuralları için `*http.Client`. |
| `WithQuickCheckoutTimeout(d)` | İstek zaman aşımı süresi (varsayılan: `30s`). |

---

## Hata Yönetimi ve Aşama Teşhisi

Hata durumlarında `*QuickCheckoutError` döndürülür ve hatanın hangi aşamada oluştuğu belirtilir:

```go
type QuickCheckoutStage string

const (
	StageParseInput        QuickCheckoutStage = "parse_input"
	StageFetchProduct      QuickCheckoutStage = "fetch_product"
	StageExtractCSRF       QuickCheckoutStage = "extract_csrf"
	StageCheckPayment      QuickCheckoutStage = "check_payment"
	StageFetchShippingForm QuickCheckoutStage = "fetch_shipping_form"
	StageProcessShipment   QuickCheckoutStage = "process_shipment"
)
```
