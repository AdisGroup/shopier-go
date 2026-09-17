# Variations & Selections API Reference

Variations define attribute axes (e.g. Size, Color), and Selections define individual options within those variations (e.g. Small, Medium, Large).

## Variations Methods

- `client.Variations.List(ctx, opts)` -> `(*PageResponse[Variation], error)`
- `client.Variations.All(ctx, opts)` -> `iter.Seq2[Variation, error]`
- `client.Variations.Get(ctx, id)` -> `(*Variation, error)`
- `client.Variations.Create(ctx, req)` -> `(*Variation, error)`
- `client.Variations.Update(ctx, id, req)` -> `(*Variation, error)`
- `client.Variations.Delete(ctx, id)` -> `error`

## Selections Methods

- `client.Selections.List(ctx, opts)` -> `(*PageResponse[Selection], error)`
- `client.Selections.All(ctx, opts)` -> `iter.Seq2[Selection, error]`
- `client.Selections.Get(ctx, id)` -> `(*Selection, error)`
- `client.Selections.Create(ctx, req)` -> `(*Selection, error)`
- `client.Selections.Update(ctx, id, req)` -> `(*Selection, error)`
- `client.Selections.Delete(ctx, id)` -> `error`

---

## Example: Creating a Variation and Options

```go
// 1. Create Variation (e.g., "Color")
variation, err := client.Variations.Create(ctx, &shopier.VariationCreateRequest{
	Title: "Color",
})
if err != nil {
	log.Fatal(err)
}

// 2. Add Selections (e.g., "Midnight Black", "Arctic White")
sel1, _ := client.Selections.Create(ctx, &shopier.SelectionCreateRequest{
	VariationID: variation.ID,
	Title:       "Midnight Black",
})

sel2, _ := client.Selections.Create(ctx, &shopier.SelectionCreateRequest{
	VariationID: variation.ID,
	Title:       "Arctic White",
})

fmt.Printf("Variation: %s | Options: %s, %s\n", variation.Title, sel1.Title, sel2.Title)
```

---

## Data Models & Types

### `Variation` {#variation-model}

Represents an attribute grouping dimension (e.g., "Size", "Material", "Color").

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Unique variation axis identifier |
| `Title` | `string` | `title` | Variation dimension name |
| `Placement` | `int` | `placement` | Display ordering sequence number |

---

### `VariationCreateRequest` & `VariationUpdateRequest`

| Field | Type | JSON Tag | Required | Description |
| :--- | :--- | :--- | :---: | :--- |
| `Title` | `string` | `title` | **Yes** | Variation title / attribute name |

