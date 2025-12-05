package constants

const (
	// General Errors
	ErrInvalidRequestBody           = "Invalid request body"
	ErrUpdatesArrayEmpty            = "updates array cannot be empty"
	ErrRequiredPermissionsFields    = "bap_id, domain, registry_env, and seller_ids are required"
	ErrFailedToUpdatePermissions    = "Failed to update permissions"
	ErrFailedToQueryPermissions     = "Failed to query permissions"

	// Catalog Sync Errors
	ErrDomainRequired               = "domain query parameter is required"
	ErrSellerIDAndDomainRequired    = "seller_id path parameter and domain query parameter are required"
	ErrInvalidLimitParameter        = "Invalid limit parameter"
	ErrInvalidOffsetParameter       = "Invalid offset parameter"
	ErrGetPendingSellers            = "Failed to get pending catalog sync sellers"
	ErrGetSyncStatus                = "Failed to get sync status"
	ErrRecordNotFound               = "Record not found for the specified seller_id, domain, and registry_env"
	
	// Registry Sync Errors
	ErrFailedToStartRegistrySync    = "Failed to start registry sync"
)
