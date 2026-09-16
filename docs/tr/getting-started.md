# Genel Bakış ve Kurulum

`shopier-go`, Shopier REST API'si için geliştirilmiş üretime hazır (production-ready) bir Go SDK ve API wrapper kütüphanesidir. Sıfır harici bağımlılık, otomatik yeniden deneme (retry) mekanizması, context yönetimi ve tam tip güvenliği sunar.

## Gereksinimler

- Go 1.23 veya üzeri

## Kurulum

Projeye standart `go get` komutuyla ekleyin:

```bash
go get github.com/AdisGroup/shopier-go
```

## Hızlı Başlangıç

İstemciyi Personal Access Token (PAT) veya OAuth2 erişim belirteci ile başlatabilirsiniz:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AdisGroup/shopier-go"
)

func main() {
	client, err := shopier.NewClient(os.Getenv("SHOPIER_TOKEN"))
	if err != nil {
		log.Fatalf("İstemci başlatılamadı: %v", err)
	}

	ctx := context.Background()

	// Güncel hesap bakiyelerini sorgula
	balances, err := client.Balance.Get(ctx)
	if err != nil {
		log.Fatalf("Bakiye sorgusu başarısız: %v", err)
	}

	for _, b := range balances {
		fmt.Printf("Bakiye: %s %s\n", b.Amount, b.Currency)
	}
}
```

## İstemci Yapılandırması (Functional Options)

İstemci oluşturulurken ağ, zaman aşımı, loglama ve yeniden deneme politikaları özelleştirilebilir:

```go
import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AdisGroup/shopier-go"
)

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

client, err := shopier.NewClient(os.Getenv("SHOPIER_TOKEN"),
	shopier.WithTimeout(15*time.Second),
	shopier.WithMaxRetries(3),
	shopier.WithRetryDelay(500*time.Millisecond, 10*time.Second),
	shopier.WithLogger(logger),
	shopier.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	shopier.WithUserAgent("ozel-servis/1.0"),
)
```

### Yapılandırma Seçenekleri

| Seçenek | Varsayılan | Açıklama |
| :--- | :--- | :--- |
| `WithBaseURL(url)` | `https://api.shopier.com/v1` | Özel REST API kök adresi |
| `WithOAuthBaseURL(url)` | `https://api.shopier.com:8443/v1` | Özel OAuth2 kök adresi (port 8443) |
| `WithTimeout(d)` | `30s` | Context deadline içermeyen istekler için varsayılan zaman aşımı |
| `WithMaxRetries(n)` | `3` | 429 ve 5xx yanıtlarında maksimum yeniden deneme adedi |
| `WithRetryDelay(initial, max)` | `500ms`, `10s` | Üstel geri çekilme (backoff) süre parametreleri |
| `WithLogger(logger)` | Sessiz (Discard) | Yapılandırılmış `*slog.Logger` nesnesi |
| `WithHTTPClient(client)` | Standart istemci | Özel `*http.Client` taşıyıcısı |
| `WithUserAgent(ua)` | `shopier-go/1.0.0` | İsteklerde iletilecek User-Agent değeri |
