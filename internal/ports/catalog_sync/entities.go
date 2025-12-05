package ports

import "time"

type Seller struct {
	SellerID      string    `json:"seller_id" gorm:"primaryKey;column:seller_id;type:text"`
	Domain        string    `json:"domain" gorm:"primaryKey;column:domain;type:text"`
	RegistryEnv   string    `json:"registry_env" gorm:"primaryKey;column:registry_env;type:text"`
	Status        string    `json:"status" gorm:"column:status;type:text"`
	Type          string    `json:"type" gorm:"column:type;type:text"`
	SubscriberURL string    `json:"subscriber_url" gorm:"column:subscriber_url;type:text"`
	Country       string    `json:"country" gorm:"column:country;type:text"`
	City          string    `json:"city" gorm:"column:city;type:text"`
	ValidFrom     time.Time `json:"valid_from" gorm:"column:valid_from;type:timestamptz"`
	ValidUntil    time.Time `json:"valid_until" gorm:"column:valid_until;type:timestamptz"`
	Active        bool      `json:"active" gorm:"column:active;type:boolean"`
	RegistryRaw   string    `json:"registry_raw" gorm:"column:registry_raw;type:jsonb"`
	LastSeenInReg time.Time `json:"last_seen_in_reg" gorm:"column:last_seen_in_reg;type:timestamptz"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Seller) TableName() string {
	return "sellers"
}

type CatalogStatus string

const (
	CatalogStatusNotSynced CatalogStatus = "NOT_SYNCED"
	CatalogStatusSyncing   CatalogStatus = "SYNCING"
	CatalogStatusSynced    CatalogStatus = "SYNCED"
	CatalogStatusFailed    CatalogStatus = "FAILED"
)

type SellerCatalogState struct {
	SellerID      string        `gorm:"primaryKey;column:seller_id;type:text"`
	Domain        string        `gorm:"primaryKey;column:domain;type:text"`
	RegistryEnv   string        `gorm:"primaryKey;column:registry_env;type:text"`
	Status        CatalogStatus `gorm:"column:status;type:text"`
	LastPullAt    *time.Time    `gorm:"column:last_pull_at;type:timestamptz"`
	LastSuccessAt *time.Time    `gorm:"column:last_success_at;type:timestamptz"`
	LastError     *string       `gorm:"column:last_error;type:text"`
	SyncVersion   int64         `gorm:"column:sync_version;type:bigint"`
	UpdatedAt     time.Time     `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`
}

func (SellerCatalogState) TableName() string {
	return "seller_catalog_state"
}
