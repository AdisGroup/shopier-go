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
