package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// Variation represents an attribute axis (e.g. Size, Color) applicable to products.
type Variation struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Placement int    `json:"placement"`
}

// VariationCreateRequest defines parameters for POST /variations.
type VariationCreateRequest struct {
	Title string `json:"title"`
}

// VariationUpdateRequest defines parameters for PUT /variations/{id}.
type VariationUpdateRequest struct {
	Title string `json:"title"`
}

// VariationService exposes product variation management endpoints.
type VariationService struct {
	client *Client
}

// List returns paginated product variations.
func (s *VariationService) List(ctx context.Context, opts *ListOptions) (*PageResponse[Variation], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Variation
	header, err := s.client.execute(ctx, http.MethodGet, "/variations", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Variation]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// All returns a Go 1.23+ iterator traversing variations across all pages.
func (s *VariationService) All(ctx context.Context, opts *ListOptions) iter.Seq2[Variation, error] {
	return func(yield func(Variation, error) bool) {
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
				yield(Variation{}, err)
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

// Get retrieves a variation by its ID.
func (s *VariationService) Get(ctx context.Context, id string) (*Variation, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: variation id is required")
	}

	var variation Variation
	_, err := s.client.execute(ctx, http.MethodGet, "/variations/"+id, nil, nil, &variation)
	if err != nil {
		return nil, err
	}
	return &variation, nil
}

// Create registers a new variation.
func (s *VariationService) Create(ctx context.Context, req *VariationCreateRequest) (*Variation, error) {
	if req == nil || req.Title == "" {
		return nil, fmt.Errorf("shopier: variation title is required")
	}

	var variation Variation
	_, err := s.client.execute(ctx, http.MethodPost, "/variations", nil, req, &variation)
	if err != nil {
		return nil, err
	}
	return &variation, nil
}

// Update updates an existing variation title.
func (s *VariationService) Update(ctx context.Context, id string, req *VariationUpdateRequest) (*Variation, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: variation id is required")
	}
	if req == nil || req.Title == "" {
		return nil, fmt.Errorf("shopier: variation title is required")
	}

	var variation Variation
	_, err := s.client.execute(ctx, http.MethodPut, "/variations/"+id, nil, req, &variation)
	if err != nil {
		return nil, err
	}
	return &variation, nil
}

// Delete permanently removes a variation.
func (s *VariationService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("shopier: variation id is required")
	}

	_, err := s.client.execute(ctx, http.MethodDelete, "/variations/"+id, nil, nil, nil)
	return err
}
