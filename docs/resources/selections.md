# Selections API Reference

Selections represent the individual values belonging to a product variation (such as "Red" under "Color" or "XL" under "Size").

## Methods

- `client.Selections.List(ctx, opts)` -> `(*PageResponse[Selection], error)`
- `client.Selections.All(ctx, opts)` -> `iter.Seq2[Selection, error]`
- `client.Selections.Get(ctx, id)` -> `(*Selection, error)`
- `client.Selections.Create(ctx, req)` -> `(*Selection, error)`
- `client.Selections.Update(ctx, id, req)` -> `(*Selection, error)`
- `client.Selections.Delete(ctx, id)` -> `error`

## Usage

```go
// List selections for a specific variation
opts := &shopier.SelectionListOptions{
	VariationIDs: []string{"var_123"},
}

res, err := client.Selections.List(ctx, opts)
if err != nil {
	log.Fatal(err)
}

for _, sel := range res.Items {
	fmt.Printf("[%s] %s (Variation ID: %s)\n", sel.ID, sel.Title, sel.VariationID)
}
```
