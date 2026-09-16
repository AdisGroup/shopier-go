# Mağaza Ayarları (Shop API)

Mağaza servisi, mağaza sahibinin yasal/iletişim profilini ve vitrinin genel görünüm/çalışma ayarlarını yönetir.

## Metotlar

- `client.Shop.GetOwner(ctx)` -> `(*ShopOwner, error)`
- `client.Shop.GetSettings(ctx)` -> `(*ShopSettings, error)`
- `client.Shop.UpdateSettings(ctx, req)` -> `(*ShopSettings, error)`

---

## 1. Mağaza Sahibi Bilgilerini Alma

```go
owner, err := client.Shop.GetOwner(ctx)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Mağaza Sahibi: %s %s (Üyelik Tipi: %s, E-posta: %s)\n",
	owner.FirstName, owner.LastName, owner.Type, owner.Contact.Email)
```

## 2. Mağaza Ayarlarını Güncelleme

```go
newTitle := "Resmi Marka Mağazası"
tatilModu := false

updated, err := client.Shop.UpdateSettings(ctx, &shopier.ShopSettingsUpdateRequest{
	Title:    &newTitle,
	Vacation: &tatilModu,
})
if err != nil {
	log.Fatal(err)
}

fmt.Println("Mağaza Başlığı:", updated.Title)
```
