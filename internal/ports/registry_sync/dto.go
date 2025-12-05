package registry_sync

// SyncRegistryRequest defines the request body for the /v1/internal/registry-sync API
type SyncRegistryRequest struct {
	RegistryEnv string   `json:"registry_env"`
	Domains     []string `json:"domains"`
}

// DomainSyncSummary provides a summary of the sync operation for a single domain
type DomainSyncSummary struct {
	Domain                 string `json:"domain"`
	NewSellers             int    `json:"new_sellers"`
	UpdatedSellers         int    `json:"updated_sellers"`
	DeactivatedSellers     int    `json:"deactivated_sellers"`
	TotalSellersInRegistry int    `json:"total_sellers_in_registry"`
}

// SyncRegistryResponse defines the response body for the /v1/internal/registry-sync API
type SyncRegistryResponse struct {
	RegistryEnv string              `json:"registry_env"`
	Domains     []DomainSyncSummary `json:"domains"`
	RunAt       string              `json:"run_at"`
}
