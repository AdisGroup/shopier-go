# İndirimler (Discounts API)

İndirim servisi, hem sepet aşamasında girilen kupon kodlarını hem de sepet tutarına/adedine göre kendiliğinden devreye giren otomatik indirim kampanyalarını yönetir.

## Metotlar

### İndirim Kodları (Discount Codes)
- `client.Discounts.ListCodes(ctx, opts)` -> `(*PageResponse[DiscountCode], error)`
- `client.Discounts.AllCodes(ctx, opts)` -> `iter.Seq2[DiscountCode, error]`
- `client.Discounts.GetCode(ctx, id)` -> `(*DiscountCode, error)`
- `client.Discounts.CreateCode(ctx, req)` -> `(*DiscountCode, error)`
- `client.Discounts.UpdateCode(ctx, id, req)` -> `(*DiscountCode, error)`
- `client.Discounts.DeleteCode(ctx, id)` -> `error`

### Otomatik İndirimler (Automatic Discounts)
- `client.Discounts.ListAutomatic(ctx, opts)` -> `(*PageResponse[AutomaticDiscount], error)`
- `client.Discounts.AllAutomatic(ctx, opts)` -> `iter.Seq2[AutomaticDiscount, error]`
- `client.Discounts.GetAutomatic(ctx, id)` -> `(*AutomaticDiscount, error)`
- `client.Discounts.CreateAutomatic(ctx, req)` -> `(*AutomaticDiscount, error)`
- `client.Discounts.UpdateAutomatic(ctx, id, req)` -> `(*AutomaticDiscount, error)`
- `client.Discounts.DeleteAutomatic(ctx, id)` -> `error`

---

## 1. İndirim Kodu Oluşturma

```go
code, err := client.Discounts.CreateCode(ctx, &shopier.DiscountCodeCreateRequest{
	Code:          "HOSGELDIN20",
	Type:          "percent", // "percent" veya "amount"
	PercentOff:    "20",
	AmountMinimum: "100.00",
	Currency:      "TRY",
	NumAvailable:  500,
	ExpiresAt:     "2026-12-31+0300",
})
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Kupon Oluşturuldu: %s (ID: %s)\n", code.Code, code.ID)
```

## 2. Otomatik İndirim Kampanyası Tanımlama

```go
autoDisc, err := client.Discounts.CreateAutomatic(ctx, &shopier.AutomaticDiscountCreateRequest{
	Title:         "250 TL Üzeri %15 İndirim",
	Scope:         "all", // "all", "selectedProducts", "selectedCategories"
	Type:          "percent",
	PercentOff:    "15",
	Requirement:   "amount", // "amount" veya "quantity"
	AmountMinimum: "250.00",
	Currency:      "TRY",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Kampanya ID:", autoDisc.ID)
```

---

## Veri Modelleri ve Tipler

### `DiscountCode` (Kupon Kodu Modeli) {#kupon-kodu-modeli-discountcode}

Ödeme adımında müşterinin girdiği indirim kuponlarını temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Benzersiz indirim kodu numarası |
| `Code` | `string` | `code` | Kupon kodu metni (örn. `"HOSGELDIN20"`) |
| `Type` | `string` | `type` | `"percent"` (yüzdelik) veya `"amount"` (sabit tutar) |
| `AmountOff` | `string` | `amountOff` | Sabit indirim tutarı (`type == "amount"` ise) |
| `PercentOff` | `string` | `percentOff` | Yüzdelik indirim oranı (`type == "percent"` ise) |
| `AmountMinimum` | `string` | `amountMinimum` | Kuponun geçerli olması için gereken minimum sepet tutarı |
| `Currency` | `string` | `currency` | Para birimi (`"TRY"`, `"USD"` vb.) |
| `NumAvailable` | `int` | `numAvailable` | Toplam kullanım hakkı kotası |
| `NumUsed` | `int` | `numUsed` | Kullanılan adet sayısı |
| `ExpiresAt` | `string` | `expiresAt` | Kuponun son kullanma tarihi |
| `DateCreated` | `string` | `dateCreated` | Kuponun oluşturulma tarihi |

---

### `AutomaticDiscount` (Otomatik İndirim Modeli) {#otomatik-indirim-modeli-automaticdiscount}

Sepet koşulları sağlandığında otomatik uygulanan promosyon kampanyalarını temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Benzersiz kampanya numarası |
| `Title` | `string` | `title` | Kampanya başlığı |
| `Scope` | `string` | `scope` | Kapsam: `"all"` (tümü), `"selectedProducts"`, `"selectedCategories"` |
| `ProductIDs` | `[]string` | `productIds` | Kapsamdaki ürün ID listesi |
| `CategoryIDs` | `[]string` | `categoryIds` | Kapsamdaki kategori ID listesi |
| `Type` | `string` | `type` | `"percent"` (yüzdelik) veya `"amount"` (sabit tutar) |
| `AmountOff` | `string` | `amountOff` | Sabit indirim tutarı |
| `PercentOff` | `string` | `percentOff` | Yüzdelik indirim oranı |
| `Requirement` | `string` | `requirement` | Koşul tipi: `"amount"` (tutar) veya `"quantity"` (adet) |
| `AmountMinimum` | `string` | `amountMinimum` | Minimum sepet tutarı eşiği |
| `QuantityMinimum` | `int` | `quantityMinimum` | Minimum sepet ürün adedi eşiği |
| `StartsAt` | `string` | `startsAt` | Kampanya başlangıç tarihi |
| `ExpiresAt` | `string` | `expiresAt` | Kampanya bitiş tarihi |

