# Products API Reference

The Products API manages merchandise listings, media images, pricing tiers, stock status, and variations.

## Methods

- `client.Products.List(ctx, opts)` -> `(*PageResponse[Product], error)`
- `client.Products.All(ctx, opts)` -> `iter.Seq2[Product, error]`
- `client.Products.Get(ctx, id)` -> `(*Product, error)`
- `client.Products.Create(ctx, req)` -> `(*Product, error)`
- `client.Products.Update(ctx, id, req)` -> `(*Product, error)`
- `client.Products.Delete(ctx, id)` -> `error`

---

## 1. List Products

```go
opts := &shopier.ProductListOptions{
	Status: "active", // "active", "draft", "outOfStock"
}

for prod, err := range client.Products.All(ctx, opts) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[%s] %s - Price: %s %s (Stock: %d)\n",
		prod.ID, prod.Title, prod.PriceData.Price, prod.PriceData.Currency, prod.StockQuantity)
}
```

## 2. Create Product

```go
product, err := client.Products.Create(ctx, &shopier.ProductCreateRequest{
	Title:       "Premium Cotton Hoodie",
	Description: "Heavyweight 100% cotton hoodie in charcoal black.",
	Type:        "physical", // "physical" or "digital"
	ShippingPayer: "sellerPays", // "sellerPays" or "buyerPays"
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
	log.Fatalf("Failed to create product: %v", err)
}

fmt.Printf("Created Product ID: %s (URL: %s)\n", product.ID, product.URL)
```

## 3. Update Product

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

fmt.Println("Updated Price:", updated.PriceData.Price)
```

## 4. Delete Product

```go
err := client.Products.Delete(ctx, "prod_123")
if err != nil {
	log.Fatal(err)
}
fmt.Println("Product permanently deleted.")
```
