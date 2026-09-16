package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AdisGroup/shopier-go/webhook"
)

func main() {
	webhookSecret := os.Getenv("SHOPIER_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Fatal("SHOPIER_WEBHOOK_SECRET environment variable is required")
	}

	// 1. Standard net/http handler adapter
	webhookHandler := webhook.NewHandler(webhookSecret, func(ctx context.Context, evt *webhook.Event) error {
		log.Printf("Received Shopier Webhook: %s (ID: %s, Account: %s)",
			evt.Header.Event,
			evt.Header.WebhookID,
			evt.Header.AccountID,
		)

		switch evt.Header.Event {
		case webhook.EventOrderCreated:
			ord, err := evt.Order()
			if err != nil {
				return err
			}
			fmt.Printf("New Order Received! #%s - Total: %s %s - Buyer: %s %s\n",
				ord.ID,
				ord.Totals.Total,
				ord.Currency,
				ord.ShippingInfo.FirstName,
				ord.ShippingInfo.LastName,
			)

		case webhook.EventOrderFulfilled:
			ord, err := evt.Order()
			if err != nil {
				return err
			}
			fmt.Printf("Order Fulfilled: #%s\n", ord.ID)

		case webhook.EventRefundRequested, webhook.EventRefundUpdated:
			ref, err := evt.Refund()
			if err != nil {
				return err
			}
			fmt.Printf("Refund Notification: #%s for Order #%s - Amount: %s %s - Status: %s\n",
				ref.ID,
				ref.OrderID,
				ref.Total,
				ref.Currency,
				ref.Status,
			)
		}

		return nil
	})

	http.Handle("/webhooks/shopier", webhookHandler)

	log.Println("Webhook listener running on :8081. Notification URL: http://localhost:8081/webhooks/shopier")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
