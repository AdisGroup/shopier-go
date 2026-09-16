# Kategoriler (Categories API)

Kategori servisi, mağazanızdaki ürünleri gruplandırmanızı ve vitrin menülerini yapılandırmanızı sağlar.

## Metotlar

- `client.Categories.List(ctx, opts)` -> `(*PageResponse[Category], error)`
- `client.Categories.All(ctx, opts)` -> `iter.Seq2[Category, error]`
- `client.Categories.Get(ctx, id)` -> `(*Category, error)`
- `client.Categories.Create(ctx, req)` -> `(*Category, error)`
- `client.Categories.Update(ctx, id, req)` -> `(*Category, error)`
- `client.Categories.Delete(ctx, id)` -> `error`

---

## 1. Kategorileri Listeleme

```go
for cat, err := range client.Categories.All(ctx, nil) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[%s] %s (Sıralama: %d)\n", cat.ID, cat.Title, cat.Placement)
}
```

## 2. Kategori Oluşturma

```go
category, err := client.Categories.Create(ctx, &shopier.CategoryCreateRequest{
	Title: "Kış Koleksiyonu",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Kategori ID:", category.ID)
```

## 3. Kategori Güncelleme

```go
updated, err := client.Categories.Update(ctx, "cat_123", &shopier.CategoryUpdateRequest{
	Title: "2026 Kış Koleksiyonu",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Yeni Başlık:", updated.Title)
```

## 4. Kategori Silme

```go
err := client.Categories.Delete(ctx, "cat_123")
if err != nil {
	log.Fatal(err)
}
```
