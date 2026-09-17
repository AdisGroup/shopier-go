package shopier

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	metaCSRFPattern1  = regexp.MustCompile(`(?i)<meta\s+[^>]*name=["']csrf-token["'][^>]*content=["']([^"']+)["']`)
	metaCSRFPattern2  = regexp.MustCompile(`(?i)<meta\s+[^>]*content=["']([^"']+)["'][^>]*name=["']csrf-token["']`)
	productURLPattern = regexp.MustCompile(`(?:https?://)?(?:www\.)?shopier\.com/([^/?#]+)/([0-9a-zA-Z_-]+)`)
)

// QuickCheckoutStage represents a step in the unofficial checkout process.
type QuickCheckoutStage string

const (
	StageParseInput        QuickCheckoutStage = "parse_input"
	StageFetchProduct      QuickCheckoutStage = "fetch_product"
	StageExtractCSRF       QuickCheckoutStage = "extract_csrf"
	StageCheckPayment      QuickCheckoutStage = "check_payment"
	StageFetchShippingForm QuickCheckoutStage = "fetch_shipping_form"
	StageProcessShipment   QuickCheckoutStage = "process_shipment"
)

// QuickCheckoutError captures failures during unofficial checkout link creation.
type QuickCheckoutError struct {
	Stage      QuickCheckoutStage `json:"stage"`
	StatusCode int                `json:"status_code,omitempty"`
	Message    string             `json:"message"`
	Body       string             `json:"body,omitempty"`
	Err        error              `json:"-"`
}

func (e *QuickCheckoutError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("quick_checkout [%s]: HTTP %d - %s", e.Stage, e.StatusCode, e.Message)
	}
	if e.Err != nil {
		return fmt.Sprintf("quick_checkout [%s]: %s: %v", e.Stage, e.Message, e.Err)
	}
	return fmt.Sprintf("quick_checkout [%s]: %s", e.Stage, e.Message)
}

func (e *QuickCheckoutError) Unwrap() error {
	return e.Err
}

// QuickCheckoutBuyer holds buyer and delivery details for the checkout session.
type QuickCheckoutBuyer struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`                // e.g. "+90 532 477 02 10" or "5324770210"
	PhoneCode string `json:"phone_code,omitempty"` // defaults to "TR"
	Country   string `json:"country,omitempty"`    // defaults to "Türkiye"
	City      string `json:"city,omitempty"`
	State     string `json:"state,omitempty"`
	Address   string `json:"address,omitempty"`
	ZipCode   string `json:"zip_code,omitempty"`
	TCIDNo    string `json:"tcid_no,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// QuickCheckoutRequest defines inputs to generate a checkout link.
type QuickCheckoutRequest struct {
	// ProductURL can be provided instead of Account and ProductID (e.g. "https://www.shopier.com/adisgroup/50900391").
	ProductURL string `json:"product_url,omitempty"`

	// Account is the seller username / shop slug (e.g. "adisgroup").
	Account string `json:"account,omitempty"`

	// ProductID is the unique numeric or alphanumeric product identifier.
	ProductID string `json:"product_id,omitempty"`

	// Quantity defaults to 1 if omitted.
	Quantity int `json:"quantity,omitempty"`

	// Options can contain optional variant or selection parameters.
	Options map[string]string `json:"options,omitempty"`

	// Buyer contains shipping and contact information.
	Buyer QuickCheckoutBuyer `json:"buyer"`
}

// QuickCheckoutResult contains the generated order ID and direct checkout link.
type QuickCheckoutResult struct {
	OrderID     string         `json:"order_id"`
	PaymentURL  string         `json:"payment_url"`
	Account     string         `json:"account"`
	ProductID   string         `json:"product_id"`
	RawResponse map[string]any `json:"raw_response,omitempty"`
}

// QuickCheckoutOption configures the QuickCheckoutService.
type QuickCheckoutOption func(*QuickCheckoutService)

