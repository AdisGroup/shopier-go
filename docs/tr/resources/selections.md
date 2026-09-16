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
