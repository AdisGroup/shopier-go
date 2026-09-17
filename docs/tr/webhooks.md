# Webhook Entegrasyonu ve İmza Doğrulama

Shopier, mağazanızda gerçekleşen olayları (sipariş, iade, ürün değişiklikleri) belirlediğiniz Bildirim URL'ine (Notification URL) HTTP POST isteği olarak iletir. İsteklerin Shopier'dan geldiğini doğrulamak için tüm bildirimler uygulamanızın webhook token'ı ile **HS256 (HMAC-SHA256)** algoritmasıyla imzalanır.

---

## Webhook Başlıkları (Headers)

| Başlık | Açıklama |
| :--- | :--- |
| `Shopier-Signature` | HMAC-SHA256 imza değeri (hex veya base64) |
| `Shopier-Event` | Olay tipi adı (örn. `order.created`) |
| `Shopier-Webhook-Id` | Bildirimin benzersiz kimlik numarası |
| `Shopier-Timestamp` | Saniye cinsinden UTC Unix zaman damgası |
| `Shopier-Account-Id` | İlgili mağazanın hesap kimliği |
| `Shopier-Api-Version` | API sürüm numarası |

---

## Desteklenen Olaylar

| Olay Sabiti | Olay Adı | Veri Modeli | Tetiklenme Anı |
| :--- | :--- | :--- | :--- |
| `webhook.EventOrderCreated` | `order.created` | [`*shopier.Order`](/tr/resources/orders#siparis-modeli-order) | Yeni bir sipariş ödendiğinde |
| `webhook.EventOrderAddressUpdated` | `order.addressUpdated` | [`*shopier.Order`](/tr/resources/orders#siparis-modeli-order) | Alıcı teslimat adresi güncellendiğinde |
| `webhook.EventOrderFulfilled` | `order.fulfilled` | [`*shopier.Order`](/tr/resources/orders#siparis-modeli-order) | Sipariş kargolandığında/tamamlandığında |
| `webhook.EventProductCreated` | `product.created` | [`*shopier.Product`](/tr/resources/products#urun-modeli-product) | Yeni bir ürün yayınlandığında |
| `webhook.EventProductUpdated` | `product.updated` | [`*shopier.Product`](/tr/resources/products#urun-modeli-product) | Mevcut ürün güncellendiğinde |
| `webhook.EventRefundRequested` | `refund.requested` | [`*shopier.Refund`](/tr/resources/refunds#iade-modeli-refund) | İade talebi oluşturulduğunda |
| `webhook.EventRefundUpdated` | `refund.updated` | [`*shopier.Refund`](/tr/resources/refunds#iade-modeli-refund) | İade onaylandığında veya reddedildiğinde |

---

## Webhook Veri Modelleri

### `webhook.Event`

```go
type Event struct {
	Header     Header // HTTP başlık üstverileri
	RawPayload []byte // Ham JSON gövdesi
}
```

#### Yardımcı Ayrıştırıcı Metotlar:
- `evt.Order()` -> `(*shopier.Order, error)`: JSON gövdesini [`shopier.Order`](/tr/resources/orders#siparis-modeli-order) modeline dönüştürür.
- `evt.Product()` -> `(*shopier.Product, error)`: JSON gövdesini [`shopier.Product`](/tr/resources/products#urun-modeli-product) modeline dönüştürür.
- `evt.Refund()` -> `(*shopier.Refund, error)`: JSON gövdesini [`shopier.Refund`](/tr/resources/refunds#iade-modeli-refund) modeline dönüştürür.

### `webhook.Header`

| Alan (Field) | Tip | HTTP Başlığı | Açıklama |
| :--- | :--- | :--- | :--- |
| `WebhookID` | `string` | `Shopier-Webhook-Id` | Bildirimin benzersiz kimlik numarası |
| `Event` | `string` | `Shopier-Event` | Olay tipi adı (örn. `order.created`) |
| `Timestamp` | `int64` | `Shopier-Timestamp` | Saniye cinsinden UTC Unix zaman damgası |
| `Signature` | `string` | `Shopier-Signature` | HMAC-SHA256 imza değeri |
| `AccountID` | `string` | `Shopier-Account-Id` | İlgili mağazanın hesap kimliği |
| `APIVersion` | `string` | `Shopier-Api-Version` | API sürüm numarası |

---

## 1. Hazır `http.Handler` ile Webhook Dinleyici

En hızlı ve standart `net/http` uyumlu webhook dinleyici kurulumu:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AdisGroup/shopier-go/webhook"
)

func main() {
	secret := os.Getenv("SHOPIER_WEBHOOK_SECRET")

	handler := webhook.NewHandler(secret, func(ctx context.Context, evt *webhook.Event) error {
		log.Printf("Olay alındı: %s (ID: %s)", evt.Header.Event, evt.Header.WebhookID)

		switch evt.Header.Event {
		case webhook.EventOrderCreated:
			order, err := evt.Order()
			if err != nil {
				return err
			}
			fmt.Printf("Yeni Sipariş: #%s (%s %s)\n", order.ID, order.Totals.Total, order.Currency)

		case webhook.EventRefundRequested:
			refund, err := evt.Refund()
			if err != nil {
				return err
			}
			fmt.Printf("İade Talebi: #%s (Sipariş: %s)\n", refund.ID, refund.OrderID)
		}

		return nil
	})

	http.Handle("/webhooks/shopier", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 2. Popüler Web Framework'leri ile Kullanım (Gin, Chi, Fiber)

### `webhook.Parse` Kullanımı (Gin, Chi)

```go
event, err := webhook.Parse(c.Request, webhookSecret)
if err != nil {
	c.String(http.StatusUnauthorized, "Yetkisiz istek: %v", err)
	return
}

order, _ := event.Order()
```

### Doğrudan İmza Doğrulama (`webhook.VerifySignature` - Fiber)

```go
app.Post("/webhooks/shopier", func(c *fiber.Ctx) error {
	sig := c.Get("Shopier-Signature")
	body := c.Body()

	if !webhook.VerifySignature(body, sig, webhookSecret) {
		return c.Status(fiber.StatusUnauthorized).SendString("Geçersiz imza")
	}

	// Payload'ı işleyin...
	return c.SendStatus(fiber.StatusOK)
})
```
