# Ürünler (Products API)

Ürün servisi, mağazanızdaki ürünlerin listelenmesini, yeni ürün eklenmesini, fiyat/stok/medya güncellemelerini ve ürün silme işlemlerini yönetir.

::: warning Ürün API Yetkilendirmesi (403 Forbidden)
Kişisel Erişim Jetonunuzda (PAT) veya API anahtarınızda `products:read` ve `products:write` yetkileri seçili olsa dahi, `/v1/products` uç noktaları (`List`, `Get`, `Create`, `Update`, `Delete`) varsayılan olarak `403 Forbidden` (*"Access to this resource on the server is denied"*) yanıtı dönebilir.

Shopier resmi dokümantasyonuna göre ürün katalog yönetimi ek mağaza onayı gerektirmektedir. Mağazanız için ürün uç noktalarını aktif ettirmek amacıyla mağaza bağlantınızla (`https://www.shopier.com/{magaza-adiniz}`) birlikte [hello@shopier.com](mailto:hello@shopier.com) adresine e-posta gönderebilirsiniz.
:::

## Metotlar

- `client.Products.List(ctx, opts)` -> `(*PageResponse[Product], error)`
- `client.Products.All(ctx, opts)` -> `iter.Seq2[Product, error]`
- `client.Products.Get(ctx, id)` -> `(*Product, error)`
- `client.Products.Create(ctx, req)` -> `(*Product, error)`
- `client.Products.Update(ctx, id, req)` -> `(*Product, error)`
- `client.Products.Delete(ctx, id)` -> `error`

---

## 1. Ürünleri Listeleme

```go
opts := &shopier.ProductListOptions{
	Status: "active", // "active", "draft", "outOfStock"
}

for prod, err := range client.Products.All(ctx, opts) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[%s] %s - Fiyat: %s %s (Stok: %d)\n",
		prod.ID, prod.Title, prod.PriceData.Price, prod.PriceData.Currency, prod.StockQuantity)
}
```

## 2. Yeni Ürün Oluşturma

```go
product, err := client.Products.Create(ctx, &shopier.ProductCreateRequest{
	Title:       "Premium Pamuk Hoodie",
	Description: "%100 pamuklu antrasit hoodie.",
	Type:        "physical", // "physical" veya "digital"
	ShippingPayer: "sellerPays", // "sellerPays" veya "buyerPays"
	PriceData: shopier.ProductPriceData{
		Currency: "TRY",
		Price:    "750.00",
		Discount: false,
	},
	StockQuantity: 50,
	Media: []shopier.ProductMedia{
		{
			Type:      "image",
			URL:       "https://mycdn.com/hoodie-front.jpg",
			Placement: 1,
		},
	},
})
if err != nil {
	log.Fatalf("Ürün oluşturulamadı: %v", err)
}

fmt.Printf("Oluşturulan Ürün ID: %s (URL: %s)\n", product.ID, product.URL)
```

## 3. Ürün Güncelleme

```go
newPrice := "799.00"
newStock := 45

updated, err := client.Products.Update(ctx, "prod_123", &shopier.ProductUpdateRequest{
	PriceData: &shopier.ProductPriceData{
		Currency: "TRY",
		Price:    newPrice,
	},
	StockQuantity: &newStock,
})
if err != nil {
	log.Fatal(err)
}

fmt.Println("Güncellenen Fiyat:", updated.PriceData.Price)
```

## 4. Ürün Silme

```go
err := client.Products.Delete(ctx, "prod_123")
if err != nil {
	log.Fatal(err)
}
fmt.Println("Ürün kalıcı olarak silindi.")
```

---

## Veri Modelleri ve Tipler

### `Product` (Ürün Modeli) {#urun-modeli-product}

