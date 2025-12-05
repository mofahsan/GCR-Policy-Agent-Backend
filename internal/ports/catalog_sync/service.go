package ports

import (
	"adapter/internal/shared/log"
	"adapter/internal/shared/utils"

	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service interface {
	GetPendingCatalogSyncSellers(domain, registryEnv, status string, limit, page, offset int) (*PendingCatalogSyncSellersResponse, error)
	GetSyncStatus(sellerID, domain, registryEnv string) (*CatalogSyncStatusResponse, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) InsertSellers(sellers []Seller) error {
	log.Info(context.Background(), fmt.Sprintf("Attempting to insert %d new sellers...", len(sellers)))
	return r.db.Create(&sellers).Error
}

func (r *GormRepository) UpdateSellers(sellers []Seller) error {
	log.Info(context.Background(), fmt.Sprintf("Attempting to update %d existing sellers...", len(sellers)))
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, seller := range sellers {
			if err := tx.Model(&Seller{}).Where("seller_id = ? AND domain = ? AND registry_env = ?", seller.SellerID, seller.Domain, seller.RegistryEnv).Updates(seller).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormRepository) GetAllSellers() ([]Seller, error) {
	var sellers []Seller
	if err := r.db.Find(&sellers).Error; err != nil {
		return nil, err
	}
	return sellers, nil
}

func (r *GormRepository) GetSellerByID(sellerID, domain, registryEnv string) (*Seller, error) {
	var seller Seller
	if err := r.db.Where("seller_id = ? AND domain = ? AND registry_env = ?", sellerID, domain, registryEnv).First(&seller).Error; err != nil {
		return nil, err
	}
	return &seller, nil
}

func (r *GormRepository) GetSellersByDomainAndRegistry(domain, registryEnv string) ([]Seller, error) {
	var sellers []Seller
	if err := r.db.Where("domain = ? AND registry_env = ? AND active = ?", domain, registryEnv, true).Find(&sellers).Error; err != nil {
		return nil, err
	}
	return sellers, nil
}

func (r *GormRepository) DeactivateSellers(sellerIDs []string, domain, registryEnv string) error {
	return r.db.Model(&Seller{}).Where("seller_id IN ? AND domain = ? AND registry_env = ?", sellerIDs, domain, registryEnv).Update("active", false).Error
}

func (r *GormRepository) UpsertCatalogState(state *SellerCatalogState) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "seller_id"}, {Name: "domain"}, {Name: "registry_env"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "last_pull_at", "last_success_at", "last_error", "sync_version", "updated_at"}),
	}).Create(state).Error
}

func (r *GormRepository) GetPendingSellers(domain, registryEnv, status string, limit, offset int) ([]SellerInfo, error) {
	var sellers []SellerInfo
	query := r.db.Table("sellers as s").
		Select("s.seller_id, scs.status, scs.last_pull_at, scs.last_success_at, scs.last_error").
		Joins("LEFT JOIN seller_catalog_state scs ON s.seller_id = scs.seller_id AND s.domain = scs.domain AND s.registry_env = scs.registry_env").
		Where("s.domain = ? AND s.registry_env = ? AND s.active = ?", domain, registryEnv, true)

	var statusConditions []string
	var statusValues []interface{}

	inputStatuses := utils.SplitAndTrim(status)

	if len(inputStatuses) == 0 {
		// Default to NOT_SYNCED and FAILED if no status is provided
		statusConditions = append(statusConditions, "scs.status = ? OR scs.status IS NULL")
		statusValues = append(statusValues, "NOT_SYNCED")
		statusConditions = append(statusConditions, "scs.status = ?")
		statusValues = append(statusValues, "FAILED")
	} else {
		// Filter by provided statuses
		for _, s := range inputStatuses {
			if s == "NOT_SYNCED" {
				// Handle NOT_SYNCED which can be explicit or NULL
				statusConditions = append(statusConditions, "scs.status = ? OR scs.status IS NULL")
				statusValues = append(statusValues, "NOT_SYNCED")
			} else if s == "FAILED" {
				// Handle FAILED explicitly
				statusConditions = append(statusConditions, "scs.status = ?")
				statusValues = append(statusValues, "FAILED")
			}
		}
	}

	if len(statusConditions) > 0 {
		// Combine all status conditions with OR
		query = query.Where(r.db.Where(utils.JoinConditions(statusConditions, " OR "), statusValues...))
	}

	err := query.Order("s.seller_id").
		Limit(limit + 1).
		Offset(offset).
		Scan(&sellers).Error

	return sellers, err
}

func (r *GormRepository) GetSellerCatalogState(sellerID, domain, registryEnv string) (*SellerCatalogState, error) {
	var state SellerCatalogState
	if err := r.db.Where("seller_id = ? AND domain = ? AND registry_env = ?", sellerID, domain, registryEnv).First(&state).Error; err != nil {
		return nil, err
	}
	return &state, nil
}
