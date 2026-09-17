# Authentication

Shopier supports two authentication models:
1. **Personal Access Tokens (PAT)**: Static tokens generated via the developer portal, ideal for scripts, backend batch jobs, and private single-merchant integrations.
2. **OAuth 2.0**: Standard authorization code flow for public multi-tenant applications where merchants grant consent to access their store resources.

---

## 1. Personal Access Tokens (PAT)

Pass the token string directly to `shopier.NewClient`:

```go
client, err := shopier.NewClient("pat_token_from_developer_portal")
if err != nil {
	log.Fatal(err)
}
```

---

## 2. OAuth 2.0 Flow

Shopier enforces OAuth 2.0 endpoints through specific ports:
- **Authorization Consent**: `https://developer.shopier.com/v1/oauth2/authorize`
- **Token Exchange & Refresh**: `https://api.shopier.com:8443/v1/oauth2/token`
- **Token Revocation**: `https://api.shopier.com:8443/v1/oauth2/revoke`

The `oauth` package manages these endpoint requirements automatically.

### Initializing OAuth Config

```go
import "github.com/AdisGroup/shopier-go/oauth"

cfg := oauth.NewConfig(
	"YOUR_CLIENT_ID",
	"YOUR_CLIENT_SECRET",
	"https://yourapp.com/oauth/callback",
)
```

### Step 1: Redirecting to Consent URL

Generate the consent page URL with requested permission scopes and a CSRF state token:

```go
consentURL := cfg.AuthCodeURL(
	"csrf_state_token",
	oauth.ScopeOrdersRead,
	oauth.ScopeOrdersWrite,
	oauth.ScopeProductsRead,
	oauth.ScopeProductsWrite,
)

// Redirect user to consentURL
```

### Step 2: Exchanging Authorization Code

In your callback HTTP handler, exchange the temporary `code` for access and refresh tokens:

```go
token, err := cfg.Exchange(ctx, code)
if err != nil {
	log.Fatalf("OAuth exchange failed: %v", err)
}

fmt.Printf("Access Token: %s (Expires in %d seconds)\n", token.AccessToken, token.ExpiresIn)
fmt.Printf("Refresh Token: %s\n", token.RefreshToken)
```

### Step 3: Refreshing Expired Tokens

Shopier access tokens expire in 3 days (259,200 seconds). Use the refresh token to request a new access token:

```go
if token.Expired() {
	newToken, err := cfg.RefreshToken(ctx, token.RefreshToken)
	if err != nil {
		log.Fatalf("Token refresh failed: %v", err)
	}
	token = newToken
}
```

### Step 4: Revoking Access

To invalidate a token when a merchant disconnects your application:

```go
err := cfg.Revoke(ctx, token.AccessToken)
if err != nil {
	log.Fatalf("Token revocation failed: %v", err)
}
```

---

## Permission Scopes Reference

| Scope Constant | Scope Name | Permitted Operations |
| :--- | :--- | :--- |
| `oauth.ScopeOrdersRead` | `orders:read` | List/Get orders, transactions, order webhooks |
| `oauth.ScopeOrdersWrite` | `orders:write` | Update order fulfillments and tracking numbers |
| `oauth.ScopeProductsRead` | `products:read` | Read products, categories, variations, selections |
| `oauth.ScopeProductsWrite` | `products:write` | Create, update, delete products and taxonomy |
| `oauth.ScopeShippingsRead` | `shippings:read` | Read shipping labels and tracking status |
| `oauth.ScopeShippingsWrite` | `shippings:write` | Generate contracted shipping codes |
| `oauth.ScopeDiscountsRead` | `discounts:read` | Read promo codes and automatic discounts |
| `oauth.ScopeDiscountsWrite` | `discounts:write` | Create and manage discount campaigns |
| `oauth.ScopePayoutsRead` | `payouts:read` | Read merchant payout ledgers and balance |
| `oauth.ScopeRefundsRead` | `refunds:read` | Read refund records and status |
| `oauth.ScopeRefundsWrite` | `refunds:write` | Initiate order refund requests |
| `oauth.ScopeShopRead` | `shop:read` | Read owner details and storefront settings |
| `oauth.ScopeShopWrite` | `shop:write` | Update storefront configuration |

---

## Data Models & Types

### `oauth.Token` {#oauthtoken-model}

Represents credentials issued by Shopier OAuth 2.0 authorization server.

| Field | Type | JSON Tag | Description |
| :--- | :--- | :--- | :--- |
| `AccessToken` | `string` | `access_token` | Bearer token used to authenticate API calls |
| `RefreshToken` | `string` | `refresh_token` | Long-lived token used to acquire new access tokens |
| `TokenType` | `string` | `token_type` | Token type (e.g. `"Bearer"`) |
| `ExpiresIn` | `int64` | `expires_in` | Validity period in seconds (typically `259200` = 3 days) |
| `Scope` | `string` | `scope` | Space-delimited string of granted permission scopes |

#### Helper Methods:
- `token.Expired() bool`: Returns `true` if current time is past expiry.

