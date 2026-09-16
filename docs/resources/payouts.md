# Payouts API Reference

The Payouts API tracks merchant bank disbursements and settlements.

## Methods

- `client.Payouts.List(ctx, opts)` -> `(*PageResponse[Payout], error)`
- `client.Payouts.All(ctx, opts)` -> `iter.Seq2[Payout, error]`
- `client.Payouts.Get(ctx, id)` -> `(*Payout, error)`
- `client.Payouts.ListTransactions(ctx, id, opts)` -> `(*PageResponse[Transaction], error)`
- `client.Payouts.AllTransactions(ctx, id, opts)` -> `iter.Seq2[Transaction, error]`

---

## Example: Listing Payouts

```go
for payout, err := range client.Payouts.All(ctx, nil) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[%s] %s %s - Status: %s (IBAN: %s)\n",
		payout.ID, payout.Amount, payout.Currency, payout.Status, payout.Destination.IBAN)
}
```
