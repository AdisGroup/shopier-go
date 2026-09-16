package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// Refund represents a payment reimbursement for an order.
type Refund struct {
	ID           string `json:"id"`
	Type         string `json:"type"`   // "full" or "partial"
	Status       string `json:"status"` // "pending", "failed", "succeeded"
	OrderID      string `json:"orderId"`
	DateCreated  string `json:"dateCreated"`
	DateRefunded string `json:"dateRefunded,omitempty"`
	Currency     string `json:"currency"`
	Total        string `json:"total"`
	Note         string `json:"note,omitempty"`
}

// RefundCreateRequest defines parameters for POST /refunds.
type RefundCreateRequest struct {
	OrderID string `json:"orderId"`
	Amount  string `json:"amount,omitempty"`
	Note    string `json:"note,omitempty"`
}

// RefundListOptions holds query filters for GET /refunds.
type RefundListOptions struct {
	ListOptions
	OrderID string `url:"orderId,omitempty"`
}

// Values compiles RefundListOptions into URL query parameters.
func (o *RefundListOptions) Values() url.Values {
	var v url.Values
	if o != nil {
		v = o.ListOptions.Values()
		if o.OrderID != "" {
			v.Set("orderId", o.OrderID)
		}
	} else {
		v = make(url.Values)
	}
	return v
}

// RefundService exposes order refund management endpoints.
type RefundService struct {
	client *Client
}

// List returns paginated refund records.
func (s *RefundService) List(ctx context.Context, opts *RefundListOptions) (*PageResponse[Refund], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Refund
	header, err := s.client.execute(ctx, http.MethodGet, "/refunds", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Refund]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// All returns a Go 1.23+ iterator traversing all refunds across pages.
func (s *RefundService) All(ctx context.Context, opts *RefundListOptions) iter.Seq2[Refund, error] {
	return func(yield func(Refund, error) bool) {
		var currentOpts RefundListOptions
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
				yield(Refund{}, err)
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

// Get retrieves a refund by its ID.
func (s *RefundService) Get(ctx context.Context, id string) (*Refund, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: refund id is required")
	}

	var refund Refund
	_, err := s.client.execute(ctx, http.MethodGet, "/refunds/"+id, nil, nil, &refund)
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

// Create initiates a full or partial refund on a paid order.
func (s *RefundService) Create(ctx context.Context, req *RefundCreateRequest) (*Refund, error) {
	if req == nil || req.OrderID == "" {
		return nil, fmt.Errorf("shopier: orderId is required")
	}

	var refund Refund
	_, err := s.client.execute(ctx, http.MethodPost, "/refunds", nil, req, &refund)
	if err != nil {
		return nil, err
	}
	return &refund, nil
}