Shopier mağazasında listelenen bir ürünü temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Benzersiz Shopier ürün numarası |
| `Title` | `string` | `title` | Ürün başlığı |
| `Description` | `string` | `description` | HTML/metin formatında ürün açıklaması |
| `Type` | `string` | `type` | `"physical"` (fiziksel ürün) veya `"digital"` (dijital indirme/kod) |
| `URL` | `string` | `url` | Ürünün doğrudan mağaza bağlantısı (`https://www.shopier.com/{magaza}/{id}`) |
| `Media` | `[]ProductMedia` | `media` | Ürüne ait görsel ve medya dosyaları |
| `PriceData` | [`ProductPriceData`](#productpricedata) | `priceData` | Fiyat, para birimi ve indirim bilgileri |
| `StockStatus` | `string` | `stockStatus` | Stok durumu (`"inStock"`, `"outOfStock"`) |
| `StockQuantity` | `int` | `stockQuantity` | Mevcut stok adedi |
| `ShippingPayer` | `string` | `shippingPayer` | `"sellerPays"` (ücretsiz kargo) veya `"buyerPays"` (alıcı öder) |
| `Categories` | `[]ProductCategoryRef` | `categories` | Bağlı olduğu kategori referansları |
| `Variants` | `[]ProductVariant` | `variants` | Ürüne tanımlı varyantlar, stok ve fiyatları |
| `Options` | `[]ProductOption` | `options` | İsteğe bağlı ekstra ürün seçenekleri |
| `DateCreated` | `string` | `dateCreated` | ISO-8601 ürün eklenme tarihi |

---

### `ProductPriceData`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `Currency` | `string` | `currency` | 3 haneli ISO para birimi kodu (`"TRY"`, `"USD"`, `"EUR"`) |
| `Price` | `string` | `price` | Standart baz satış fiyatı (örn. `"750.00"`) |
| `Discount` | `bool` | `discount` | İndirim aktifse `true` döner |
| `DiscountedPrice` | `string` | `discountedPrice` | İndirimli satış fiyatı |
| `ShippingPrice` | `string` | `shippingPrice` | Alıcı öderse sabit kargo ücreti |

---

### `ProductMedia`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Medya dosya kimliği |
| `Type` | `string` | `type` | Medya türü (`"image"`) |
| `URL` | `string` | `url` | Görselin doğrudan CDN bağlantısı |
| `Placement` | `int` | `placement` | Görsel sıralama indeksi (`1` = kapak görseli) |

---

### `ProductVariant`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `VariationID` | `string` | `variationId` | Bağlı varyasyon seti ID'si |
| `VariationTitle` | `string` | `variationTitle` | Varyasyon başlığı (örn. `"Beden"`, `"Renk"`) |
| `SelectionID` | `any` | `selectionId` | Seçilen seçenek ID'si / ID dizisi |
| `StockQuantity` | `int` | `stockQuantity` | Bu varyanta özel stok adedi |
| `PriceData` | `*ProductPriceData` | `priceData` | Bu varyanta özel fiyat ezme (override) nesnesi |

---

### `ProductCreateRequest`

| Alan (Field) | Tip | JSON Etiketi | Zorunlu | Açıklama |
| :--- | :--- | :--- | :---: | :--- |
| `Title` | `string` | `title` | **Evet** | Ürün başlığı |
| `Type` | `string` | `type` | **Evet** | `"physical"` veya `"digital"` |
| `PriceData` | `ProductPriceData` | `priceData` | **Evet** | Para birimi ve fiyat tanımı |
| `ShippingPayer` | `string` | `shippingPayer` | **Evet** | `"sellerPays"` veya `"buyerPays"` |
| `Media` | `[]ProductMedia` | `media` | **Evet** | En az bir ürün görseli |
| `Description` | `string` | `description` | Hayır | Detaylı ürün açıklaması |
| `StockQuantity` | `int` | `stockQuantity` | Hayır | Başlangıç stok adedi |
| `Categories` | `[]ProductCategoryRef` | `categories` | Hayır | Bağlı kategori referansları |
| `Variants` | `[]ProductVariant` | `variants` | Hayır | Varyant listesi |

