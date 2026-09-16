package shopier

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"net/url"
)

// Balance represents funds held in a specific currency in the merchant's account.
type Balance struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// Transaction represents a financial ledger entry (charge, fee, adjustment, payout item).
type Transaction struct {
	OrderID     string            `json:"orderId,omitempty"`
	Type        string            `json:"type"`
	Description string            `json:"description,omitempty"`
	DateCreated string            `json:"dateCreated"`
	Gross       TransactionAmount `json:"gross"`
	Fee         TransactionFee    `json:"fee"`
	Net         TransactionAmount `json:"net"`
}

// BalanceListTransactionsOptions holds query filters for GET /balance/transactions.
type BalanceListTransactionsOptions struct {
	ListOptions
	DateStart string `url:"dateStart,omitempty"`
	DateEnd   string `url:"dateEnd,omitempty"`
}

// Values compiles BalanceListTransactionsOptions into URL query parameters.
func (o *BalanceListTransactionsOptions) Values() url.Values {
	var v url.Values
	if o != nil {
		v = o.ListOptions.Values()
		if o.DateStart != "" {
			v.Set("dateStart", o.DateStart)
		}
		if o.DateEnd != "" {
			v.Set("dateEnd", o.DateEnd)
		}
	} else {
		v = make(url.Values)
	}
	return v
}

// BalanceService exposes balance inquiries and transaction ledger queries.
type BalanceService struct {
	client *Client
}

// Get retrieves current net balances across all active currencies.
func (s *BalanceService) Get(ctx context.Context) ([]Balance, error) {
	var balances []Balance
	_, err := s.client.execute(ctx, http.MethodGet, "/balance", nil, nil, &balances)
	if err != nil {
		return nil, err
	}
	return balances, nil
}

// ListTransactions returns paginated balance transactions contributing to net funds.
func (s *BalanceService) ListTransactions(ctx context.Context, opts *BalanceListTransactionsOptions) (*PageResponse[Transaction], error) {
	var query url.Values
	if opts != nil {
		query = opts.Values()
	}

	var items []Transaction
	header, err := s.client.execute(ctx, http.MethodGet, "/balance/transactions", query, nil, &items)
	if err != nil {
		return nil, err
	}

	return &PageResponse[Transaction]{
		Items:      items,
		Pagination: ParsePaginationHeaders(header),
	}, nil
}

// AllTransactions returns a Go 1.23+ iterator traversing all balance transactions across pages.
func (s *BalanceService) AllTransactions(ctx context.Context, opts *BalanceListTransactionsOptions) iter.Seq2[Transaction, error] {
	return func(yield func(Transaction, error) bool) {
		var currentOpts BalanceListTransactionsOptions
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
			res, err := s.ListTransactions(ctx, &currentOpts)
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

// GetTransaction retrieves a specific balance transaction by its associated order ID.
func (s *BalanceService) GetTransaction(ctx context.Context, orderID string) (*Transaction, error) {
	if orderID == "" {
		return nil, fmt.Errorf("shopier: order id is required")
	}

	var tx Transaction
	_, err := s.client.execute(ctx, http.MethodGet, "/balance/transactions/"+orderID, nil, nil, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
