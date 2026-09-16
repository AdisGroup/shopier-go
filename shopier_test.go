package shopier_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AdisGroup/shopier-go"
)

func TestNewClient_MissingToken(t *testing.T) {
	_, err := shopier.NewClient("")
	if !errors.Is(err, shopier.ErrMissingToken) {
		t.Fatalf("expected ErrMissingToken, got %v", err)
	}
}

func TestClient_RateLimit_RetryAfter(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"rate_limit_exceeded","message":"Rate limit reached"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"currency":"TRY","amount":"1500.00"}]`))
	}))
	defer server.Close()

	client, err := shopier.NewClient("test_token",
		shopier.WithBaseURL(server.URL),
		shopier.WithMaxRetries(2),
		shopier.WithRetryDelay(10*time.Millisecond, 50*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("unexpected NewClient error: %v", err)
	}

	balances, err := client.Balance.Get(context.Background())
	if err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}

	if len(balances) != 1 || balances[0].Amount != "1500.00" {
		t.Errorf("unexpected balances: %+v", balances)
	}

	if atomic.LoadInt32(&attempts) != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestClient_ServerError_Backoff(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attempts, 1)
		if count <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"message":"temporarily unavailable"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"prod_123","title":"Widget","type":"physical","media":[],"priceData":{"currency":"TRY","price":"10.00","discount":false},"shippingPayer":"sellerPays"}`))
	}))
	defer server.Close()

	client, err := shopier.NewClient("test_token",
		shopier.WithBaseURL(server.URL),
		shopier.WithMaxRetries(3),
		shopier.WithRetryDelay(5*time.Millisecond, 20*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("unexpected NewClient error: %v", err)
	}

	prod, err := client.Products.Get(context.Background(), "prod_123")
	if err != nil {
		t.Fatalf("expected recovery after 503 retries, got: %v", err)
	}

	if prod.ID != "prod_123" || prod.Title != "Widget" {
		t.Errorf("unexpected product data: %+v", prod)
	}

	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := shopier.NewClient("test_token",
		shopier.WithBaseURL(server.URL),
		shopier.WithMaxRetries(0),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = client.Balance.Get(ctx)
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
	}
}

func TestClient_APIError_Parsing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_parameter","message":"dateStart cannot be in future"}`))
	}))
	defer server.Close()

	client, err := shopier.NewClient("test_token",
		shopier.WithBaseURL(server.URL),
		shopier.WithMaxRetries(0),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.Orders.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected APIError, got nil")
	}

	var apiErr *shopier.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected error to be *shopier.APIError, got %T: %v", err, err)
	}

	if apiErr.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
	if apiErr.ErrorCode != "invalid_parameter" {
		t.Errorf("expected error code invalid_parameter, got %s", apiErr.ErrorCode)
	}
}
