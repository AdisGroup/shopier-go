package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AdisGroup/shopier-go"
)

func main() {
	// QuickCheckoutService can be used standalone (without an API token)
	// or accessed via client.QuickCheckout from a configured shopier.Client.
	qc := shopier.NewQuickCheckoutService(
		shopier.WithQuickCheckoutTimeout(30 * time.Second),
	)

	ctx := context.Background()

	// 1. Prepare request with product URL or Account + ProductID
	req := &shopier.QuickCheckoutRequest{
		ProductURL: "https://www.shopier.com/adisgroup/50900391",
		Quantity:   1,
		Buyer: shopier.QuickCheckoutBuyer{
			Email:     "buyer@example.com",
			FirstName: "Özgün Deniz",
			LastName:  "Küçük",
			Phone:     "+90 532 477 02 10",
			Country:   "Türkiye",
			City:      "İstanbul",
			Address:   "Merkez Mah. No: 1",
			Comment:   "Lütfen mesai saatlerinde teslim ediniz.",
		},
	}

	fmt.Printf("Generating direct payment link for %s...\n", req.ProductURL)

	// 2. Execute 4-step frontend session to retrieve direct payment URL
	res, err := qc.Create(ctx, req)
	if err != nil {
		if qErr, ok := err.(*shopier.QuickCheckoutError); ok {
			log.Fatalf("Quick checkout failed at stage [%s]: HTTP %d - %s",
				qErr.Stage, qErr.StatusCode, qErr.Message)
		}
		log.Fatalf("Quick checkout failed: %v", err)
	}

	// 3. Output checkout result
	fmt.Printf("Order ID   : %s\n", res.OrderID)
	fmt.Printf("Payment URL: %s\n", res.PaymentURL)
}
