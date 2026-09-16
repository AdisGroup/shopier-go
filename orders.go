package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// Order represents an order placed through Shopier.
type Order struct {
	ID            string            `json:"id"`
	Status        string            `json:"status"`
	PaymentStatus string            `json:"paymentStatus"`
	Installments  bool              `json:"installments"`
	DateCreated   string            `json:"dateCreated"`
	Currency      string            `json:"currency"`
	PaymentMethod string            `json:"paymentMethod"`
	Totals        OrderTotals       `json:"totals"`
	Discounts     []OrderDiscount   `json:"discounts,omitempty"`
	ShippingInfo  OrderShippingInfo `json:"shippingInfo"`
	BillingInfo   *OrderBillingInfo `json:"billingInfo,omitempty"`
	Note          string            `json:"note,omitempty"`
	LineItems     []OrderLineItem   `json:"lineItems"`
}

// OrderTotals specifies price breakdowns for an order.
type OrderTotals struct {
	Subtotal string `json:"subtotal"`
	Shipping string `json:"shipping"`
	Discount string `json:"discount"`
	Total    string `json:"total"`
}

// OrderDiscount describes a discount applied to an order.
type OrderDiscount struct {
	ID     string `json:"id"`
	Method string `json:"method"`
}

// OrderShippingInfo represents the recipient and destination delivery details.
type OrderShippingInfo struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	NationalID string `json:"nationalId,omitempty"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Company    string `json:"company,omitempty"`
	Address    string `json:"address"`
	District   string `json:"district"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	Postcode   string `json:"postcode,omitempty"`
	Country    string `json:"country"`
}

// OrderBillingInfo represents corporate or separate invoicing details.
type OrderBillingInfo struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	NationalID string `json:"nationalId,omitempty"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Company    string `json:"company,omitempty"`
	TaxOffice  string `json:"taxOffice,omitempty"`
	TaxNumber  string `json:"taxNumber,omitempty"`
	Address    string `json:"address"`
	District   string `json:"district"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	Postcode   string `json:"postcode,omitempty"`
	Country    string `json:"country"`
}

// OrderLineItem represents an individual item within an order.
type OrderLineItem struct {
	ProductID string                   `json:"productId"`
	Title     string                   `json:"title"`
	Type      string                   `json:"type"`
	Selection []OrderLineItemSelection `json:"selection,omitempty"`
	Options   []OrderLineItemOption    `json:"options,omitempty"`
	Quantity  int                      `json:"quantity"`
	Price     string                   `json:"price"`
}

// OrderLineItemSelection represents variant choices for a line item.
type OrderLineItemSelection struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	VariationTitle string `json:"variationTitle"`
}

// OrderLineItemOption represents buyer-selected addon options for a line item.
type OrderLineItemOption struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Price string `json:"price"`
}

// OrderFulfillments defines delivery execution status and carrier details.
type OrderFulfillments struct {
	ProductType     string `json:"productType,omitempty"`
	ShippingCompany string `json:"shippingCompany,omitempty"`
	TrackingNumber  string `json:"trackingNumber,omitempty"`
	Note            string `json:"note,omitempty"`
}

// OrderUpdateRequest defines payload parameters accepted by PUT /orders/{id}.
type OrderUpdateRequest struct {
	Fulfillments *OrderFulfillments `json:"fulfillments,omitempty"`
	ShippingInfo *OrderShippingInfo `json:"shippingInfo,omitempty"`
}

// OrderTransaction provides financial breakdown details for an order payment.
type OrderTransaction struct {
	OrderID     string             `json:"orderId"`
	Type        string             `json:"type"`
	Description string             `json:"description,omitempty"`
	DateCreated string             `json:"dateCreated"`
	Gross       TransactionAmount  `json:"gross"`
	Fee         TransactionFee     `json:"fee"`
	Net         TransactionAmount  `json:"net"`
}

// TransactionAmount contains amount details in original and settlement currencies.
type TransactionAmount struct {
	OriginCurrency string `json:"originCurrency"`
	OriginAmount   string `json:"originAmount"`
	Currency       string `json:"currency"`
	Amount         string `json:"amount"`
}

// TransactionFee contains fee deductions applied to a transaction.
type TransactionFee struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// OrderListOptions holds query filters accepted by GET /orders.
// Shopier defaults dateStart/dateEnd to the last 60 days when omitted.
type OrderListOptions struct {
	ListOptions
	DateStart         string `url:"dateStart,omitempty"`
	DateEnd           string `url:"dateEnd,omitempty"`
	FulfillmentStatus string `url:"fulfillmentStatus,omitempty"`
	RefundType        string `url:"refundType,omitempty"`
	CustomerEmail     string `url:"customerEmail,omitempty"`
	CustomerPhone     string `url:"customerPhone,omitempty"`
	ProductID         string `url:"productId,omitempty"`
}

// Values compiles OrderListOptions into URL query parameters.
func (o *OrderListOptions) Values() url.Values {
	var v url.Values
	if o != nil {
		v = o.ListOptions.Values()
		if o.DateStart != "" {
			v.Set("dateStart", o.DateStart)
		}
		if o.DateEnd != "" {
			v.Set("dateEnd", o.DateEnd)
		}
		if o.FulfillmentStatus != "" {
			v.Set("fulfillmentStatus", o.FulfillmentStatus)
		}
		if o.RefundType != "" {
			v.Set("refundType", o.RefundType)
		}
		if o.CustomerEmail != "" {
			v.Set("customerEmail", o.CustomerEmail)
		}
		if o.CustomerPhone != "" {
			v.Set("customerPhone", o.CustomerPhone)
		}
		if o.ProductID != "" {
			v.Set("productId", o.ProductID)
		}
	} else {
		v = make(url.Values)
	}
	return v
}

// OrderService exposes order management endpoints.
type OrderService struct {
	client *Client
}

// List queries orders matching the specified filter parameters.
func (s *OrderService) List(ctx context.Context, opts *OrderListOptions) (*PageResponse[Order], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Order
	header, err := s.client.execute(ctx, http.MethodGet, "/orders", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Order]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// All returns a Go 1.23+ range iterator traversing every order matching opts.
// Pages are fetched lazily as the loop progresses.
func (s *OrderService) All(ctx context.Context, opts *OrderListOptions) iter.Seq2[Order, error] {
	return func(yield func(Order, error) bool) {
		var currentOpts OrderListOptions
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
				yield(Order{}, err)
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

// Get retrieves an order by its unique Shopier ID.
func (s *OrderService) Get(ctx context.Context, id string) (*Order, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: order id is required")
	}

	var order Order
	_, err := s.client.execute(ctx, http.MethodGet, "/orders/"+id, nil, nil, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Update updates an order's fulfillment tracking or buyer shipping address.
func (s *OrderService) Update(ctx context.Context, id string, req *OrderUpdateRequest) (*Order, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: order id is required")
	}
	if req == nil {
		return nil, fmt.Errorf("shopier: update request payload is required")
	}

	var updated Order
	_, err := s.client.execute(ctx, http.MethodPut, "/orders/"+id, nil, req, &updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// GetTransaction retrieves settlement and fee breakdown for an order.
func (s *OrderService) GetTransaction(ctx context.Context, orderID string) (*OrderTransaction, error) {
	if orderID == "" {
		return nil, fmt.Errorf("shopier: order id is required")
	}

	var tx OrderTransaction
	_, err := s.client.execute(ctx, http.MethodGet, "/orders/transactions/"+orderID, nil, nil, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
