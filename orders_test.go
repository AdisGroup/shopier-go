package shopier_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdisGroup/shopier-go"
)

func TestOrderService_ListAndIterator(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" || page == "1" {
			w.Header().Set("Shopier-Pagination-Page", "1")
			w.Header().Set("Shopier-Pagination-Limit", "2")
			w.Header().Set("Shopier-Pagination-Total-Pages", "2")
			w.Header().Set("Shopier-Pagination-Total-Items", "3")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[
				{"id":"ord_1","status":"unfulfilled","currency":"TRY","totals":{"total":"100.00"},"shippingInfo":{"firstName":"Ali","lastName":"Yilmaz","email":"ali@example.com","phone":"+905551234567","address":"Istiklal","city":"Istanbul","country":"Turkiye"},"lineItems":[]},
				{"id":"ord_2","status":"fulfilled","currency":"TRY","totals":{"total":"200.00"},"shippingInfo":{"firstName":"Ayse","lastName":"Demir","email":"ayse@example.com","phone":"+905559876543","address":"Kordon","city":"Izmir","country":"Turkiye"},"lineItems":[]}
			]`))
			return
		}

		if page == "2" {
			w.Header().Set("Shopier-Pagination-Page", "2")
			w.Header().Set("Shopier-Pagination-Limit", "2")
			w.Header().Set("Shopier-Pagination-Total-Pages", "2")
			w.Header().Set("Shopier-Pagination-Total-Items", "3")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[
				{"id":"ord_3","status":"unfulfilled","currency":"TRY","totals":{"total":"300.00"},"shippingInfo":{"firstName":"Mehmet","lastName":"Kaya","email":"mehmet@example.com","phone":"+905553332211","address":"Tunali","city":"Ankara","country":"Turkiye"},"lineItems":[]}
			]`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client, err := shopier.NewClient("test_token",
		shopier.WithBaseURL(server.URL),
		shopier.WithMaxRetries(0),
	)
	if err != nil {
		t.Fatalf("unexpected NewClient error: %v", err)
	}

	// Test List method with page metadata
	res, err := client.Orders.List(context.Background(), &shopier.OrderListOptions{
		ListOptions: shopier.ListOptions{Page: 1, Limit: 2},
	})
	if err != nil {
		t.Fatalf("Orders.List failed: %v", err)
	}

	if len(res.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(res.Items))
	}
	if !res.Pagination.HasNextPage() {
		t.Error("expected HasNextPage to be true")
	}
	if res.Pagination.TotalItems != 3 {
		t.Errorf("expected TotalItems 3, got %d", res.Pagination.TotalItems)
	}

	// Test Go 1.23+ range iterator across multiple pages
	var iteratedIDs []string
	for order, iterErr := range client.Orders.All(context.Background(), nil) {
		if iterErr != nil {
			t.Fatalf("iteration error: %v", iterErr)
		}
		iteratedIDs = append(iteratedIDs, order.ID)
	}

	if len(iteratedIDs) != 3 {
		t.Fatalf("expected 3 orders iterated, got %d: %v", len(iteratedIDs), iteratedIDs)
	}
	if iteratedIDs[0] != "ord_1" || iteratedIDs[1] != "ord_2" || iteratedIDs[2] != "ord_3" {
		t.Errorf("unexpected order sequence: %v", iteratedIDs)
	}
}

func TestOrderService_GetAndUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/orders/ord_123" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"ord_123","status":"unfulfilled","currency":"TRY","totals":{"total":"150.00"},"shippingInfo":{"firstName":"Deniz","lastName":"Akin","email":"deniz@example.com","phone":"+905550001122","address":"Bagdat Cad","city":"Istanbul","country":"Turkiye"},"lineItems":[]}`))
			return
		}
		if r.Method == http.MethodPut && r.URL.Path == "/orders/ord_123" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"ord_123","status":"fulfilled","currency":"TRY","totals":{"total":"150.00"},"shippingInfo":{"firstName":"Deniz","lastName":"Akin","email":"deniz@example.com","phone":"+905550001122","address":"Bagdat Cad","city":"Istanbul","country":"Turkiye"},"lineItems":[]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, _ := shopier.NewClient("token", shopier.WithBaseURL(server.URL), shopier.WithMaxRetries(0))

	ord, err := client.Orders.Get(context.Background(), "ord_123")
	if err != nil {
		t.Fatalf("Orders.Get failed: %v", err)
	}
	if ord.Status != "unfulfilled" {
		t.Errorf("expected status unfulfilled, got %s", ord.Status)
	}

	updated, err := client.Orders.Update(context.Background(), "ord_123", &shopier.OrderUpdateRequest{
		Fulfillments: &shopier.OrderFulfillments{
			ShippingCompany: "yurtici",
			TrackingNumber:  "1234567890",
		},
	})
	if err != nil {
		t.Fatalf("Orders.Update failed: %v", err)
	}
	if updated.Status != "fulfilled" {
		t.Errorf("expected status fulfilled, got %s", updated.Status)
	}
}
