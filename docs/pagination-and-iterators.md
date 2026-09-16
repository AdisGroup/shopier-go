# Pagination & Iterators

Shopier returns paginated lists on all `List` endpoints. Metadata about pagination is delivered via HTTP response headers:
- `Shopier-Pagination-Page`
- `Shopier-Pagination-Limit`
- `Shopier-Pagination-Total-Pages`
- `Shopier-Pagination-Total-Items`

The SDK provides two distinct ways to work with paginated data:
1. **Manual Page Requests (`List`)**: Full control over page numbers, limits, and header metadata.
2. **Go 1.23+ Range-Over-Func Iterators (`All`)**: Automated traversal across all pages using native `for ... range` syntax.

---

## 1. Modern Go 1.23+ Range Iterators (`All`)

The `All` method returns an `iter.Seq2[T, error]`. It transparently requests subsequent pages as the loop iterates over items, eliminating manual loop boilerplate:

```go
ctx := context.Background()

// Stream every order across all pages with zero pagination boilerplate
for order, err := range client.Orders.All(ctx, nil) {
	if err != nil {
		log.Fatalf("Iteration error: %v", err)
	}

	fmt.Printf("Order #%s: %s %s\n", order.ID, order.Totals.Total, order.Currency)
}
```

You can pass filtering options to `All` as well:

```go
opts := &shopier.OrderListOptions{
	FulfillmentStatus: "unfulfilled",
}

for order, err := range client.Orders.All(ctx, opts) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pending fulfillment:", order.ID)
}
```

---

## 2. Manual Page Requests (`List`)

When building user interfaces or explicit batch workers that require exact page navigation:

```go
opts := &shopier.OrderListOptions{
	ListOptions: shopier.ListOptions{
		Page:  1,
		Limit: 20, // Min: 1, Max: 50, Default: 10
		Sort:  "desc",
	},
	FulfillmentStatus: "unfulfilled",
}

res, err := client.Orders.List(ctx, opts)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Current Page: %d / %d (Total Items: %d)\n",
	res.Pagination.Page,
	res.Pagination.TotalPages,
	res.Pagination.TotalItems,
)

for _, order := range res.Items {
	fmt.Printf("- Order: %s\n", order.ID)
}

// Check if more pages exist
if res.Pagination.HasNextPage() {
	opts.Page = res.Pagination.NextPage()
	// Fetch next page...
}
```

### `PaginationInfo` Struct

```go
type PaginationInfo struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalPages int  `json:"total_pages"`
	TotalItems int  `json:"total_items"`
}

func (p PaginationInfo) HasNextPage() bool
func (p PaginationInfo) NextPage() int
```
