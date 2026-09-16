package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// DiscountCode represents a promotional code redeemed at checkout.
type DiscountCode struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	DateCreated   string `json:"dateCreated"`
	Type          string `json:"type"` // "amount" or "percent"
	AmountOff     string `json:"amountOff,omitempty"`
	PercentOff    string `json:"percentOff,omitempty"`
	AmountMinimum string `json:"amountMinimum,omitempty"`
	Currency      string `json:"currency"`
	NumAvailable  int    `json:"numAvailable"`
	NumUsed       int    `json:"numUsed"`
	ExpiresAt     string `json:"expiresAt,omitempty"`
}

// DiscountCodeCreateRequest defines parameters for POST /discounts/codes.
type DiscountCodeCreateRequest struct {
	Code          string `json:"code"`
	Type          string `json:"type"` // "amount" or "percent"
	AmountOff     string `json:"amountOff,omitempty"`
	PercentOff    string `json:"percentOff,omitempty"`
	AmountMinimum string `json:"amountMinimum,omitempty"`
	Currency      string `json:"currency"`
	NumAvailable  int    `json:"numAvailable"`
	ExpiresAt     string `json:"expiresAt,omitempty"`
}

// DiscountCodeUpdateRequest defines parameters for PUT /discounts/codes/{id}.
type DiscountCodeUpdateRequest struct {
	Code          string `json:"code,omitempty"`
	NumAvailable  int    `json:"numAvailable,omitempty"`
	ExpiresAt     string `json:"expiresAt,omitempty"`
	AmountMinimum string `json:"amountMinimum,omitempty"`
}

// AutomaticDiscount represents a rule-based promotion applied without user entry.
type AutomaticDiscount struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Scope           string   `json:"scope"` // "all", "selectedProducts", "selectedCategories"
	ProductIDs      []string `json:"productIds,omitempty"`
	CategoryIDs     []string `json:"categoryIds,omitempty"`
	DateCreated     string   `json:"dateCreated"`
	Type            string   `json:"type"` // "amount" or "percent"
	AmountOff       string   `json:"amountOff,omitempty"`
	PercentOff      string   `json:"percentOff,omitempty"`
	Requirement     string   `json:"requirement"` // "amount" or "quantity"
	AmountMinimum   string   `json:"amountMinimum,omitempty"`
	QuantityMinimum int      `json:"quantityMinimum,omitempty"`
	Currency        string   `json:"currency"`
	StartsAt        string   `json:"startsAt,omitempty"`
	ExpiresAt       string   `json:"expiresAt,omitempty"`
}

// AutomaticDiscountCreateRequest defines parameters for POST /discounts/automatic.
type AutomaticDiscountCreateRequest struct {
	Title           string   `json:"title"`
	Scope           string   `json:"scope"`
	ProductIDs      []string `json:"productIds,omitempty"`
	CategoryIDs     []string `json:"categoryIds,omitempty"`
	Type            string   `json:"type"`
	AmountOff       string   `json:"amountOff,omitempty"`
	PercentOff      string   `json:"percentOff,omitempty"`
	Requirement     string   `json:"requirement"`
	AmountMinimum   string   `json:"amountMinimum,omitempty"`
	QuantityMinimum int      `json:"quantityMinimum,omitempty"`
	Currency        string   `json:"currency"`
	StartsAt        string   `json:"startsAt,omitempty"`
	ExpiresAt       string   `json:"expiresAt,omitempty"`
}

// AutomaticDiscountUpdateRequest defines parameters for PUT /discounts/automatic/{id}.
type AutomaticDiscountUpdateRequest struct {
	Title           string   `json:"title,omitempty"`
	Scope           string   `json:"scope,omitempty"`
	ProductIDs      []string `json:"productIds,omitempty"`
	CategoryIDs     []string `json:"categoryIds,omitempty"`
	AmountMinimum   string   `json:"amountMinimum,omitempty"`
	QuantityMinimum int      `json:"quantityMinimum,omitempty"`
	StartsAt        string   `json:"startsAt,omitempty"`
	ExpiresAt       string   `json:"expiresAt,omitempty"`
}

// DiscountService exposes promotional code and automatic discount endpoints.
type DiscountService struct {
	client *Client
}

