package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AdisGroup/shopier-go"
)

func main() {
	patToken := os.Getenv("SHOPIER_PAT_TOKEN")
	if patToken == "" {
		log.Fatal("SHOPIER_PAT_TOKEN environment variable is required")
	}

	client, err := shopier.NewClient(patToken)
	if err != nil {
		log.Fatalf("Client initialization failed: %v", err)
	}

	ctx := context.Background()

	fmt.Println("Iterating through all orders using Go 1.23+ Range-over-func...")

	// client.Orders.All handles pagination under the hood, requesting subsequent pages
	// only when the loop reaches the end of current page items.
	orderCount := 0
	for order, err := range client.Orders.All(ctx, nil) {
		if err != nil {
			log.Fatalf("Error encountered during iteration: %v", err)
		}

		orderCount++
		fmt.Printf("[%d] Order #%s | Status: %s | Total: %s %s\n",
			orderCount,
			order.ID,
			order.Status,
			order.Totals.Total,
			order.Currency,
		)
	}

	fmt.Printf("Successfully iterated over all %d orders.\n", orderCount)
}
