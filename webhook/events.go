package webhook

import (
	"encoding/json"
	"fmt"

	"github.com/AdisGroup/shopier-go"
)

// Supported webhook event names dispatched by Shopier.
const (
	EventProductCreated      = "product.created"
	EventProductUpdated      = "product.updated"
	EventOrderCreated        = "order.created"
	EventOrderAddressUpdated = "order.addressUpdated"
	EventOrderFulfilled      = "order.fulfilled"
	EventRefundRequested     = "refund.requested"
	EventRefundUpdated       = "refund.updated"
)

// Header contains contextual metadata parsed from incoming Shopier webhook HTTP headers.
type Header struct {
	WebhookID  string
	Event      string
	Timestamp  int64
	Signature  string
	AccountID  string
	APIVersion string
}

// Event wraps an authenticated webhook notification with its typed payload models.
type Event struct {
	Header     Header
	RawPayload []byte
}

// Order decodes the raw JSON payload into an Order model.
// Valid for: order.created, order.addressUpdated, order.fulfilled.
func (e *Event) Order() (*shopier.Order, error) {
	var ord shopier.Order
	if err := json.Unmarshal(e.RawPayload, &ord); err != nil {
		return nil, fmt.Errorf("shopier/webhook: failed to unmarshal order payload: %w", err)
	}
	return &ord, nil
}

// Product decodes the raw JSON payload into a Product model.
// Valid for: product.created, product.updated.
func (e *Event) Product() (*shopier.Product, error) {
	var prod shopier.Product
	if err := json.Unmarshal(e.RawPayload, &prod); err != nil {
		return nil, fmt.Errorf("shopier/webhook: failed to unmarshal product payload: %w", err)
	}
	return &prod, nil
}

// Refund decodes the raw JSON payload into a Refund model.
// Valid for: refund.requested, refund.updated.
func (e *Event) Refund() (*shopier.Refund, error) {
	var ref shopier.Refund
	if err := json.Unmarshal(e.RawPayload, &ref); err != nil {
		return nil, fmt.Errorf("shopier/webhook: failed to unmarshal refund payload: %w", err)
	}
	return &ref, nil
}
