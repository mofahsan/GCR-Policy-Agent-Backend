package ports

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) UpsertBaps(baps map[string]Bap) error {
	var bapList []Bap
	for _, b := range baps {
		bapList = append(bapList, b)
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "bap_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_seen_at"}),
	}).Create(&bapList).Error
}

func (r *GormRepository) UpsertBapAccessPolicies(policies []BapAccessPolicy) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "seller_id"}, {Name: "domain"}, {Name: "registry_env"}, {Name: "bap_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"decision", "decision_source", "decided_at", "expires_at", "reason", "updated_at"}),
	}).Create(&policies).Error
}
func (r *GormRepository) FindBapByID(bapID string) (*Bap, error) {
	var bap Bap
	if err := r.db.First(&bap, "bap_id = ?", bapID).Error; err != nil {
		return nil, err
	}
	return &bap, nil
}

func (r *GormRepository) QueryBapAccessPolicies(bapID, domain, registryEnv string, sellerIDs []string) ([]BapAccessPolicy, error) {
	var policies []BapAccessPolicy
	if err := r.db.Where("bap_id = ? AND domain = ? AND registry_env = ? AND seller_id IN ?", bapID, domain, registryEnv, sellerIDs).Find(&policies).Error; err != nil {
		return nil, err
	}
	return policies, nil
}
