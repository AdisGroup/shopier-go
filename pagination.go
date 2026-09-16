package shopier

import (
	"net/http"
	"net/url"
	"strconv"
)

// PaginationInfo holds pagination metadata extracted from Shopier response headers.
// Shopier delivers pagination state via custom HTTP headers rather than response bodies:
// - Shopier-Pagination-Page
// - Shopier-Pagination-Limit
// - Shopier-Pagination-Total-Pages
// - Shopier-Pagination-Total-Items
type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

// HasNextPage returns true if more pages are available after the current page.
func (p PaginationInfo) HasNextPage() bool {
	return p.TotalPages > 0 && p.Page < p.TotalPages
}

// NextPage returns the index for the next page, or 0 if no subsequent pages exist.
func (p PaginationInfo) NextPage() int {
	if !p.HasNextPage() {
		return 0
	}
	return p.Page + 1
}

// ParsePaginationHeaders reads Shopier-Pagination-* headers from an HTTP response.
func ParsePaginationHeaders(header http.Header) PaginationInfo {
	var info PaginationInfo

	if v := header.Get("Shopier-Pagination-Page"); v != "" {
		info.Page, _ = strconv.Atoi(v)
	}
	if v := header.Get("Shopier-Pagination-Limit"); v != "" {
		info.Limit, _ = strconv.Atoi(v)
	}
	if v := header.Get("Shopier-Pagination-Total-Pages"); v != "" {
		info.TotalPages, _ = strconv.Atoi(v)
	}
	if v := header.Get("Shopier-Pagination-Total-Items"); v != "" {
		info.TotalItems, _ = strconv.Atoi(v)
	}

	return info
}

// ListOptions specifies common query parameters for all Shopier list endpoints.
type ListOptions struct {
	// Page specifies the page of results to return (default: 1, min: 1).
	Page int `url:"page,omitempty"`

	// Limit specifies the maximum number of items returned (default: 10, min: 1, max: 50).
	Limit int `url:"limit,omitempty"`

	// Sort specifies sort order or field where supported by the endpoint.
	Sort string `url:"sort,omitempty"`
}

// Values converts ListOptions into standard url.Values.
func (o *ListOptions) Values() url.Values {
	v := make(url.Values)
	if o == nil {
		return v
	}
	if o.Page > 0 {
		v.Set("page", strconv.Itoa(o.Page))
	}
	if o.Limit > 0 {
		v.Set("limit", strconv.Itoa(o.Limit))
	}
	if o.Sort != "" {
		v.Set("sort", o.Sort)
	}
	return v
}

// PageResponse wraps an item slice alongside its Shopier pagination metadata.
type PageResponse[T any] struct {
	Items      []T            `json:"items"`
	Pagination PaginationInfo `json:"pagination"`
}
