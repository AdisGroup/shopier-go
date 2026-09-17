package shopier

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestQuickCheckout_SuccessFlow(t *testing.T) {
	mockCSRF := "mock-csrf-token-123456"
	mockFormCSRF := "mock-form-csrf-789012"
	expectedOrderID := "987654321"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/adisgroup/50900391":
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `<!DOCTYPE html><html><head><meta name="csrf-token" content="%s"></head><body><h1>Product</h1></body></html>`, mockCSRF)

		case "/s/api/v1/check_payment_progress/adisgroup":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if r.Header.Get("X-CSRF-Token") != mockCSRF {
				http.Error(w, "invalid csrf", http.StatusForbidden)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))

		case "/s/shipping/adisgroup":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `<!DOCTYPE html><html><head><meta name="csrf-token" content="%s"></head><body><form></form></body></html>`, mockFormCSRF)

		case "/s/api/v1/shipment_form_process/adisgroup":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if r.Header.Get("X-CSRF-Token") != mockFormCSRF {
				http.Error(w, "invalid form csrf", http.StatusForbidden)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"order_id": "%s", "status": "ok"}`, expectedOrderID)

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	svc := NewQuickCheckoutService(
		WithQuickCheckoutBaseURL(server.URL),
		WithQuickCheckoutTimeout(5*time.Second),
	)

	ctx := context.Background()
	req := &QuickCheckoutRequest{
		ProductURL: fmt.Sprintf("%s/adisgroup/50900391", server.URL),
		Buyer: QuickCheckoutBuyer{
			Email:     "buyer@example.com",
			FirstName: "Ahmet",
			LastName:  "Yılmaz",
			Phone:     "+90 532 111 22 33",
			Country:   "Türkiye",
			City:      "İstanbul",
		},
	}

	res, err := svc.Create(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.OrderID != expectedOrderID {
		t.Errorf("expected order_id %q, got %q", expectedOrderID, res.OrderID)
	}

	expectedPaymentURL := fmt.Sprintf("%s/s/payment/adisgroup/%s", server.URL, expectedOrderID)
	if res.PaymentURL != expectedPaymentURL {
		t.Errorf("expected payment_url %q, got %q", expectedPaymentURL, res.PaymentURL)
	}
}

func TestQuickCheckout_ExplicitAccountAndProductID(t *testing.T) {
	mockCSRF := "csrf-token-abc"
	expectedOrderID := "112233"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/myshop/item100":
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<html><head><meta content="%s" name="csrf-token"></head></html>`, mockCSRF)
		case "/s/api/v1/check_payment_progress/myshop":
			w.WriteHeader(http.StatusOK)
		case "/s/shipping/myshop":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<html><body>No meta tag</body></html>`))
		case "/s/api/v1/shipment_form_process/myshop":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"order_id": "%s"}`, expectedOrderID)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := NewClient("dummy_token", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to init client: %v", err)
	}
	client.QuickCheckout.baseURL = server.URL

	ctx := context.Background()
	req := &QuickCheckoutRequest{
		Account:   "myshop",
		ProductID: "item100",
		Quantity:  2,
		Buyer: QuickCheckoutBuyer{
			Email:     "buyer@example.com",
			FirstName: "Mehmet",
			LastName:  "Demir",
			Phone:     "05329998877",
		},
	}

	res, err := client.QuickCheckout.Create(ctx, req)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	if res.OrderID != expectedOrderID {
		t.Errorf("expected order_id %q, got %q", expectedOrderID, res.OrderID)
	}
}

func TestQuickCheckout_ValidationErrors(t *testing.T) {
	svc := NewQuickCheckoutService()
	ctx := context.Background()

	// 1. Nil request
	_, err := svc.Create(ctx, nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}

	// 2. Missing Account and ProductID
	_, err = svc.Create(ctx, &QuickCheckoutRequest{
		Buyer: QuickCheckoutBuyer{Email: "a@b.com", FirstName: "A", LastName: "B", Phone: "123"},
	})
	if err == nil {
		t.Fatal("expected error for missing product identification, got nil")
	}

	// 3. Missing Buyer fields
	_, err = svc.Create(ctx, &QuickCheckoutRequest{
		Account:   "shop",
		ProductID: "123",
	})
	if err == nil {
		t.Fatal("expected error for missing buyer details, got nil")
	}
}

func TestQuickCheckout_CSRFNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>No CSRF</title></head></html>`))
	}))
	defer server.Close()

	svc := NewQuickCheckoutService(WithQuickCheckoutBaseURL(server.URL))
	ctx := context.Background()

	req := &QuickCheckoutRequest{
		Account:   "shop",
		ProductID: "123",
		Buyer: QuickCheckoutBuyer{
			Email:     "a@b.com",
			FirstName: "A",
			LastName:  "B",
			Phone:     "123",
		},
	}

	_, err := svc.Create(ctx, req)
	if err == nil {
		t.Fatal("expected error when CSRF tag is missing, got nil")
	}

	var qErr *QuickCheckoutError
	if !errors.As(err, &qErr) {
		t.Fatalf("expected *QuickCheckoutError, got %T", err)
	}

	if qErr.Stage != StageExtractCSRF {
		t.Errorf("expected stage %s, got %s", StageExtractCSRF, qErr.Stage)
	}
}

func TestQuickCheckout_ConcurrencyIsolation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/s/api/v1/shipment_form_process/concurrent-shop":
			r.ParseForm()
			email := r.Form.Get("Email")
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"order_id": "ord_%s"}`, email)
		case r.URL.Path == "/s/api/v1/check_payment_progress/concurrent-shop":
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/s/shipping/concurrent-shop":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<html><meta name="csrf-token" content="tok"></html>`))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<html><meta name="csrf-token" content="tok"></html>`))
		}
	}))
	defer server.Close()

	svc := NewQuickCheckoutService(WithQuickCheckoutBaseURL(server.URL))
	ctx := context.Background()

	var wg sync.WaitGroup
	concurrentReqs := 10
	errChan := make(chan error, concurrentReqs)

	for i := 0; i < concurrentReqs; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			email := fmt.Sprintf("user%d@example.com", idx)
			res, err := svc.Create(ctx, &QuickCheckoutRequest{
				Account:   "concurrent-shop",
				ProductID: "item1",
				Buyer: QuickCheckoutBuyer{
					Email:     email,
					FirstName: "First",
					LastName:  "Last",
					Phone:     "5550001122",
				},
			})
			if err != nil {
				errChan <- err
				return
			}
			expectedID := fmt.Sprintf("ord_%s", email)
			if res.OrderID != expectedID {
				errChan <- fmt.Errorf("concurrency mismatch: expected %s, got %s", expectedID, res.OrderID)
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("concurrency test error: %v", err)
	}
}
