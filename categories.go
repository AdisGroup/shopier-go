package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// Category represents a grouping container for products.
type Category struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Placement int    `json:"placement"`
}

// CategoryCreateRequest defines parameters for POST /categories.
type CategoryCreateRequest struct {
	Title string `json:"title"`
}

// CategoryUpdateRequest defines parameters for PUT /categories/{id}.
type CategoryUpdateRequest struct {
	Title string `json:"title"`
}

// CategoryService exposes product category management endpoints.
type CategoryService struct {
	client *Client
}

// List returns paginated product categories.
func (s *CategoryService) List(ctx context.Context, opts *ListOptions) (*PageResponse[Category], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Category
	header, err := s.client.execute(ctx, http.MethodGet, "/categories", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Category]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// All returns a Go 1.23+ iterator traversing every category across all pages.
func (s *CategoryService) All(ctx context.Context, opts *ListOptions) iter.Seq2[Category, error] {
	return func(yield func(Category, error) bool) {
		var currentOpts ListOptions
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
				yield(Category{}, err)
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

// Get retrieves a category by its ID.
func (s *CategoryService) Get(ctx context.Context, id string) (*Category, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: category id is required")
	}

	var cat Category
	_, err := s.client.execute(ctx, http.MethodGet, "/categories/"+id, nil, nil, &cat)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// Create registers a new product category.
func (s *CategoryService) Create(ctx context.Context, req *CategoryCreateRequest) (*Category, error) {
	if req == nil || req.Title == "" {
		return nil, fmt.Errorf("shopier: category title is required")
	}

	var cat Category
	_, err := s.client.execute(ctx, http.MethodPost, "/categories", nil, req, &cat)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// Update renames an existing product category.
func (s *CategoryService) Update(ctx context.Context, id string, req *CategoryUpdateRequest) (*Category, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: category id is required")
	}
	if req == nil || req.Title == "" {
		return nil, fmt.Errorf("shopier: category title is required")
	}

	var cat Category
	_, err := s.client.execute(ctx, http.MethodPut, "/categories/"+id, nil, req, &cat)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// Delete permanently removes a product category.
func (s *CategoryService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("shopier: category id is required")
	}

	_, err := s.client.execute(ctx, http.MethodDelete, "/categories/"+id, nil, nil, nil)
	return err
}
