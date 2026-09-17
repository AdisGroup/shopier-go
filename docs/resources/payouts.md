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

---

## Data Models & Types

### `Payout` {#payout-model}

Represents a bank settlement disbursement transferred to the merchant.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Unique payout transfer ID |
| `Status` | `string` | `status` | Payout status (`"pending"` or `"paid"`) |
| `Amount` | `string` | `amount` | Transferred total amount |
| `Currency` | `string` | `currency` | Currency code (`"TRY"`, `"USD"`, etc.) |
| `DateCreated` | `string` | `dateCreated` | ISO-8601 disbursement date |
| `Destination` | `PayoutDestination` | `destination` | Receiving bank destination details |

---

### `PayoutDestination`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `Type` | `string` | `type` | Destination category (`"bankAccount"`) |
| `IBAN` | `string` | `iban` | Masked destination IBAN string |
