# Hata Yönetimi ve Yeniden Deneme (Retry)

SDK, Shopier HTTP hata yanıtlarını tip güvenli Go hatalarına dönüştürür ve geçici ağ/sunucu aksaklıkları ile rate limit durumlarında otomatik yeniden deneme (retry) desteği sunar.

---

## Hata Tipleri ve Yakalama

Tüm API hata yanıtları standart Go `error` arayüzünü uygular ve `errors.As` ile tiplerine ayrıştırılabilir:

```go
import (
	"errors"
	"github.com/AdisGroup/shopier-go"
)

res, err := client.Orders.Get(ctx, "gecersiz_id")
if err != nil {
	var rateLimitErr *shopier.RateLimitError
	var authErr *shopier.AuthenticationError
	var apiErr *shopier.APIError

	switch {
	case errors.As(err, &rateLimitErr):
		// 429 Too Many Requests
		fmt.Printf("İstek limiti aşıldı. Beklenmesi gereken süre: %s\n", rateLimitErr.RetryAfter)

	case errors.As(err, &authErr):
		// 401 Unauthorized
		fmt.Println("Geçersiz veya süresi dolmuş token.")

	case errors.As(err, &apiErr):
		// 400 Bad Request, 404 Not Found, 403 Forbidden, 5xx vb.
		fmt.Printf("API Hatası [%d]: %s (Açıklama: %s)\n",
			apiErr.StatusCode,
			apiErr.ErrorCode,
			apiErr.Message,
		)
		// Ham sunucu yanıtını inceleyin
		fmt.Println("Ham Yanıt:", string(apiErr.RawBody))

	default:
		// Context zaman aşımı veya ağ kopması
		fmt.Printf("Ağ/Taşıyıcı Hatası: %v\n", err)
	}
}
```

---

## Sık Karşılaşılan HTTP Hata Kodları

| Durum Kodu | Hata Tipi | Neden ve Çözüm |
| :--- | :--- | :--- |
| `400 Bad Request` | `invalid_request` / `validation_failed` | Eksik veya hatalı istek parametresi. İstek nesnesi alanlarını kontrol edin. |
| `401 Unauthorized` | `unauthorized` | API jetonu / PAT eksik, hatalı veya süresi dolmuş. |
| `403 Forbidden` | `forbidden` | Jeton yetkileri yetersiz veya uç nokta (örn. `/v1/products`) mağaza onayına tabi. `hello@shopier.com` ile iletişime geçin. |
| `404 Not Found` | `not_found` | İstenen kaynak (sipariş ID, ürün ID vb.) bulunamadı. |
| `429 Too Many Requests` | `rate_limit_exceeded` | İstek kotası aşıldı (200 istek/dk). SDK tarafından otomatik olarak bekletilip yeniden denenir. |
| `500 / 503 / 504` | `server_error` | Shopier sunucu kesintisi veya aşırı yüklenme. Exponential backoff ile otomatik yinelenir. |

---

## Rate Limit (İstek Limiti) Politikası

Shopier API çağrılarını **60 saniyelik kayan pencere (sliding window)** içinde sınırlar:
- Standart kota: Uygulama/kullanıcı çifti başına **dakikada 200 istek**.
- Kota aşıldığında Shopier `HTTP 429 Too Many Requests` döner ve yanıta `Retry-After: <saniye>` başlığını ekler.

### Otomatik Yeniden Deneme

SDK, `429` yanıtını yakaladığında `Retry-After` başlığındaki süre kadar goroutine'i bekletir ve isteği `WithMaxRetries` sınırına kadar otomatik olarak tekrarlar.

```go
client, err := shopier.NewClient("token",
	shopier.WithMaxRetries(3), // 429 ve 5xx durumlarında 3 defaya kadar yeniden dene
)
```

---

## 5xx Sunucu Hatalarında Üstel Geri Çekilme (Backoff)

Shopier geçici bir sunucu hatası (`500`, `502`, `503`, `504`) döndürdüğünde, istemci rastgele dağıtılmış jitter içeren üstel geri çekilme uygular:

$$\text{gecikme} = \min(\text{ilkGecikme} \times 2^{\text{deneme}}, \text{maksGecikme}) + \text{jitter}$$
