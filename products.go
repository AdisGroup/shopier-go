package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// Product represents a merchandise listing in a Shopier storefront.
type Product struct {
	ID               string               `json:"id"`
	Title            string               `json:"title"`
	Description      string               `json:"description,omitempty"`
	Type             string               `json:"type"`
	DateCreated      string               `json:"dateCreated,omitempty"`
	URL              string               `json:"url,omitempty"`
	Media            []ProductMedia       `json:"media"`
	PriceData        ProductPriceData     `json:"priceData"`
	StockStatus      string               `json:"stockStatus,omitempty"`
	StockQuantity    int                  `json:"stockQuantity"`
	ShippingPayer    string               `json:"shippingPayer"`
	Categories       []ProductCategoryRef `json:"categories,omitempty"`
	Variants         []ProductVariant     `json:"variants,omitempty"`
	Options          []ProductOption      `json:"options,omitempty"`
	SingleOption     bool                 `json:"singleOption"`
	CustomListing    bool                 `json:"customListing"`
	CustomNote       string               `json:"customNote,omitempty"`
	PlacementScore   int                  `json:"placementScore,omitempty"`
	DispatchDuration int                  `json:"dispatchDuration,omitempty"`
}

// ProductMedia describes an image asset attached to a product listing.
type ProductMedia struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	Placement int    `json:"placement"`
}

// ProductPriceData holds currency, unit price, discounts, and shipping cost.
type ProductPriceData struct {
	Currency        string `json:"currency"`
	Price           string `json:"price"`
	Discount        bool   `json:"discount"`
	DiscountedPrice string `json:"discountedPrice,omitempty"`
	ShippingPrice   string `json:"shippingPrice,omitempty"`
}

// ProductCategoryRef links a product to an existing category.
type ProductCategoryRef struct {
	ID         string `json:"id,omitempty"`
	CategoryID string `json:"categoryId,omitempty"`
	Title      string `json:"title,omitempty"`
}

// ProductVariant defines variant stock and pricing tied to variation selections.
type ProductVariant struct {
	VariationID    string            `json:"variationId,omitempty"`
	VariationTitle string            `json:"variationTitle,omitempty"`
	SelectionID    any               `json:"selectionId,omitempty"` // array of strings on input, string or array on response
	SelectionTitle string            `json:"selectionTitle,omitempty"`
	StockStatus    string            `json:"stockStatus,omitempty"`
	StockQuantity  int               `json:"stockQuantity"`
	Media          []ProductMedia    `json:"media,omitempty"`
	PriceData      *ProductPriceData `json:"priceData,omitempty"`
	Primary        bool              `json:"primary,omitempty"`
}

// ProductOption defines an optional addon item configurable by the buyer.
type ProductOption struct {
	ID          string `json:"id,omitempty"`
	OptionID    string `json:"optionId,omitempty"`
	Title       string `json:"title,omitempty"`
	OptionTitle string `json:"optionTitle,omitempty"`
	Price       string `json:"price,omitempty"`
	OptionPrice string `json:"optionPrice,omitempty"`
}

// ProductCreateRequest defines parameters required by POST /products.
type ProductCreateRequest struct {
	Title            string               `json:"title"`
	Description      string               `json:"description,omitempty"`
	Type             string               `json:"type"`
	Media            []ProductMedia       `json:"media"`
	PriceData        ProductPriceData     `json:"priceData"`
	StockQuantity    int                  `json:"stockQuantity,omitempty"`
	ShippingPayer    string               `json:"shippingPayer"`
	Categories       []ProductCategoryRef `json:"categories,omitempty"`
	Variants         []ProductVariant     `json:"variants,omitempty"`
	Options          []ProductOption      `json:"options,omitempty"`
	SingleOption     bool                 `json:"singleOption,omitempty"`
	CustomListing    bool                 `json:"customListing,omitempty"`
	CustomNote       string               `json:"customNote,omitempty"`
	PlacementScore   int                  `json:"placementScore,omitempty"`
	DispatchDuration int                  `json:"dispatchDuration,omitempty"`
}

