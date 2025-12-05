package ports

import "time"


// PendingCatalogSyncSellersResponse defines the response body for the pending catalog sync sellers API
type PendingCatalogSyncSellersResponse struct {
	Domain       string       `json:"domain"`
	RegistryEnv  string       `json:"registry_env"`
	StatusFilter []string     `json:"status_filter"`
	Sellers      []SellerInfo `json:"sellers"`
	Page         PageInfo     `json:"page"`
}

// SellerInfo defines the structure for a seller's catalog sync information
type SellerInfo struct {
	SellerID      string     `json:"seller_id"`
	Status        string     `json:"status"`
	LastPullAt    *time.Time `json:"last_pull_at"`
	LastSuccessAt *time.Time `json:"last_success_at"`
	LastError     *string    `json:"last_error"`
}

// PageInfo defines the structure for pagination information
type PageInfo struct {
	Limit    int    `json:"limit"`
	Page     int    `json:"page"`
	HasMore  bool   `json:"has_more"`
}
// CatalogSyncStatusResponse defines the response body for the catalog sync status API
type CatalogSyncStatusResponse struct {
	SellerID           string     `json:"seller_id"`
	Domain             string     `json:"domain"`
	RegistryEnv        string     `json:"registry_env"`
	Status             string     `json:"status"`
	LastPullAt         *time.Time `json:"last_pull_at"`
	LastSuccessAt      *time.Time `json:"last_success_at"`
	LastError          *string    `json:"last_error"`
	SyncVersion        int64      `json:"sync_version"`
	RegistryLastSeenAt time.Time  `json:"registry_last_seen_at"`
}
