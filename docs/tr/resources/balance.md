# Bakiye (Balance API)

Bakiye servisi, mağazanızın kullanılabilir net bakiyesini ve bakiyeye katkı sağlayan hesap hareketlerini sorgulamanızı sağlar.

## Metotlar

- `client.Balance.Get(ctx)` -> `([]Balance, error)`
- `client.Balance.ListTransactions(ctx, opts)` -> `(*PageResponse[Transaction], error)`
- `client.Balance.AllTransactions(ctx, opts)` -> `iter.Seq2[Transaction, error]`
- `client.Balance.GetTransaction(ctx, orderID)` -> `(*Transaction, error)`

---

## 1. Güncel Bakiyeyi Sorgulama

```go
balances, err := client.Balance.Get(ctx)
if err != nil {
	log.Fatal(err)
}

for _, b := range balances {
	fmt.Printf("Kullanılabilir Tutar: %s %s\n", b.Amount, b.Currency)
}
```

## 2. Bakiye Hareketlerini Listeleme

```go
opts := &shopier.BalanceListTransactionsOptions{
	DateStart: "2026-01-01T00:00:00Z",
}

for tx, err := range client.Balance.AllTransactions(ctx, opts) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("İşlem #%s | Tip: %s | Net: %s %s\n",
		tx.OrderID, tx.Type, tx.Net.Amount, tx.Net.Currency)
}
```

---

## Veri Modelleri ve Tipler

### `Balance` (Bakiye Modeli) {#bakiye-modeli-balance}

Belirli bir para birimindeki kullanılabilir mağaza bakiyesini temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `Currency` | `string` | `currency` | 3 haneli para birimi kodu (örn. `"TRY"`, `"USD"`) |
| `Amount` | `string` | `amount` | Kullanılabilir bakiye tutarı (örn. `"1250.50"`) |

---

### `Transaction` (Bakiye Hareketi Modeli) {#bakiye-hareketi-modeli-transaction}

Mağaza cari hesabına yansıyan bir finansal işlemi temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `OrderID` | `string` | `orderId` | İlgili sipariş numarası |
| `Type` | `string` | `type` | İşlem tipi (`"sale"`, `"refund"`, `"adjustment"`) |
| `Description` | `string` | `description` | İşlem açıklaması |
| `DateCreated` | `string` | `dateCreated` | ISO-8601 işlem tarihi |
| `Gross` | `TransactionAmount` | `gross` | İşlem brüt tutarı |
| `Fee` | `TransactionFee` | `fee` | Shopier komisyon ve hizmet kesintisi |
| `Net` | `TransactionAmount` | `net` | Hesaba yansıyan net tutar |

