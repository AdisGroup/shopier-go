# Varyasyonlar ve Seçenekler (Variations & Selections)

Varyasyonlar ürün özellik eksenlerini (Renk, Beden vb.), Seçenekler ise bu varyasyonların altındaki somut tercihleri (Kırmızı, XL vb.) tanımlar.

## Varyasyon Metotları

- `client.Variations.List(ctx, opts)` -> `(*PageResponse[Variation], error)`
- `client.Variations.All(ctx, opts)` -> `iter.Seq2[Variation, error]`
- `client.Variations.Get(ctx, id)` -> `(*Variation, error)`
- `client.Variations.Create(ctx, req)` -> `(*Variation, error)`
- `client.Variations.Update(ctx, id, req)` -> `(*Variation, error)`
- `client.Variations.Delete(ctx, id)` -> `error`

## Seçenek Metotları

- `client.Selections.List(ctx, opts)` -> `(*PageResponse[Selection], error)`
- `client.Selections.All(ctx, opts)` -> `iter.Seq2[Selection, error]`
- `client.Selections.Get(ctx, id)` -> `(*Selection, error)`
- `client.Selections.Create(ctx, req)` -> `(*Selection, error)`
- `client.Selections.Update(ctx, id, req)` -> `(*Selection, error)`
- `client.Selections.Delete(ctx, id)` -> `error`

---

## Örnek: Varyasyon ve Seçenek Tanımlama

```go
// 1. Varyasyon Oluştur (Örn: "Beden")
variation, err := client.Variations.Create(ctx, &shopier.VariationCreateRequest{
	Title: "Beden",
})
if err != nil {
	log.Fatal(err)
}

// 2. Seçenekleri Ekle (Örn: "M", "L")
sel1, _ := client.Selections.Create(ctx, &shopier.SelectionCreateRequest{
	VariationID: variation.ID,
	Title:       "M",
})

sel2, _ := client.Selections.Create(ctx, &shopier.SelectionCreateRequest{
	VariationID: variation.ID,
	Title:       "L",
})

fmt.Printf("Varyasyon: %s | Seçenekler: %s, %s\n", variation.Title, sel1.Title, sel2.Title)
```