// ListCodes returns paginated discount codes.
func (s *DiscountService) ListCodes(ctx context.Context, opts *ListOptions) (*PageResponse[DiscountCode], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []DiscountCode
	header, err := s.client.execute(ctx, http.MethodGet, "/discounts/codes", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[DiscountCode]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// AllCodes returns a Go 1.23+ iterator traversing all discount codes across pages.
func (s *DiscountService) AllCodes(ctx context.Context, opts *ListOptions) iter.Seq2[DiscountCode, error] {
	return func(yield func(DiscountCode, error) bool) {
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
			res, err := s.ListCodes(ctx, &currentOpts)
			if err != nil {
				yield(DiscountCode{}, err)
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

// GetCode retrieves a discount code by ID.
func (s *DiscountService) GetCode(ctx context.Context, id string) (*DiscountCode, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: discount code id is required")
	}

	var code DiscountCode
	_, err := s.client.execute(ctx, http.MethodGet, "/discounts/codes/"+id, nil, nil, &code)
	if err != nil {
		return nil, err
	}
	return &code, nil
}

// CreateCode publishes a new discount code.
func (s *DiscountService) CreateCode(ctx context.Context, req *DiscountCodeCreateRequest) (*DiscountCode, error) {
	if req == nil || req.Code == "" {
		return nil, fmt.Errorf("shopier: discount code payload is required")
	}

	var code DiscountCode
	_, err := s.client.execute(ctx, http.MethodPost, "/discounts/codes", nil, req, &code)
	if err != nil {
		return nil, err
	}
	return &code, nil
}

// UpdateCode updates attributes of an existing discount code.
func (s *DiscountService) UpdateCode(ctx context.Context, id string, req *DiscountCodeUpdateRequest) (*DiscountCode, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: discount code id is required")
	}
	if req == nil {
		return nil, fmt.Errorf("shopier: update payload is required")
	}

	var code DiscountCode
	_, err := s.client.execute(ctx, http.MethodPut, "/discounts/codes/"+id, nil, req, &code)
	if err != nil {
		return nil, err
	}
	return &code, nil
}

// DeleteCode permanently deletes a discount code.
func (s *DiscountService) DeleteCode(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("shopier: discount code id is required")
	}

	_, err := s.client.execute(ctx, http.MethodDelete, "/discounts/codes/"+id, nil, nil, nil)
	return err
}

// ListAutomatic returns paginated automatic discount campaigns.
func (s *DiscountService) ListAutomatic(ctx context.Context, opts *ListOptions) (*PageResponse[AutomaticDiscount], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []AutomaticDiscount
	header, err := s.client.execute(ctx, http.MethodGet, "/discounts/automatic", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[AutomaticDiscount]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// AllAutomatic returns a Go 1.23+ iterator traversing all automatic discounts.
func (s *DiscountService) AllAutomatic(ctx context.Context, opts *ListOptions) iter.Seq2[AutomaticDiscount, error] {
	return func(yield func(AutomaticDiscount, error) bool) {
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
			res, err := s.ListAutomatic(ctx, &currentOpts)
			if err != nil {
				yield(AutomaticDiscount{}, err)
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

// GetAutomatic retrieves an automatic discount by ID.
func (s *DiscountService) GetAutomatic(ctx context.Context, id string) (*AutomaticDiscount, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: automatic discount id is required")
	}

	var disc AutomaticDiscount
	_, err := s.client.execute(ctx, http.MethodGet, "/discounts/automatic/"+id, nil, nil, &disc)
	if err != nil {
		return nil, err
	}
	return &disc, nil
}

// CreateAutomatic creates a new automatic discount campaign.
func (s *DiscountService) CreateAutomatic(ctx context.Context, req *AutomaticDiscountCreateRequest) (*AutomaticDiscount, error) {
	if req == nil || req.Title == "" {
		return nil, fmt.Errorf("shopier: automatic discount title is required")
	}

	var disc AutomaticDiscount
	_, err := s.client.execute(ctx, http.MethodPost, "/discounts/automatic", nil, req, &disc)
	if err != nil {
		return nil, err
	}
	return &disc, nil
}

// UpdateAutomatic updates an existing automatic discount campaign.
func (s *DiscountService) UpdateAutomatic(ctx context.Context, id string, req *AutomaticDiscountUpdateRequest) (*AutomaticDiscount, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: automatic discount id is required")
	}
	if req == nil {
		return nil, fmt.Errorf("shopier: update payload is required")
	}

	var disc AutomaticDiscount
	_, err := s.client.execute(ctx, http.MethodPut, "/discounts/automatic/"+id, nil, req, &disc)
	if err != nil {
		return nil, err
	}
	return &disc, nil
}

// DeleteAutomatic permanently deletes an automatic discount campaign.
func (s *DiscountService) DeleteAutomatic(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("shopier: automatic discount id is required")
	}

	_, err := s.client.execute(ctx, http.MethodDelete, "/discounts/automatic/"+id, nil, nil, nil)
	return err
}
