# Seçenekler (Selections API)

Seçenekler, bir varyasyona (Örn. "Renk") ait değerleri (Örn. "Kırmızı", "Mavi") temsil eder.

## Metotlar

- `client.Selections.List(ctx, opts)` -> `(*PageResponse[Selection], error)`
- `client.Selections.All(ctx, opts)` -> `iter.Seq2[Selection, error]`
- `client.Selections.Get(ctx, id)` -> `(*Selection, error)`
- `client.Selections.Create(ctx, req)` -> `(*Selection, error)`
- `client.Selections.Update(ctx, id, req)` -> `(*Selection, error)`
- `client.Selections.Delete(ctx, id)` -> `error`

## Kullanım

```go
opts := &shopier.SelectionListOptions{
	VariationIDs: []string{"var_123"},
}

res, err := client.Selections.List(ctx, opts)
if err != nil {
	log.Fatal(err)
}

for _, sel := range res.Items {
	fmt.Printf("[%s] %s (Varyasyon ID: %s)\n", sel.ID, sel.Title, sel.VariationID)
}
```

---

## Veri Modelleri ve Tipler

### `Selection` (Seçenek Modeli) {#secenek-modeli-selection}

Bir varyasyonun altındaki somut seçenek değerini (Örn. "XL", "Lacivert") temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Benzersiz seçenek numarası |
| `Title` | `string` | `title` | Seçenek adı / değeri |
| `VariationID` | `string` | `variationId` | Bağlı olduğu varyasyon grubu numarası |

---

### `SelectionCreateRequest`

| Alan (Field) | Tip | JSON Etiketi | Zorunlu | Açıklama |
| :--- | :--- | :--- | :---: | :--- |
| `VariationID` | `string` | `variationId` | **Evet** | Hedef varyasyon ID |
| `Title` | `string` | `title` | **Evet** | Seçenek adı |

---

### `SelectionListOptions`

| Alan (Field) | Tip | URL Parametresi | Açıklama |
| :--- | :--- | :--- | :--- |
| `VariationIDs` | `[]string` | `variationId` | Belirli varyasyonlara ait seçenekleri filtreler |
| `Page` | `int` | `page` | Sayfa indeksi (1 tabanlı) |
| `Limit` | `int` | `limit` | Sayfa başına kayıt sayısı (varsayılan: 50) |

