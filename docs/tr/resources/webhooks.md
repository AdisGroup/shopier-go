# Webhook Abonelikleri (Webhooks API)

Webhook API servisi, mağazanızdaki olay bildirim URL'lerini programatik olarak oluşturmanızı, listelemenizi ve silmenizi sağlar.

> [!NOTE]
> Gelen webhook isteklerini doğrulamak ve işlemek için [Webhook Entegrasyonu ve İmza Doğrulama Rehberi](/tr/webhooks) sayfasını inceleyin.

## Metotlar

- `client.Webhooks.List(ctx)` -> `([]WebhookSubscription, error)`
- `client.Webhooks.Create(ctx, req)` -> `(*WebhookSubscription, error)`
- `client.Webhooks.Delete(ctx, id)` -> `error`

---

## 1. Yeni Webhook Aboneliği Tanımlama

```go
sub, err := client.Webhooks.Create(ctx, &shopier.WebhookCreateRequest{
	Event: "order.created",
	URL:   "https://api.markam.com/webhooks/shopier",
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Abonelik Oluşturuldu! ID: %s\n", sub.ID)
// ÖNEMLİ: sub.Token değerini güvenli bir şekilde saklayın. Sadece ilk oluşturma yanıtında döner.
fmt.Printf("İmza Gizli Anahtarı (Secret): %s\n", sub.Token)
```

## 2. Aktif Abonelikleri Listeleme

```go
subs, err := client.Webhooks.List(ctx)
if err != nil {
	log.Fatal(err)
}

for _, s := range subs {
	fmt.Printf("[%s] Olay: %s -> %s\n", s.ID, s.Event, s.URL)
}
```

## 3. Webhook Aboneliğini Silme

```go
err := client.Webhooks.Delete(ctx, "wh_12345")
if err != nil {
	log.Fatal(err)
}
```

---

## Veri Modelleri ve Tipler

### `WebhookSubscription` (Webhook Abonelik Modeli) {#webhook-abonelik-modeli-webhooksubscription}

Kayıtlı bir olay bildirim URL'ini temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Benzersiz webhook abonelik kimliği |
| `Event` | `string` | `event` | Abone olunan olay adı (örn. `"order.created"`) |
| `URL` | `string` | `url` | Bildirimlerin iletildiği HTTPS uç noktası |
| `Token` | `string` | `token` | HMAC imza gizli anahtarı (*yalnızca oluşturulduğunda döner*) |

---

### `WebhookCreateRequest`

| Alan (Field) | Tip | JSON Etiketi | Zorunlu | Açıklama |
| :--- | :--- | :--- | :---: | :--- |
| `Event` | `string` | `event` | **Evet** | Abone olunacak olay adı |
| `URL` | `string` | `url` | **Evet** | HTTPS hedef bildirim adresi |

