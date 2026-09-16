# Balance API Reference

The Balance API provides real-time access to available account balances and transaction ledger entries.

## Methods

- `client.Balance.Get(ctx)` -> `([]Balance, error)`
- `client.Balance.ListTransactions(ctx, opts)` -> `(*PageResponse[Transaction], error)`
- `client.Balance.AllTransactions(ctx, opts)` -> `iter.Seq2[Transaction, error]`
- `client.Balance.GetTransaction(ctx, orderID)` -> `(*Transaction, error)`

---

## 1. Query Current Balance

```go
balances, err := client.Balance.Get(ctx)
if err != nil {
	log.Fatal(err)
}

for _, b := range balances {
	fmt.Printf("Available: %s %s\n", b.Amount, b.Currency)
}
```

## 2. Stream Balance Transactions

```go
opts := &shopier.BalanceListTransactionsOptions{
	DateStart: "2026-01-01T00:00:00Z",
}

for tx, err := range client.Balance.AllTransactions(ctx, opts) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction #%s | Type: %s | Net: %s %s\n",
		tx.OrderID, tx.Type, tx.Net.Amount, tx.Net.Currency)
}
```
