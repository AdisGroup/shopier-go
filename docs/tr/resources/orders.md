# Siparişler (Orders API)

Sipariş servisi, mağazanız üzerinden oluşturulan tüm siparişleri listelemenizi, filtrelemenizi ve kargo takip numarası ile güncellemenizi sağlar.

## Metotlar

- `client.Orders.List(ctx, opts)` -> `(*PageResponse[Order], error)`
- `client.Orders.All(ctx, opts)` -> `iter.Seq2[Order, error]`
- `client.Orders.Get(ctx, id)` -> `(*Order, error)`
- `client.Orders.Update(ctx, id, req)` -> `(*Order, error)`
- `client.Orders.GetTransaction(ctx, orderID)` -> `(*OrderTransaction, error)`

---

## 1. Siparişleri Listeleme

```go
opts := &shopier.OrderListOptions{
	ListOptions: shopier.ListOptions{
		Page:  1,
		Limit: 20,
	},
	FulfillmentStatus: "unfulfilled", // "unfulfilled" veya "fulfilled"
	RefundType:        "none",        // "none", "partial", "full"
	DateStart:         "2026-01-01T00:00:00Z",
}

res, err := client.Orders.List(ctx, opts)
if err != nil {
	log.Fatal(err)
}

for _, ord := range res.Items {
	fmt.Printf("Sipariş #%s: Tutar %s %s\n", ord.ID, ord.Totals.Total, ord.Currency)
}
```

## 2. Tekil Sipariş Sorgulama

```go
order, err := client.Orders.Get(ctx, "123456789")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Alıcı: %s %s (%s)\n", order.ShippingInfo.FirstName, order.ShippingInfo.LastName, order.ShippingInfo.Email)
```

## 3. Sipariş Güncelleme (Kargolama)

```go
updated, err := client.Orders.Update(ctx, "123456789", &shopier.OrderUpdateRequest{
	Fulfillments: &shopier.OrderFulfillments{
		ShippingCompany: "yurtici", // yurtici, mng, ptt, aras, surat, ups, dhl vb.
		TrackingNumber:  "9876543210",
	},
})
if err != nil {
	log.Fatal(err)
}

fmt.Println("Yeni durum:", updated.Status) // "fulfilled"
```

## 4. Sipariş Finansal İşlemi (Transaction)

```go
tx, err := client.Orders.GetTransaction(ctx, "123456789")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Brüt: %s %s | Komisyon: %s %s | Net: %s %s\n",
	tx.Gross.Amount, tx.Gross.Currency,
	tx.Fee.Amount, tx.Fee.Currency,
	tx.Net.Amount, tx.Net.Currency,
)
```

---

## Veri Modelleri ve Tipler

### `Order` (Sipariş Modeli) {#siparis-modeli-order}

