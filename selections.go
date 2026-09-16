package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// Selection represents a concrete option within an attribute variation (e.g. "Red" in "Color").
type Selection struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	VariationID string `json:"variationId"`
}

// SelectionCreateRequest defines parameters for POST /selections.
type SelectionCreateRequest struct {
	VariationID string `json:"variationId"`
	Title       string `json:"title"`
}

// SelectionUpdateRequest defines parameters for PUT /selections/{id}.
type SelectionUpdateRequest struct {
	Title string `json:"title"`
}

// SelectionListOptions holds query filters for GET /selections.
type SelectionListOptions struct {
	ListOptions
	VariationIDs []string `url:"variationId,omitempty"`
}

// Values compiles SelectionListOptions into URL query parameters.
func (o *SelectionListOptions) Values() url.Values {
	var v url.Values
	if o != nil {
		v = o.ListOptions.Values()
		for _, vid := range o.VariationIDs {
			v.Add("variationId", vid)
		}
	} else {
		v = make(url.Values)
	}
	return v
}

// SelectionService exposes variation selection endpoints.
type SelectionService struct {
	client *Client
}

// List returns paginated variation selections.
func (s *SelectionService) List(ctx context.Context, opts *SelectionListOptions) (*PageResponse[Selection], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Selection
	header, err := s.client.execute(ctx, http.MethodGet, "/selections", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Selection]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// All returns a Go 1.23+ iterator traversing selections across all pages.
func (s *SelectionService) All(ctx context.Context, opts *SelectionListOptions) iter.Seq2[Selection, error] {
	return func(yield func(Selection, error) bool) {
		var currentOpts SelectionListOptions
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
				yield(Selection{}, err)
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

// Get retrieves a selection by its ID.
func (s *SelectionService) Get(ctx context.Context, id string) (*Selection, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: selection id is required")
	}

	var sel Selection
	_, err := s.client.execute(ctx, http.MethodGet, "/selections/"+id, nil, nil, &sel)
	if err != nil {
		return nil, err
	}
	return &sel, nil
}

// Create adds a new option value to a variation.
func (s *SelectionService) Create(ctx context.Context, req *SelectionCreateRequest) (*Selection, error) {
	if req == nil || req.VariationID == "" || req.Title == "" {
		return nil, fmt.Errorf("shopier: variationId and title are required")
	}

	var sel Selection
	_, err := s.client.execute(ctx, http.MethodPost, "/selections", nil, req, &sel)
	if err != nil {
		return nil, err
	}
	return &sel, nil
}

// Update renames an existing selection.
func (s *SelectionService) Update(ctx context.Context, id string, req *SelectionUpdateRequest) (*Selection, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: selection id is required")
	}
	if req == nil || req.Title == "" {
		return nil, fmt.Errorf("shopier: selection title is required")
	}

	var sel Selection
	_, err := s.client.execute(ctx, http.MethodPut, "/selections/"+id, nil, req, &sel)
	if err != nil {
		return nil, err
	}
	return &sel, nil
}

// Delete permanently removes a selection.
func (s *SelectionService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("shopier: selection id is required")
	}

	_, err := s.client.execute(ctx, http.MethodDelete, "/selections/"+id, nil, nil, nil)
	return err
}
