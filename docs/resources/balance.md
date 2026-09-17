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

---

## Data Models & Types

### `Balance` {#balance-model}

Represents available account balance in a specific currency.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `Currency` | `string` | `currency` | 3-letter currency code (e.g. `"TRY"`, `"USD"`) |
| `Amount` | `string` | `amount` | Current available balance string (e.g. `"1250.50"`) |

---

### `Transaction` {#transaction-model}

Represents a financial movement in the merchant's Shopier ledger.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `OrderID` | `string` | `orderId` | Associated order identifier |
| `Type` | `string` | `type` | Transaction type (`"sale"`, `"refund"`, `"adjustment"`) |
| `Description` | `string` | `description` | Explanatory description |
| `DateCreated` | `string` | `dateCreated` | ISO-8601 ledger entry timestamp |
| `Gross` | `TransactionAmount` | `gross` | Gross incoming amount |
| `Fee` | `TransactionFee` | `fee` | Deduction / commission fee |
| `Net` | `TransactionAmount` | `net` | Net balance credited or debited |