Shopier mağazanız üzerinden verilen bir siparişi temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Benzersiz Shopier sipariş numarası (örn. `"123456789"`) |
| `Status` | `string` | `status` | Teslimat/kargo durumu (`"unfulfilled"`, `"fulfilled"`, `"cancelled"`) |
| `PaymentStatus` | `string` | `paymentStatus` | Ödeme durumu (`"success"`, `"failed"`, `"pending"`) |
| `Installments` | `bool` | `installments` | Taksitli ödeme yapıldıysa `true` döner |
| `DateCreated` | `string` | `dateCreated` | ISO-8601 formatında sipariş oluşturulma tarihi |
| `Currency` | `string` | `currency` | 3 haneli para birimi kodu (`"TRY"`, `"USD"`, `"EUR"`) |
| `PaymentMethod` | `string` | `paymentMethod` | Ödeme yöntemi (`"credit_card"`) |
| `Totals` | [`OrderTotals`](#ordertotals) | `totals` | Sipariş tutar kırılımları |
| `Discounts` | `[]OrderDiscount` | `discounts` | Uygulanan indirim kodları veya otomatik indirimler |
| `ShippingInfo` | [`OrderShippingInfo`](#ordershippinginfo) | `shippingInfo` | Alıcı bilgileri ve teslimat adresi |
| `BillingInfo` | `*OrderBillingInfo` | `billingInfo` | Fatura adresi ve kurumsal vergi bilgileri |
| `Note` | `string` | `note` | Müşterinin ödeme sırasında bıraktığı sipariş notu |
| `LineItems` | `[]OrderLineItem` | `lineItems` | Sipariş edilen ürünler, varyantlar ve adetler |

---

### `OrderTotals`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `Subtotal` | `string` | `subtotal` | Ürünlerin indirim ve kargo öncesi ara toplamı (örn. `"200.00"`) |
| `Shipping` | `string` | `shipping` | Alıcıya yansıtılan kargo ücreti (örn. `"25.00"`) |
| `Discount` | `string` | `discount` | Toplam indirim tutarı (örn. `"20.00"`) |
| `Total` | `string` | `total` | Alıcının ödediği nihai toplam tutar (örn. `"205.00"`) |

---

### `OrderShippingInfo`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `FirstName` | `string` | `firstName` | Alıcı adı |
| `LastName` | `string` | `lastName` | Alıcı soyadı |
| `NationalID` | `string` | `nationalId` | T.C. Kimlik Numarası veya pasaport no |
| `Email` | `string` | `email` | Alıcı e-posta adresi |
| `Phone` | `string` | `phone` | Alıcı telefon numarası (örn. `"05551234567"`) |
| `Company` | `string` | `company` | Şirket / firma adı (varsa) |
| `Address` | `string` | `address` | Açık adres satırı |
| `District` | `string` | `district` | İlçe / mahalle |
| `City` | `string` | `city` | İl / şehir |
| `State` | `string` | `state` | Eyalet / bölge (yurtdışı siparişler için) |
| `Postcode` | `string` | `postcode` | Posta kodu |
| `Country` | `string` | `country` | 2 haneli ISO ülke kodu (`"TR"`, `"US"`) |

---

### `OrderBillingInfo`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `FirstName` | `string` | `firstName` | Fatura kesilecek kişi adı |
| `LastName` | `string` | `lastName` | Fatura kesilecek kişi soyadı |
| `TaxOffice` | `string` | `taxOffice` | Vergi dairesi adı |
| `TaxNumber` | `string` | `taxNumber` | Vergi kimlik numarası / TC |
| `Address` | `string` | `address` | Fatura adresi |
| `City` | `string` | `city` | Fatura şehri |
| `Country` | `string` | `country` | Fatura ülkesi |

---

### `OrderLineItem`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ProductID` | `string` | `productId` | Satın alınan Shopier ürün kimliği |
| `Title` | `string` | `title` | Ürün başlığı |
| `Type` | `string` | `type` | `"physical"` (fiziksel) veya `"digital"` (dijital) |
| `Quantity` | `int` | `quantity` | Satın alınan adet |
| `Price` | `string` | `price` | Birim satış fiyatı |
| `Selection` | `[]OrderLineItemSelection` | `selection` | Seçilen varyant bilgileri (beden, renk vb.) |
| `Options` | `[]OrderLineItemOption` | `options` | Özel ekstra seçenekler |

---

### `OrderUpdateRequest`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `Fulfillments` | `*OrderFulfillments` | `fulfillments` | Kargo firması (`"yurtici"`, `"mng"`, `"ptt"`, `"aras"`, `"surat"`, `"ups"`, `"dhl"`) ve takip numarası |
| `ShippingInfo` | `*OrderShippingInfo` | `shippingInfo` | Güncellenen alıcı teslimat adresi |

---

### `OrderTransaction`

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `OrderID` | `string` | `orderId` | İlgili sipariş numarası |
| `Type` | `string` | `type` | İşlem türü (`"sale"`, `"refund"`) |
| `DateCreated` | `string` | `dateCreated` | Muhasebe işlem tarihi |
| `Gross` | `TransactionAmount` | `gross` | Alıcının ödediği brüt tutar |
| `Fee` | `TransactionFee` | `fee` | Shopier komisyon ve işlem kesintisi |
| `Net` | `TransactionAmount` | `net` | Satıcı hesabına aktarılan net hakediş |

