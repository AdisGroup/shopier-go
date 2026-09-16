# Categories API Reference

The Categories API manages product grouping and navigation taxonomy.

## Methods

- `client.Categories.List(ctx, opts)` -> `(*PageResponse[Category], error)`
- `client.Categories.All(ctx, opts)` -> `iter.Seq2[Category, error]`
- `client.Categories.Get(ctx, id)` -> `(*Category, error)`
- `client.Categories.Create(ctx, req)` -> `(*Category, error)`
- `client.Categories.Update(ctx, id, req)` -> `(*Category, error)`
- `client.Categories.Delete(ctx, id)` -> `error`

---

## 1. List Categories

```go
for cat, err := range client.Categories.All(ctx, nil) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[%s] %s (Placement: %d)\n", cat.ID, cat.Title, cat.Placement)
}
```

## 2. Create Category

```go
category, err := client.Categories.Create(ctx, &shopier.CategoryCreateRequest{
	Title: "Winter Collection",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Created Category ID:", category.ID)
```

## 3. Update Category

```go
updated, err := client.Categories.Update(ctx, "cat_123", &shopier.CategoryUpdateRequest{
	Title: "Winter 2026 Collection",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("New Title:", updated.Title)
```

## 4. Delete Category

```go
err := client.Categories.Delete(ctx, "cat_123")
if err != nil {
	log.Fatal(err)
}
```