// ProductUpdateRequest defines parameters accepted by PUT /products/{id}.
type ProductUpdateRequest struct {
	Title            *string              `json:"title,omitempty"`
	Description      *string              `json:"description,omitempty"`
	Type             *string              `json:"type,omitempty"`
	Media            []ProductMedia       `json:"media,omitempty"`
	PriceData        *ProductPriceData    `json:"priceData,omitempty"`
	StockQuantity    *int                 `json:"stockQuantity,omitempty"`
	ShippingPayer    *string              `json:"shippingPayer,omitempty"`
	Categories       []ProductCategoryRef `json:"categories,omitempty"`
	Variants         []ProductVariant     `json:"variants,omitempty"`
	Options          []ProductOption      `json:"options,omitempty"`
	SingleOption     *bool                `json:"singleOption,omitempty"`
	CustomListing    *bool                `json:"customListing,omitempty"`
	CustomNote       *string              `json:"customNote,omitempty"`
	PlacementScore   *int                 `json:"placementScore,omitempty"`
	DispatchDuration *int                 `json:"dispatchDuration,omitempty"`
}

// ProductListOptions holds filter parameters accepted by GET /products.
type ProductListOptions struct {
	ListOptions
	Title    string `url:"title,omitempty"`
	Status   string `url:"status,omitempty"`
	Category string `url:"category,omitempty"`
}

// Values compiles ProductListOptions into URL query parameters.
func (o *ProductListOptions) Values() url.Values {
	var v url.Values
	if o != nil {
		v = o.ListOptions.Values()
		if o.Title != "" {
			v.Set("title", o.Title)
		}
		if o.Status != "" {
			v.Set("status", o.Status)
		}
		if o.Category != "" {
			v.Set("category", o.Category)
		}
	} else {
		v = make(url.Values)
	}
	return v
}

// ProductService exposes product catalog management endpoints.
type ProductService struct {
	client *Client
}

// List queries products matching the filter criteria.
func (s *ProductService) List(ctx context.Context, opts *ProductListOptions) (*PageResponse[Product], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Product
	header, err := s.client.execute(ctx, http.MethodGet, "/products", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Product]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// All returns a Go 1.23+ range iterator across all products matching opts.
func (s *ProductService) All(ctx context.Context, opts *ProductListOptions) iter.Seq2[Product, error] {
	return func(yield func(Product, error) bool) {
		var currentOpts ProductListOptions
		if opts != nil {
			currentOpts = *opts
		}
		if currentOpts.Page <= 0 {
			currentOpts.Page = 1
		}
		if currentOpts.Limit <= 0 {
			currentOpts.Limit = 50
		}

		for {
			res, err := s.List(ctx, &currentOpts)
			if err != nil {
				yield(Product{}, err)
				return
			}

			for _, item := range res.Items {
				if !yield(item, nil) {
					return
				}
			}

			if !res.Pagination.HasNextPage() {
				return
			}
			currentOpts.Page = res.Pagination.NextPage()
		}
	}
}

// Get retrieves a single product by its unique Shopier ID.
func (s *ProductService) Get(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: product id is required")
	}

	var product Product
	_, err := s.client.execute(ctx, http.MethodGet, "/products/"+id, nil, nil, &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Create publishes a new product to the shop.
func (s *ProductService) Create(ctx context.Context, req *ProductCreateRequest) (*Product, error) {
	if req == nil {
		return nil, fmt.Errorf("shopier: product create payload is required")
	}

	var created Product
	_, err := s.client.execute(ctx, http.MethodPost, "/products", nil, req, &created)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

// Update updates an existing product's attributes.
func (s *ProductService) Update(ctx context.Context, id string, req *ProductUpdateRequest) (*Product, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: product id is required")
	}
	if req == nil {
		return nil, fmt.Errorf("shopier: product update payload is required")
	}

	var updated Product
	_, err := s.client.execute(ctx, http.MethodPut, "/products/"+id, nil, req, &updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete permanently removes a product listing from the storefront.
func (s *ProductService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("shopier: product id is required")
	}

	_, err := s.client.execute(ctx, http.MethodDelete, "/products/"+id, nil, nil, nil)
	return err
}
