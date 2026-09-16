package oauth

// Shopier permission scopes governing access to merchant resources.
const (
	ScopeOrdersRead     = "orders:read"
	ScopeOrdersWrite    = "orders:write"
	ScopeProductsRead   = "products:read"
	ScopeProductsWrite  = "products:write"
	ScopeShippingsRead  = "shippings:read"
	ScopeShippingsWrite = "shippings:write"
	ScopeDiscountsRead  = "discounts:read"
	ScopeDiscountsWrite = "discounts:write"
	ScopePayoutsRead    = "payouts:read"
	ScopeRefundsRead    = "refunds:read"
	ScopeRefundsWrite   = "refunds:write"
	ScopeShopRead       = "shop:read"
	ScopeShopWrite      = "shop:write"
)

// AllScopes returns a slice containing every valid Shopier permission scope.
func AllScopes() []string {
	return []string{
		ScopeOrdersRead,
		ScopeOrdersWrite,
		ScopeProductsRead,
		ScopeProductsWrite,
		ScopeShippingsRead,
		ScopeShippingsWrite,
		ScopeDiscountsRead,
		ScopeDiscountsWrite,
		ScopePayoutsRead,
		ScopeRefundsRead,
		ScopeRefundsWrite,
		ScopeShopRead,
		ScopeShopWrite,
	}
}
