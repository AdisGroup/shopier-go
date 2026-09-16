package shopier

import (
	"context"
	"net/http"
)

// ShopOwner represents identity and registration details of the merchant.
type ShopOwner struct {
	ID          string               `json:"id"`
	Type        string               `json:"type"` // "personal" or "business"
	FirstName   string               `json:"firstName"`
	LastName    string               `json:"lastName"`
	Contact     ShopOwnerContact     `json:"contact"`
	Company     *ShopOwnerCompany    `json:"company,omitempty"`
	BankAccount ShopOwnerBankAccount `json:"bankAccount"`
}

// ShopOwnerContact represents merchant communication coordinates.
type ShopOwnerContact struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	District string `json:"district"`
	City     string `json:"city"`
	State    string `json:"state,omitempty"`
	Country  string `json:"country"`
}

// ShopOwnerCompany specifies business legal registration details.
type ShopOwnerCompany struct {
	Name      string `json:"name"`
	TaxOffice string `json:"taxOffice"`
	TaxNumber string `json:"taxNumber"`
}

// ShopOwnerBankAccount specifies registered merchant settlement account.
type ShopOwnerBankAccount struct {
	AccountHolder string `json:"accountHolder"`
	IBAN          string `json:"iban"`
}

// ShopSettings encapsulates public storefront display and behavior configuration.
type ShopSettings struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	Title        string `json:"title"`
	Slogan       string `json:"slogan,omitempty"`
	Announcement string `json:"announcement,omitempty"`
	Confirmation string `json:"confirmation,omitempty"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Access       bool   `json:"access"`
	Cart         bool   `json:"cart"`
	MobileView   string `json:"mobileView"` // "singleColumn" or "doubleColumn"
	Filter       bool   `json:"filter"`
	OutOfStock   bool   `json:"outOfStock"`
	Language     string `json:"language"` // "TR" or "EN"
	Vacation     bool   `json:"vacation"`
}

// ShopSettingsUpdateRequest defines mutable parameters for PUT /shop/settings.
type ShopSettingsUpdateRequest struct {
	Title        *string `json:"title,omitempty"`
	Slogan       *string `json:"slogan,omitempty"`
	Announcement *string `json:"announcement,omitempty"`
	Confirmation *string `json:"confirmation,omitempty"`
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Access       *bool   `json:"access,omitempty"`
	Cart         *bool   `json:"cart,omitempty"`
	MobileView   *string `json:"mobileView,omitempty"`
	Filter       *bool   `json:"filter,omitempty"`
	OutOfStock   *bool   `json:"outOfStock,omitempty"`
	Language     *string `json:"language,omitempty"`
	Vacation     *bool   `json:"vacation,omitempty"`
}

// ShopService exposes merchant profile and store preference endpoints.
type ShopService struct {
	client *Client
}

// GetOwner retrieves legal entity and contact information of the merchant.
func (s *ShopService) GetOwner(ctx context.Context) (*ShopOwner, error) {
	var owner ShopOwner
	_, err := s.client.execute(ctx, http.MethodGet, "/shop/owner", nil, nil, &owner)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

// GetSettings retrieves current storefront display and configuration settings.
func (s *ShopService) GetSettings(ctx context.Context) (*ShopSettings, error) {
	var settings ShopSettings
	_, err := s.client.execute(ctx, http.MethodGet, "/shop/settings", nil, nil, &settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// UpdateSettings modifies storefront preferences and contact coordinates.
func (s *ShopService) UpdateSettings(ctx context.Context, req *ShopSettingsUpdateRequest) (*ShopSettings, error) {
	var updated ShopSettings
	_, err := s.client.execute(ctx, http.MethodPut, "/shop/settings", nil, req, &updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}
