package testutil_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/AdisGroup/shopier-go/testutil"
)

func TestMockServer(t *testing.T) {
	server := testutil.NewMockServer()
	defer server.Close()

	server.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"currency":"TRY","amount":"2500.00"}]`))
	})

	client := server.Client()
	balances, err := client.Balance.Get(context.Background())
	if err != nil {
		t.Fatalf("Balance.Get failed on mock server: %v", err)
	}

	if len(balances) != 1 || balances[0].Amount != "2500.00" {
		t.Fatalf("unexpected balances: %+v", balances)
	}
}