// WithQuickCheckoutBaseURL sets a custom base URL (default: "https://www.shopier.com").
func WithQuickCheckoutBaseURL(baseURL string) QuickCheckoutOption {
	return func(s *QuickCheckoutService) {
		s.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithQuickCheckoutUserAgent sets a custom User-Agent header.
func WithQuickCheckoutUserAgent(ua string) QuickCheckoutOption {
	return func(s *QuickCheckoutService) {
		s.userAgent = ua
	}
}

// WithQuickCheckoutHTTPClient provides a custom base HTTP client (transport/proxy settings).
func WithQuickCheckoutHTTPClient(client *http.Client) QuickCheckoutOption {
	return func(s *QuickCheckoutService) {
		s.baseClient = client
	}
}

// WithQuickCheckoutTimeout sets request timeout.
func WithQuickCheckoutTimeout(timeout time.Duration) QuickCheckoutOption {
	return func(s *QuickCheckoutService) {
		s.timeout = timeout
	}
}

// QuickCheckoutService generates direct checkout and payment links.
type QuickCheckoutService struct {
	baseURL    string
	userAgent  string
	baseClient *http.Client
	timeout    time.Duration
}

// NewQuickCheckoutService creates an isolated QuickCheckoutService instance.
func NewQuickCheckoutService(opts ...QuickCheckoutOption) *QuickCheckoutService {
	s := &QuickCheckoutService{
		baseURL:   "https://www.shopier.com",
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		timeout:   30 * time.Second,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Create executes the 4-step frontend flow and returns a direct payment URL.
// Each execution uses a fresh, isolated cookie jar for total concurrency safety.
func (s *QuickCheckoutService) Create(ctx context.Context, req *QuickCheckoutRequest) (*QuickCheckoutResult, error) {
	if req == nil {
		return nil, &QuickCheckoutError{
			Stage:   StageParseInput,
			Message: "request payload cannot be nil",
		}
	}

	account := strings.TrimSpace(req.Account)
	productID := strings.TrimSpace(req.ProductID)

	if req.ProductURL != "" {
		parsedAccount, parsedID, err := parseProductURL(req.ProductURL)
		if err == nil {
			if account == "" {
				account = parsedAccount
			}
			if productID == "" {
				productID = parsedID
			}
		}
	}

	if account == "" || productID == "" {
		return nil, &QuickCheckoutError{
			Stage:   StageParseInput,
			Message: "both account and product_id are required (or provide a valid product_url)",
		}
	}

	if req.Buyer.Email == "" || req.Buyer.FirstName == "" || req.Buyer.LastName == "" || req.Buyer.Phone == "" {
		return nil, &QuickCheckoutError{
			Stage:   StageParseInput,
			Message: "buyer email, first_name, last_name, and phone are required",
		}
	}

	qty := req.Quantity
	if qty <= 0 {
		qty = 1
	}

	phoneCode := req.Buyer.PhoneCode
	if phoneCode == "" {
		phoneCode = "TR"
	}

	country := req.Buyer.Country
	if country == "" {
		country = "Türkiye"
	}

	// Create a per-request isolated HTTP client with fresh cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:   StageParseInput,
			Message: "failed to create isolated cookie jar",
			Err:     err,
		}
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: s.timeout,
	}
	if s.baseClient != nil && s.baseClient.Transport != nil {
		client.Transport = s.baseClient.Transport
	}

	productURL := fmt.Sprintf("%s/%s/%s", s.baseURL, account, productID)
	checkPaymentURL := fmt.Sprintf("%s/s/api/v1/check_payment_progress/%s", s.baseURL, account)
	shippingURL := fmt.Sprintf("%s/s/shipping/%s", s.baseURL, account)
	shipmentFormURL := fmt.Sprintf("%s/s/api/v1/shipment_form_process/%s", s.baseURL, account)

	// Step 1: GET Product Page to retrieve initial CSRF token and session cookies
	getReq, err := http.NewRequestWithContext(ctx, http.MethodGet, productURL, nil)
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:   StageFetchProduct,
			Message: "failed to build product page request",
			Err:     err,
		}
	}
	getReq.Header.Set("User-Agent", s.userAgent)
	getReq.Header.Set("Accept-Language", "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7")

	resp, err := client.Do(getReq)
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:   StageFetchProduct,
			Message: "failed to access product page",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:      StageFetchProduct,
			StatusCode: resp.StatusCode,
			Message:    "failed to read product page response body",
			Err:        err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &QuickCheckoutError{
			Stage:      StageFetchProduct,
			StatusCode: resp.StatusCode,
			Message:    "product page returned non-200 status",
			Body:       string(bodyBytes),
		}
	}

	csrfToken := extractCSRFToken(bodyBytes)
	if csrfToken == "" {
		return nil, &QuickCheckoutError{
			Stage:      StageExtractCSRF,
			StatusCode: resp.StatusCode,
			Message:    "csrf-token meta tag not found in product page HTML",
			Body:       string(bodyBytes),
		}
	}

	// Step 2: POST check_payment_progress
	checkPaymentValues := url.Values{}
	checkPaymentValues.Set("product_id", productID)
	checkPaymentValues.Set("quantity", strconv.Itoa(qty))
	for k, v := range req.Options {
		checkPaymentValues.Set(k, v)
	}

	checkReq, err := http.NewRequestWithContext(ctx, http.MethodPost, checkPaymentURL, strings.NewReader(checkPaymentValues.Encode()))
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:   StageCheckPayment,
			Message: "failed to build check_payment_progress request",
			Err:     err,
		}
	}
	checkReq.Header.Set("User-Agent", s.userAgent)
	checkReq.Header.Set("Accept", "*/*")
	checkReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	checkReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	checkReq.Header.Set("X-CSRF-Token", csrfToken)
	checkReq.Header.Set("Referer", productURL)

	checkResp, err := client.Do(checkReq)
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:   StageCheckPayment,
			Message: "check_payment_progress request failed",
			Err:     err,
		}
	}
	defer checkResp.Body.Close()

	checkBody, _ := io.ReadAll(checkResp.Body)
	if checkResp.StatusCode != http.StatusOK && checkResp.StatusCode != http.StatusFound {
		return nil, &QuickCheckoutError{
			Stage:      StageCheckPayment,
			StatusCode: checkResp.StatusCode,
			Message:    "check_payment_progress returned unexpected status",
			Body:       string(checkBody),
		}
	}

	// Step 3: POST shipping page to retrieve shipping form and form-level CSRF token
	shippingReq, err := http.NewRequestWithContext(ctx, http.MethodPost, shippingURL, strings.NewReader(checkPaymentValues.Encode()))
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:   StageFetchShippingForm,
			Message: "failed to build shipping page request",
			Err:     err,
		}
	}
	shippingReq.Header.Set("User-Agent", s.userAgent)
	shippingReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	shippingReq.Header.Set("Referer", productURL)

	shippingResp, err := client.Do(shippingReq)
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:   StageFetchShippingForm,
			Message: "failed to access shipping form page",
			Err:     err,
		}
	}
	defer shippingResp.Body.Close()

	shippingBody, err := io.ReadAll(shippingResp.Body)
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:      StageFetchShippingForm,
			StatusCode: shippingResp.StatusCode,
			Message:    "failed to read shipping form body",
			Err:        err,
		}
	}

	if shippingResp.StatusCode != http.StatusOK {
		return nil, &QuickCheckoutError{
			Stage:      StageFetchShippingForm,
			StatusCode: shippingResp.StatusCode,
			Message:    "shipping form page returned non-200 status",
			Body:       string(shippingBody),
		}
	}

	formCSRFToken := extractCSRFToken(shippingBody)
	if formCSRFToken == "" {
		formCSRFToken = csrfToken
	}

	// Step 4: POST shipment_form_process
	shipmentValues := url.Values{}
	shipmentValues.Set("Email", req.Buyer.Email)
	shipmentValues.Set("phone-contact-select", phoneCode)
	cleanPhone := cleanPhoneNumber(req.Buyer.Phone)
	shipmentValues.Set("formControlPhone", cleanPhone)
	shipmentValues.Set("Phone", req.Buyer.Phone)
	shipmentValues.Set("FirstName", req.Buyer.FirstName)
	shipmentValues.Set("LastName", req.Buyer.LastName)
	shipmentValues.Set("country", country)
	if req.Buyer.City != "" {
		shipmentValues.Set("City", req.Buyer.City)
	}
	if req.Buyer.State != "" {
		shipmentValues.Set("State", req.Buyer.State)
	}
	if req.Buyer.Address != "" {
		shipmentValues.Set("Address", req.Buyer.Address)
	}
	if req.Buyer.ZipCode != "" {
		shipmentValues.Set("ZipCode", req.Buyer.ZipCode)
	}
	if req.Buyer.TCIDNo != "" {
		shipmentValues.Set("TCIDNo", req.Buyer.TCIDNo)
	}
	if req.Buyer.Comment != "" {
		shipmentValues.Set("Comment", req.Buyer.Comment)
	}

	shipmentReq, err := http.NewRequestWithContext(ctx, http.MethodPost, shipmentFormURL, strings.NewReader(shipmentValues.Encode()))
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:   StageProcessShipment,
			Message: "failed to build shipment_form_process request",
			Err:     err,
		}
	}
	shipmentReq.Header.Set("User-Agent", s.userAgent)
	shipmentReq.Header.Set("Accept", "*/*")
	shipmentReq.Header.Set("Accept-Language", "tr-TR,tr;q=0.9")
	shipmentReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	shipmentReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	shipmentReq.Header.Set("X-CSRF-Token", formCSRFToken)
	shipmentReq.Header.Set("Referer", shippingURL)

	shipmentResp, err := client.Do(shipmentReq)
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:   StageProcessShipment,
			Message: "shipment_form_process request failed",
			Err:     err,
		}
	}
	defer shipmentResp.Body.Close()

	shipmentBody, err := io.ReadAll(shipmentResp.Body)
	if err != nil {
		return nil, &QuickCheckoutError{
			Stage:      StageProcessShipment,
			StatusCode: shipmentResp.StatusCode,
			Message:    "failed to read shipment_form_process response body",
			Err:        err,
		}
	}

	if shipmentResp.StatusCode != http.StatusOK {
		return nil, &QuickCheckoutError{
			Stage:      StageProcessShipment,
			StatusCode: shipmentResp.StatusCode,
			Message:    "shipment_form_process returned non-200 status",
			Body:       string(shipmentBody),
		}
	}

	var jsonResult map[string]any
	if err := json.Unmarshal(shipmentBody, &jsonResult); err != nil {
		return nil, &QuickCheckoutError{
			Stage:      StageProcessShipment,
			StatusCode: shipmentResp.StatusCode,
			Message:    "failed to parse shipment_form_process JSON response",
			Body:       string(shipmentBody),
			Err:        err,
		}
	}

	orderID := ""
	if rawID, ok := jsonResult["order_id"]; ok {
		orderID = fmt.Sprintf("%v", rawID)
	}

	if orderID == "" {
		return nil, &QuickCheckoutError{
			Stage:      StageProcessShipment,
			StatusCode: shipmentResp.StatusCode,
			Message:    "order_id not found in shipment_form_process response",
			Body:       string(shipmentBody),
		}
	}

	paymentURL := fmt.Sprintf("%s/s/payment/%s/%s", s.baseURL, account, orderID)

	return &QuickCheckoutResult{
		OrderID:     orderID,
		PaymentURL:  paymentURL,
		Account:     account,
		ProductID:   productID,
		RawResponse: jsonResult,
	}, nil
}

