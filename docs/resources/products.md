# Products API Reference

The Products API manages merchandise listings, media images, pricing tiers, stock status, and variations.

::: warning Product Catalog API Access (403 Forbidden)
By default, access to `/v1/products` endpoints (`List`, `Get`, `Create`, `Update`, `Delete`) may return `403 Forbidden` (*"Access to this resource on the server is denied"*) even if your Personal Access Token includes `products:read` and `products:write` scopes.

As documented in Shopier's API reference, product catalog endpoints require additional store-level activation. To enable product endpoints for your store, contact Shopier developer support at [hello@shopier.com](mailto:hello@shopier.com) with your store URL (`https://www.shopier.com/{your-store}`).
:::

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

---

## Data Models & Types

### `Product` {#product-model}

Represents a merchandise item listed on Shopier.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Unique Shopier product ID |
| `Title` | `string` | `title` | Product title |
| `Description` | `string` | `description` | HTML/text product description |
| `Type` | `string` | `type` | `"physical"` (tangible goods) or `"digital"` (instant downloads/keys) |
| `URL` | `string` | `url` | Direct storefront product URL (e.g. `https://www.shopier.com/{shop}/{id}`) |
| `Media` | `[]ProductMedia` | `media` | Uploaded product gallery images |
| `PriceData` | [`ProductPriceData`](#productpricedata) | `priceData` | Pricing, currency, and discount structure |
| `StockStatus` | `string` | `stockStatus` | `"inStock"` or `"outOfStock"` |
| `StockQuantity` | `int` | `stockQuantity` | Available inventory count |
| `ShippingPayer` | `string` | `shippingPayer` | `"sellerPays"` (free shipping) or `"buyerPays"` |
| `Categories` | `[]ProductCategoryRef` | `categories` | Associated category IDs and titles |
| `Variants` | `[]ProductVariant` | `variants` | Variant stock, prices, and variation attributes |
| `Options` | `[]ProductOption` | `options` | Optional add-on choices |
| `DateCreated` | `string` | `dateCreated` | ISO-8601 creation timestamp |

---

### `ProductPriceData`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `Currency` | `string` | `currency` | 3-letter ISO currency code (`"TRY"`, `"USD"`, `"EUR"`) |
| `Price` | `string` | `price` | Regular base price (e.g. `"750.00"`) |
| `Discount` | `bool` | `discount` | `true` if promotional discount is active |
| `DiscountedPrice` | `string` | `discountedPrice` | Discounted sale price |
| `ShippingPrice` | `string` | `shippingPrice` | Fixed shipping fee if buyer pays |

---

### `ProductMedia`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Media asset identifier |
| `Type` | `string` | `type` | Asset type (`"image"`) |
| `URL` | `string` | `url` | Full image CDN URL |
| `Placement` | `int` | `placement` | Display order index (`1` = cover image) |

---

### `ProductVariant`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `VariationID` | `string` | `variationId` | Associated variation set ID |
| `VariationTitle` | `string` | `variationTitle` | Variation title (e.g. `"Size"`, `"Color"`) |
| `SelectionID` | `any` | `selectionId` | Selected option ID(s) |
| `StockQuantity` | `int` | `stockQuantity` | Inventory quantity for this specific variant |
| `PriceData` | `*ProductPriceData` | `priceData` | Optional variant-specific override price |

---

### `ProductCreateRequest`

| Field | Type | JSON Tag | Required | Description |
| :--- | :--- | :--- | :---: | :--- |
| `Title` | `string` | `title` | **Yes** | Product listing title |
| `Type` | `string` | `type` | **Yes** | `"physical"` or `"digital"` |
| `PriceData` | `ProductPriceData` | `priceData` | **Yes** | Currency and price specification |
| `ShippingPayer` | `string` | `shippingPayer` | **Yes** | `"sellerPays"` or `"buyerPays"` |
| `Media` | `[]ProductMedia` | `media` | **Yes** | At least one image asset |
| `Description` | `string` | `description` | No | Detailed listing description |
| `StockQuantity` | `int` | `stockQuantity` | No | Initial inventory level |
| `Categories` | `[]ProductCategoryRef` | `categories` | No | Associated category references |
| `Variants` | `[]ProductVariant` | `variants` | No | Variant inventory and pricing |

