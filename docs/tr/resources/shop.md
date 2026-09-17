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

---

## Veri Modelleri ve Tipler

### `ShopOwner` (Mağaza Sahibi Modeli) {#magaza-sahibi-modeli-shopowner}

Satıcının yasal kimlik, iletişim ve banka hakediş hesap bilgilerini temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `id` | Benzersiz satıcı hesap kimliği |
| `Type` | `string` | `type` | Hesap türü (`"personal"` bireysel veya `"business"` kurumsal) |
| `FirstName` | `string` | `firstName` | Satıcı adı |
| `LastName` | `string` | `lastName` | Satıcı soyadı |
| `Contact` | `ShopOwnerContact` | `contact` | İletişim ve adres bilgileri |
| `Company` | `*ShopOwnerCompany` | `company` | Şirket unvanı ve vergi dairesi bilgileri |
| `BankAccount` | `ShopOwnerBankAccount` | `bankAccount` | Tanımlı hakediş IBAN bilgisi |

---

### `ShopSettings` (Mağaza Ayarları Modeli) {#magaza-ayarlari-modeli-shopsettings}

Mağaza vitrininin görünüm, dil, sepet ve operasyonel ayarlarını temsil eder.

| Alan (Field) | Tip | JSON Etiketi | Açıklama |
| :--- | :--- | :--- | :--- |
| `Name` | `string` | `name` | Mağaza kullanıcı adı / slug (örn. `"adisgroup"`) |
| `URL` | `string` | `url` | Mağaza vitrin web adresi |
| `Title` | `string` | `title` | Vitrin başlığı |
| `Slogan` | `string` | `slogan` | Mağaza sloganı |
| `Announcement` | `string` | `announcement` | Tepe duyuru bandı metni |
| `Email` | `string` | `email` | Müşteri destek e-postası |
| `Phone` | `string` | `phone` | Müşteri destek telefonu |
| `Language` | `string` | `language` | Mağaza dili (`"TR"` veya `"EN"`) |
| `Vacation` | `bool` | `vacation` | Tatil modu açık mı? |
| `Cart` | `bool` | `cart` | Çoklu ürün sepeti aktif mi? |
| `MobileView` | `string` | `mobileView` | Mobil görünüm düzeni (`"singleColumn"`, `"doubleColumn"`) |
| `OutOfStock` | `bool` | `outOfStock` | Tükendiğinde ürünleri gizleme/gösterme ayarı |

