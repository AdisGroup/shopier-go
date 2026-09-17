# Hakedişler (Payouts API)

Hakediş servisi, mağazanızın banka hesabına aktarılan periyodik ödemeleri ve bu ödemeleri oluşturan işlem detaylarını listelemenizi sağlar.

## Metotlar

- `client.Payouts.List(ctx, opts)` -> `(*PageResponse[Payout], error)`
- `client.Payouts.All(ctx, opts)` -> `iter.Seq2[Payout, error]`
- `client.Payouts.Get(ctx, id)` -> `(*Payout, error)`
- `client.Payouts.ListTransactions(ctx, id, opts)` -> `(*PageResponse[Transaction], error)`
- `client.Payouts.AllTransactions(ctx, id, opts)` -> `iter.Seq2[Transaction, error]`

---

## Örnek: Hakediş Ödemelerini Listeleme

```go
for payout, err := range client.Payouts.All(ctx, nil) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[%s] %s %s - Durum: %s (IBAN: %s)\n",
		payout.ID, payout.Amount, payout.Currency, payout.Status, payout.Destination.IBAN)
}
```

---

## Veri Modelleri ve Tipler

### `Payout` (Hakediş Modeli) {#hakedis-modeli-payout}

Satıcının banka hesabına transfer edilen ödeme kaydını temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Benzersiz hakediş transfer numarası |
| `Status` | `string` | `status` | Hakediş durumu (`"pending"` veya `"paid"`) |
| `Amount` | `string` | `amount` | Transfer edilen toplam net tutar |
| `Currency` | `string` | `currency` | Para birimi (`"TRY"`, `"USD"` vb.) |
| `DateCreated` | `string` | `dateCreated` | Hakediş aktarım tarihi |
| `Destination` | `PayoutDestination` | `destination` | Aktarılan banka hesap bilgisi |

---

### `PayoutDestination`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `Type` | `string` | `type` | Hesap türü (`"bankAccount"`) |
| `IBAN` | `string` | `iban` | Maskelenmiş alıcı IBAN numarası |
