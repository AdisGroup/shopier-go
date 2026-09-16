package webhook

import (
	"context"
	"net/http"
)

// EventHandler is a callback invoked when a valid, authenticated Shopier event is received.
type EventHandler func(ctx context.Context, event *Event) error

// NewHandler wraps an EventHandler in a standard http.HandlerFunc.
// It verifies the Shopier-Signature header and returns HTTP 200 within Shopier's required 5-second deadline.
func NewHandler(token string, handler EventHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		event, err := Parse(r, token)
		if err != nil {
			switch err {
			case ErrMissingSignature, ErrMissingEvent, ErrInvalidSignature:
				http.Error(w, err.Error(), http.StatusUnauthorized)
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

		if err := handler(r.Context(), event); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}
