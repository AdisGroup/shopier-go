package testutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/AdisGroup/shopier-go"
)

// MockServer wraps an httptest.Server pre-configured to simulate Shopier API responses.
type MockServer struct {
	*httptest.Server
	mux *http.ServeMux
}

// NewMockServer initializes a new mock HTTP server.
func NewMockServer() *MockServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	return &MockServer{
		Server: server,
		mux:    mux,
	}
}

// Handle registers an HTTP handler for a specific path pattern.
func (m *MockServer) Handle(pattern string, handler http.Handler) {
	m.mux.Handle(pattern, handler)
}

// HandleFunc registers a handler function for a specific path pattern.
func (m *MockServer) HandleFunc(pattern string, fn http.HandlerFunc) {
	m.mux.HandleFunc(pattern, fn)
}

// Client returns a Shopier Client pointed at this mock server.
func (m *MockServer) Client(opts ...shopier.Option) *shopier.Client {
	baseOpts := []shopier.Option{
		shopier.WithBaseURL(m.URL),
		shopier.WithOAuthBaseURL(m.URL),
		shopier.WithMaxRetries(0),
	}
	baseOpts = append(baseOpts, opts...)

	c, _ := shopier.NewClient("test_pat_token", baseOpts...)
	return c
}