func extractCSRFToken(body []byte) string {
	if m := metaCSRFPattern1.FindSubmatch(body); len(m) > 1 {
		return strings.TrimSpace(string(m[1]))
	}
	if m := metaCSRFPattern2.FindSubmatch(body); len(m) > 1 {
		return strings.TrimSpace(string(m[1]))
	}
	return ""
}

func parseProductURL(rawURL string) (account, productID string, err error) {
	raw := strings.TrimSpace(rawURL)
	if raw == "" {
		return "", "", errors.New("product url cannot be empty")
	}

	// Try standard url.Parse if it contains scheme
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		u, parseErr := url.Parse(raw)
		if parseErr == nil && u.Path != "" {
			parts := strings.Split(strings.Trim(u.Path, "/"), "/")
			if len(parts) >= 2 {
				return parts[len(parts)-2], parts[len(parts)-1], nil
			}
		}
	}

	// Clean path segments
	clean := strings.Trim(raw, "/")
	// If domain without scheme was given (e.g. shopier.com/acc/id or www.shopier.com/acc/id)
	if idx := strings.Index(clean, "/"); idx != -1 {
		prefix := clean[:idx]
		if strings.Contains(prefix, "shopier.com") || strings.Contains(prefix, "localhost") || strings.Contains(prefix, "127.0.0.1") {
			clean = clean[idx+1:]
		}
	}

	parts := strings.Split(strings.Trim(clean, "/"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2], parts[len(parts)-1], nil
	}

	return "", "", errors.New("invalid product url format")
}

func cleanPhoneNumber(phone string) string {
	var buf bytes.Buffer
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			buf.WriteRune(r)
		}
	}
	digits := buf.String()
	if strings.HasPrefix(digits, "90") && len(digits) > 10 {
		return digits[2:]
	}
	return phone
}
