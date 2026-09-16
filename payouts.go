package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// Payout represents a fund disbursement transferred to the merchant's bank account.
type Payout struct {
	ID          string            `json:"id"`
	Status      string            `json:"status"` // "pending" or "paid"
	DateCreated string            `json:"dateCreated"`
	Amount      string            `json:"amount"`
	Currency    string            `json:"currency"`
	Destination PayoutDestination `json:"destination"`
}

// PayoutDestination specifies the bank details receiving transferred funds.
type PayoutDestination struct {
	Type string `json:"type"` // "bankAccount"
	IBAN string `json:"iban"`
}

// PayoutService exposes merchant disbursement endpoints.
type PayoutService struct {
	client *Client
}

// List returns paginated payout transfer records.
func (s *PayoutService) List(ctx context.Context, opts *ListOptions) (*PageResponse[Payout], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Payout
	header, err := s.client.execute(ctx, http.MethodGet, "/payouts", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Payout]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// All returns a Go 1.23+ iterator traversing all payouts across pages.
func (s *PayoutService) All(ctx context.Context, opts *ListOptions) iter.Seq2[Payout, error] {
	return func(yield func(Payout, error) bool) {
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
				yield(Payout{}, err)
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

// Get retrieves a payout record by its ID.
func (s *PayoutService) Get(ctx context.Context, id string) (*Payout, error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: payout id is required")
	}

	var payout Payout
	_, err := s.client.execute(ctx, http.MethodGet, "/payouts/"+id, nil, nil, &payout)
	if err != nil {
		return nil, err
	}
	return &payout, nil
}

// ListTransactions returns paginated transactions included in a specific payout.
func (s *PayoutService) ListTransactions(ctx context.Context, id string, opts *ListOptions) (*PageResponse[Transaction], error) {
	if id == "" {
		return nil, fmt.Errorf("shopier: payout id is required")
	}

	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Transaction
	header, err := s.client.execute(ctx, http.MethodGet, "/payouts/transactions/"+id, query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Transaction]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// AllTransactions returns a Go 1.23+ iterator traversing all transactions of a payout across pages.
func (s *PayoutService) AllTransactions(ctx context.Context, id string, opts *ListOptions) iter.Seq2[Transaction, error] {
	return func(yield func(Transaction, error) bool) {
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
			res, err := s.ListTransactions(ctx, id, &currentOpts)
			if err != nil {
				yield(Transaction{}, err)
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
