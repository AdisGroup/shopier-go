package shopier_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdisGroup/shopier-go"
)

func TestProductService_CRUD(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/products":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"prod_new","title":"Awesome T-Shirt","type":"physical","media":[],"priceData":{"currency":"TRY","price":"250.00","discount":false},"shippingPayer":"sellerPays"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/products/prod_new":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"prod_new","title":"Awesome T-Shirt","type":"physical","media":[],"priceData":{"currency":"TRY","price":"250.00","discount":false},"shippingPayer":"sellerPays"}`))
		case r.Method == http.MethodPut && r.URL.Path == "/products/prod_new":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"prod_new","title":"Awesome T-Shirt V2","type":"physical","media":[],"priceData":{"currency":"TRY","price":"299.00","discount":false},"shippingPayer":"sellerPays"}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/products/prod_new":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, _ := shopier.NewClient("token", shopier.WithBaseURL(server.URL), shopier.WithMaxRetries(0))

	// Create
	created, err := client.Products.Create(context.Background(), &shopier.ProductCreateRequest{
		Title: "Awesome T-Shirt",
		Type:  "physical",
		PriceData: shopier.ProductPriceData{
			Currency: "TRY",
			Price:    "250.00",
		},
		ShippingPayer: "sellerPays",
	})
	if err != nil {
		t.Fatalf("Products.Create failed: %v", err)
	}
	if created.ID != "prod_new" {
		t.Errorf("expected id prod_new, got %s", created.ID)
	}

	// Get
	got, err := client.Products.Get(context.Background(), "prod_new")
	if err != nil {
		t.Fatalf("Products.Get failed: %v", err)
	}
	if got.Title != "Awesome T-Shirt" {
		t.Errorf("expected title Awesome T-Shirt, got %s", got.Title)
	}

	// Update
	newTitle := "Awesome T-Shirt V2"
	updated, err := client.Products.Update(context.Background(), "prod_new", &shopier.ProductUpdateRequest{
		Title: &newTitle,
	})
	if err != nil {
		t.Fatalf("Products.Update failed: %v", err)
	}
	if updated.Title != "Awesome T-Shirt V2" {
		t.Errorf("expected updated title, got %s", updated.Title)
	}

	// Delete
	if err := client.Products.Delete(context.Background(), "prod_new"); err != nil {
		t.Fatalf("Products.Delete failed: %v", err)
	}
}
