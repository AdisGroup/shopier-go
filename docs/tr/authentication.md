# Kimlik Doğrulama (Authentication)

Shopier iki farklı kimlik doğrulama modelini destekler:
1. **Personal Access Token (PAT)**: Geliştirici panelinden üretilen statik belirteçler. Arka plan servisleri, cron görevleri ve tekil mağaza entegrasyonları için uygundur.
2. **OAuth 2.0**: Mağaza sahiplerinin uygulamanıza yetki verdiği çok kullanıcılı (multi-tenant) uygulamalar için standart yetkilendirme kodu akışı.

---

## 1. Personal Access Token (PAT) ile Kullanım

Geliştirici portalından aldığınız belirteci doğrudan `shopier.NewClient` fonksiyonuna iletin:

```go
client, err := shopier.NewClient("panelden_alinan_pat_token")
if err != nil {
	log.Fatal(err)
}
```

---

## 2. OAuth 2.0 Entegrasyonu

Shopier OAuth 2.0 uç noktaları özel port kurallarına sahiptir:
- **Yetkilendirme Sayfası**: `https://developer.shopier.com/v1/oauth2/authorize`
- **Token Takası & Yenileme**: `https://api.shopier.com:8443/v1/oauth2/token`
- **Token İptali (Revoke)**: `https://api.shopier.com:8443/v1/oauth2/revoke`

`oauth` paketi bu port ve adres yapılandırmasını varsayılan olarak doğru yönetir.

### OAuth Yapılandırması

```go
import "github.com/AdisGroup/shopier-go/oauth"

cfg := oauth.NewConfig(
	"CLIENT_ID",
	"CLIENT_SECRET",
	"https://uygulamaniz.com/oauth/callback",
)
```

### 1. Adım: Yetkilendirme URL'ine Yönlendirme

İstenen izin kapsamları (scopes) ve CSRF koruması için state parametresiyle onay URL'ini oluşturun:

```go
consentURL := cfg.AuthCodeURL(
	"csrf_state_token",
	oauth.ScopeOrdersRead,
	oauth.ScopeOrdersWrite,
	oauth.ScopeProductsRead,
	oauth.ScopeProductsWrite,
)

// Kullanıcıyı consentURL adresine yönlendirin
```

### 2. Adım: Yetki Kodunu (Code) Erişim Belirtecine Dönüştürme

Callback endpoint'inize gelen geçici `code` parametresini erişim ve yenileme belirteçlerine dönüştürün:

```go
token, err := cfg.Exchange(ctx, code)
if err != nil {
	log.Fatalf("OAuth token takası başarısız: %v", err)
}

fmt.Printf("Erişim Belirteci: %s (Geçerlilik süresi: %d saniye)\n", token.AccessToken, token.ExpiresIn)
fmt.Printf("Yenileme Belirteci: %s\n", token.RefreshToken)
```

### 3. Adım: Süresi Dolan Belirteçleri Yenileme

Shopier erişim belirteçleri 3 gün (259.200 saniye) geçerlidir. Süre dolduğunda yenileme belirteci ile yeni bir erişim belirteci alın:

```go
if token.Expired() {
	newToken, err := cfg.RefreshToken(ctx, token.RefreshToken)
	if err != nil {
		log.Fatalf("Token yenileme başarısız: %v", err)
	}
	token = newToken
}
```

### 4. Adım: Belirteç İptali (Revoke)

Mağaza sahibi uygulamanızın bağlantısını kestiğinde belirteci geçersiz kılın:

```go
err := cfg.Revoke(ctx, token.AccessToken)
if err != nil {
	log.Fatalf("Token iptali başarısız: %v", err)
}
```

---

## İzin Kapsamları (Scopes) Referansı

| Sabit Adı | Kapsam Değeri | İzin Verilen İşlemler |
| :--- | :--- | :--- |
| `oauth.ScopeOrdersRead` | `orders:read` | Sipariş listeleme/detay, işlem hareketleri, sipariş webhook'ları |
| `oauth.ScopeOrdersWrite` | `orders:write` | Sipariş kargo takip ve teslimat güncelleme |
| `oauth.ScopeProductsRead` | `products:read` | Ürünler, kategoriler, varyasyonlar, seçenekleri okuma |
| `oauth.ScopeProductsWrite` | `products:write` | Ürün ve taksonomi oluşturma, güncelleme, silme |
| `oauth.ScopeShippingsRead` | `shippings:read` | Kargo etiket ve takip durumlarını okuma |
| `oauth.ScopeShippingsWrite` | `shippings:write` | Anlaşmalı kargo kodu üretme ve iptal |
| `oauth.ScopeDiscountsRead` | `discounts:read` | İndirim kodları ve otomatik indirimleri okuma |
| `oauth.ScopeDiscountsWrite` | `discounts:write` | İndirim kampanyaları oluşturma ve yönetme |
| `oauth.ScopePayoutsRead` | `payouts:read` | Hakediş ödemeleri ve bakiye hareketlerini okuma |
| `oauth.ScopeRefundsRead` | `refunds:read` | İade kayıtlarını okuma |
| `oauth.ScopeRefundsWrite` | `refunds:write` | Sipariş iade talebi başlatma |
| `oauth.ScopeShopRead` | `shop:read` | Mağaza sahibi ve mağaza ayarlarını okuma |
| `oauth.ScopeShopWrite` | `shop:write` | Mağaza ayarlarını güncelleme |

---

## Veri Modelleri ve Tipler

### `oauth.Token` {#oauthtoken-model}

Shopier OAuth 2.0 sunucusu tarafından üretilen kimlik doğrulama belirteç nesnesidir.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `AccessToken` | `string` | `access_token` | API çağrılarında `Authorization: Bearer <token>` olarak iletilen erişim belirteci |
| `RefreshToken` | `string` | `refresh_token` | Süresi dolan erişim belirtecini yenilemek için kullanılan uzun ömürlü belirteç |
| `TokenType` | `string` | `token_type` | Belirteç türü (`"Bearer"`) |
| `ExpiresIn` | `int64` | `expires_in` | Saniye cinsinden geçerlilik süresi (genelde `259200` = 3 gün) |
| `Scope` | `string` | `scope` | Boşlukla ayrılmış yetkilendirilmiş izin kapsamları listesi |

#### Yardımcı Metotlar:
- `token.Expired() bool`: Belirtecin geçerlilik süresinin dolup dolmadığını kontrol eder.

