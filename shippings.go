package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// Shipping represents shipment tracking and carrier dispatch status.
type Shipping struct {
	OrderID        string `json:"orderId"`
	Status         string `json:"status"`
	Method         string `json:"method"`
	Type           string `json:"type"`
	DateCreated    string `json:"dateCreated"`
	DateDispatched string `json:"dateDispatched,omitempty"`
	Company        string `json:"company"`
	Code           string `json:"code"`
	TrackingNumber string `json:"trackingNumber,omitempty"`
	TrackingURL    string `json:"trackingUrl,omitempty"`
	Size           string `json:"size,omitempty"`
	SizeUnit       string `json:"sizeUnit,omitempty"`
	Weight         string `json:"weight,omitempty"`
	WeightUnit     string `json:"weightUnit,omitempty"`
	Cost           string `json:"cost,omitempty"`
	Currency       string `json:"currency,omitempty"`
}

// ShippingCreateRequest defines parameters for POST /shippings.
type ShippingCreateRequest struct {
	OrderID string `json:"orderId"`
	Company string `json:"company"`
}

// ShippingListOptions holds query filters for GET /shippings.
type ShippingListOptions struct {
	ListOptions
	Status  string `url:"status,omitempty"`
	Method  string `url:"method,omitempty"`
	Type    string `url:"type,omitempty"`
	OrderID string `url:"orderId,omitempty"`
}

// Values compiles ShippingListOptions into URL query parameters.
func (o *ShippingListOptions) Values() url.Values {
	var v url.Values
	if o != nil {
		v = o.ListOptions.Values()
		if o.Status != "" {
			v.Set("status", o.Status)
		}
		if o.Method != "" {
			v.Set("method", o.Method)
		}
		if o.Type != "" {
			v.Set("type", o.Type)
		}
		if o.OrderID != "" {
			v.Set("orderId", o.OrderID)
		}
	} else {
		v = make(url.Values)
	}
	return v
}

// ShippingService exposes shipping code and tracking management endpoints.
type ShippingService struct {
	client *Client
}

// List returns paginated shipping records.
func (s *ShippingService) List(ctx context.Context, opts *ShippingListOptions) (*PageResponse[Shipping], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Shipping
	header, err := s.client.execute(ctx, http.MethodGet, "/shippings", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Shipping]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// All returns a Go 1.23+ iterator traversing all shipping records.
func (s *ShippingService) All(ctx context.Context, opts *ShippingListOptions) iter.Seq2[Shipping, error] {
	return func(yield func(Shipping, error) bool) {
		var currentOpts ShippingListOptions
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
				yield(Shipping{}, err)
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

// Get retrieves a shipping record by its unique shipping code.
func (s *ShippingService) Get(ctx context.Context, code string) (*Shipping, error) {
	if code == "" {
		return nil, fmt.Errorf("shopier: shipping code is required")
	}

	var shipping Shipping
	_, err := s.client.execute(ctx, http.MethodGet, "/shippings/"+code, nil, nil, &shipping)
	if err != nil {
		return nil, err
	}
	return &shipping, nil
}

// Create generates a new contracted shipping code for an order.
func (s *ShippingService) Create(ctx context.Context, req *ShippingCreateRequest) (*Shipping, error) {
	if req == nil || req.OrderID == "" || req.Company == "" {
		return nil, fmt.Errorf("shopier: orderId and company are required")
	}

	var shipping Shipping
	_, err := s.client.execute(ctx, http.MethodPost, "/shippings", nil, req, &shipping)
	if err != nil {
		return nil, err
	}
	return &shipping, nil
}

// Delete cancels an unmanifested shipping code.
func (s *ShippingService) Delete(ctx context.Context, code string) error {
	if code == "" {
		return fmt.Errorf("shopier: shipping code is required")
	}

	_, err := s.client.execute(ctx, http.MethodDelete, "/shippings/"+code, nil, nil, nil)
	return err
}
