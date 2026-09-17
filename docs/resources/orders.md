# Orders API Reference

The Orders API allows you to retrieve, filter, and update orders placed through Shopier storefronts.

## Methods

- `client.Orders.List(ctx, opts)` -> `(*PageResponse[Order], error)`
- `client.Orders.All(ctx, opts)` -> `iter.Seq2[Order, error]`
- `client.Orders.Get(ctx, id)` -> `(*Order, error)`
- `client.Orders.Update(ctx, id, req)` -> `(*Order, error)`
- `client.Orders.GetTransaction(ctx, orderID)` -> `(*OrderTransaction, error)`

---

## 1. List Orders

```go
opts := &shopier.OrderListOptions{
	ListOptions: shopier.ListOptions{
		Page:  1,
		Limit: 20,
	},
	FulfillmentStatus: "unfulfilled", // "unfulfilled" or "fulfilled"
	RefundType:        "none",        // "none", "partial", "full"
	DateStart:         "2026-01-01T00:00:00Z",
}

res, err := client.Orders.List(ctx, opts)
if err != nil {
	log.Fatal(err)
}

for _, ord := range res.Items {
	fmt.Printf("Order #%s: Total %s %s\n", ord.ID, ord.Totals.Total, ord.Currency)
}
```

## 2. Get Order

```go
order, err := client.Orders.Get(ctx, "123456789")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Buyer: %s %s (%s)\n", order.ShippingInfo.FirstName, order.ShippingInfo.LastName, order.ShippingInfo.Email)
```

## 3. Update Order (Fulfill & Track)

```go
updated, err := client.Orders.Update(ctx, "123456789", &shopier.OrderUpdateRequest{
	Fulfillments: &shopier.OrderFulfillments{
		ShippingCompany: "yurtici", // yurtici, mng, ptt, aras, surat, ups, dhl, etc.
		TrackingNumber:  "9876543210",
	},
})
if err != nil {
	log.Fatal(err)
}

fmt.Println("New status:", updated.Status) // "fulfilled"
```

## 4. Get Order Transaction

```go
tx, err := client.Orders.GetTransaction(ctx, "123456789")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Gross: %s %s | Fee: %s %s | Net: %s %s\n",
	tx.Gross.Amount, tx.Gross.Currency,
	tx.Fee.Amount, tx.Fee.Currency,
	tx.Net.Amount, tx.Net.Currency,
)
```

---

## Data Models & Types

### `Order` {#order-model}

Represents an order placed on a Shopier storefront.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Unique Shopier order identifier (e.g. `"123456789"`) |
| `Status` | `string` | `status` | Fulfillment status (`"unfulfilled"`, `"fulfilled"`, `"cancelled"`) |
| `PaymentStatus` | `string` | `paymentStatus` | Payment status (`"success"`, `"failed"`, `"pending"`) |
| `Installments` | `bool` | `installments` | `true` if paid via credit card installments |
| `DateCreated` | `string` | `dateCreated` | ISO-8601 creation timestamp |
| `Currency` | `string` | `currency` | 3-letter ISO currency code (`"TRY"`, `"USD"`, `"EUR"`) |
| `PaymentMethod` | `string` | `paymentMethod` | Method used (`"credit_card"`, etc.) |
| `Totals` | [`OrderTotals`](#ordertotals) | `totals` | Price breakdown and order sum |
| `Discounts` | `[]OrderDiscount` | `discounts` | Applied discount vouchers or automatic discounts |
| `ShippingInfo` | [`OrderShippingInfo`](#ordershippinginfo) | `shippingInfo` | Delivery address and recipient contact details |
| `BillingInfo` | `*OrderBillingInfo` | `billingInfo` | Invoicing information (corporate or individual) |
| `Note` | `string` | `note` | Optional customer note left during checkout |
| `LineItems` | `[]OrderLineItem` | `lineItems` | Purchased products, selected variants, and quantities |

---

### `OrderTotals`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `Subtotal` | `string` | `subtotal` | Sum of item prices before discounts/shipping (e.g. `"200.00"`) |
| `Shipping` | `string` | `shipping` | Shipping cost charged to buyer (e.g. `"25.00"`) |
| `Discount` | `string` | `discount` | Total discount amount deducted (e.g. `"20.00"`) |
| `Total` | `string` | `total` | Final payable amount (e.g. `"205.00"`) |

---

### `OrderShippingInfo`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `FirstName` | `string` | `firstName` | Recipient given name |
| `LastName` | `string` | `lastName` | Recipient surname |
| `NationalID` | `string` | `nationalId` | TC Identity Number (Turkiye) or passport number |
| `Email` | `string` | `email` | Buyer email address |
| `Phone` | `string` | `phone` | Recipient phone number (e.g. `"05551234567"`) |
| `Company` | `string` | `company` | Optional company name |
| `Address` | `string` | `address` | Street address line |
| `District` | `string` | `district` | District / neighborhood |
| `City` | `string` | `city` | City / province |
| `State` | `string` | `state` | State / province (international) |
| `Postcode` | `string` | `postcode` | Postal / zip code |
| `Country` | `string` | `country` | 2-letter ISO country code (`"TR"`, `"US"`) |

---

### `OrderBillingInfo`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `FirstName` | `string` | `firstName` | Invoicing given name |
| `LastName` | `string` | `lastName` | Invoicing surname |
| `TaxOffice` | `string` | `taxOffice` | Corporate tax administration office |
| `TaxNumber` | `string` | `taxNumber` | Corporate tax identification number |
| `Address` | `string` | `address` | Official billing street address |
| `City` | `string` | `city` | Billing city |
| `Country` | `string` | `country` | Billing country |

---

### `OrderLineItem`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `ProductID` | `string` | `productId` | Shopier product ID |
| `Title` | `string` | `title` | Product title at the time of purchase |
| `Type` | `string` | `type` | `"physical"` or `"digital"` |
| `Quantity` | `int` | `quantity` | Number of units purchased |
| `Price` | `string` | `price` | Unit price in store currency |
| `Selection` | `[]OrderLineItemSelection` | `selection` | Chosen variant attributes (size, color, etc.) |
| `Options` | `[]OrderLineItemOption` | `options` | Extra customized options |

---

### `OrderUpdateRequest`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `Fulfillments` | `*OrderFulfillments` | `fulfillments` | Tracking carrier (`"yurtici"`, `"mng"`, `"ptt"`, `"aras"`, `"surat"`, `"ups"`, `"dhl"`) and tracking number |
| `ShippingInfo` | `*OrderShippingInfo` | `shippingInfo` | Updated destination delivery address |

---

### `OrderTransaction`

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `OrderID` | `string` | `orderId` | Associated order identifier |
| `Type` | `string` | `type` | Transaction type (`"sale"`, `"refund"`) |
| `DateCreated` | `string` | `dateCreated` | Settlement creation timestamp |
| `Gross` | `TransactionAmount` | `gross` | Total transaction amount paid by buyer |
| `Fee` | `TransactionFee` | `fee` | Shopier commission and service deduction |
| `Net` | `TransactionAmount` | `net` | Net payable balance deposited to merchant |

