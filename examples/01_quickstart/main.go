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
		log.Fatalf("Failed to initialize Shopier client: %v", err)
	}

	ctx := context.Background()

	// 1. Query merchant balances
	balances, err := client.Balance.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch balances: %v", err)
	}
	for _, b := range balances {
		fmt.Printf("Balance: %s %s\n", b.Amount, b.Currency)
	}

	// 2. Fetch the most recent 5 orders
	ordersRes, err := client.Orders.List(ctx, &shopier.OrderListOptions{
		ListOptions: shopier.ListOptions{Limit: 5},
	})
	if err != nil {
		log.Fatalf("Failed to list orders: %v", err)
	}

	fmt.Printf("Fetched %d orders (Total: %d)\n", len(ordersRes.Items), ordersRes.Pagination.TotalItems)
	for _, ord := range ordersRes.Items {
		fmt.Printf("Order #%s: %s %s - Status: %s (Buyer: %s %s)\n",
			ord.ID,
			ord.Totals.Total,
			ord.Currency,
			ord.Status,
			ord.ShippingInfo.FirstName,
			ord.ShippingInfo.LastName,
		)
	}
}
