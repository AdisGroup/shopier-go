# Shop Settings & Profile API Reference

The Shop API manages merchant profile information and public storefront display preferences.

## Methods

- `client.Shop.GetOwner(ctx)` -> `(*ShopOwner, error)`
- `client.Shop.GetSettings(ctx)` -> `(*ShopSettings, error)`
- `client.Shop.UpdateSettings(ctx, req)` -> `(*ShopSettings, error)`

---

## 1. Retrieve Shop Owner Profile

```go
owner, err := client.Shop.GetOwner(ctx)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Merchant: %s %s (Type: %s, Email: %s)\n",
	owner.FirstName, owner.LastName, owner.Type, owner.Contact.Email)
```

## 2. Update Storefront Preferences

```go
newTitle := "Official Brand Store"
vacationMode := false

updated, err := client.Shop.UpdateSettings(ctx, &shopier.ShopSettingsUpdateRequest{
	Title:    &newTitle,
	Vacation: &vacationMode,
})
if err != nil {
	log.Fatal(err)
}

fmt.Println("Store Title:", updated.Title)
```

---

## Data Models & Types

### `ShopOwner` {#shopowner-model}

Represents identity, contact, and bank settlement details of the merchant.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Unique shop owner ID |
| `Type` | `string` | `type` | Account type (`"personal"` or `"business"`) |
| `FirstName` | `string` | `firstName` | Owner given name |
| `LastName` | `string` | `lastName` | Owner surname |
| `Contact` | `ShopOwnerContact` | `contact` | Address and telephone details |
| `Company` | `*ShopOwnerCompany` | `company` | Legal entity registration if business |
| `BankAccount` | `ShopOwnerBankAccount` | `bankAccount` | Registered settlement payout IBAN |

---

### `ShopSettings` {#shopsettings-model}

Encapsulates public storefront display and behavior configuration.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `Name` | `string` | `name` | Shop slug/subdomain (e.g. `"adisgroup"`) |
| `URL` | `string` | `url` | Full storefront URL |
| `Title` | `string` | `title` | Public shop title |
| `Slogan` | `string` | `slogan` | Shop slogan |
| `Announcement` | `string` | `announcement` | Top banner announcement text |
| `Email` | `string` | `email` | Customer service email |
| `Phone` | `string` | `phone` | Customer service phone |
| `Language` | `string` | `language` | Storefront language (`"TR"` or `"EN"`) |
| `Vacation` | `bool` | `vacation` | `true` if shop is temporarily paused in vacation mode |
| `Cart` | `bool` | `cart` | Multi-product shopping cart enabled |
| `MobileView` | `string` | `mobileView` | `"singleColumn"` or `"doubleColumn"` |
| `OutOfStock` | `bool` | `outOfStock` | Allow display of out-of-stock items |

