package ports

import "time"

// PermissionsUpdateRequest defines the structure for a single permission update
type PermissionsUpdateRequest struct {
	SellerID       string     `json:"seller_id"`
	Domain         string     `json:"domain"`
	RegistryEnv    string     `json:"registry_env"`
	BapID          string     `json:"bap_id"`
	Decision       string     `json:"decision"`
	DecisionSource string     `json:"decision_source"`
	Reason         *string    `json:"reason"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

// PermissionsUpdateResponse defines the structure for a single permission update result
type PermissionsUpdateResponse struct {
	SellerID    string `json:"seller_id"`
	Domain      string `json:"domain"`
	RegistryEnv string `json:"registry_env"`
	BapID       string `json:"bap_id"`
	Decision    string `json:"decision"`
	Stored      bool   `json:"stored"`
}

// PermissionsQueryRequest defines the request body for the /v1/permissions/query API
type PermissionsQueryRequest struct {
	BapID           string   `json:"bap_id"`
	Domain          string   `json:"domain"`
	RegistryEnv     string   `json:"registry_env"`
	SellerIDs       []string `json:"seller_ids"`
	IncludeNoPolicy bool     `json:"include_no_policy"`
}

// PermissionDetail provides detailed permission information for a single seller
type PermissionDetail struct {
	SellerID       string     `json:"seller_id"`
	Domain         string     `json:"domain"`
	RegistryEnv    string     `json:"registry_env"`
	BapID          string     `json:"bap_id"`
	Decision       string     `json:"decision"`
	DecisionSource *string    `json:"decision_source,omitempty"`
	DecidedAt      *time.Time `json:"decided_at,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

// PermissionsQueryResponse defines the response body for the /v1/permissions/query API
type PermissionsQueryResponse struct {
	BapStatus   string             `json:"bap_status"`
	Domain      string             `json:"domain"`
	RegistryEnv string             `json:"registry_env"`
	Permissions []PermissionDetail `json:"permissions"`
}
