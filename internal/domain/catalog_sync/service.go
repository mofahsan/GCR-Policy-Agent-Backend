package catalogPorts

import (
	catalogPorts "adapter/internal/ports/catalog_sync"
	"gorm.io/gorm"
	"strings"
)

type CatalogSyncService struct {
	repo catalogPorts.SellerRepository
}

func NewCatalogSyncService(repo catalogPorts.SellerRepository) *CatalogSyncService {
	return &CatalogSyncService{repo: repo}
}


func (s *CatalogSyncService) GetPendingCatalogSyncSellers(domain, registryEnv, status string, limit, page, offset int) (*catalogPorts.PendingCatalogSyncSellersResponse, error) {
	sellers, err := s.repo.GetPendingSellers(domain, registryEnv, status, limit, offset)
	if err != nil {
		return nil, err
	}

	hasMore := len(sellers) > limit
	if hasMore {
		sellers = sellers[:limit] // Trim the extra record fetched for hasMore check
	}

	var statusFilter []string
	if status == "" {
		statusFilter = []string{"NOT_SYNCED", "FAILED"}
	} else {
		statusFilter = strings.Split(status, ",")
	}

	return &catalogPorts.PendingCatalogSyncSellersResponse{
		Domain:       domain,
		RegistryEnv:  registryEnv,
		StatusFilter: statusFilter,
		Sellers:      sellers,
		Page: catalogPorts.PageInfo{
			Limit:    limit,
			Page:     page,
			HasMore:  hasMore,
		},
	}, nil
}
func (s *CatalogSyncService) GetSyncStatus(sellerID, domain, registryEnv string) (*catalogPorts.CatalogSyncStatusResponse, error) {
	seller, err := s.repo.GetSellerByID(sellerID, domain, registryEnv)
	if err != nil {
		return nil, err // Could be gorm.ErrRecordNotFound
	}

	state, err := s.repo.GetSellerCatalogState(sellerID, domain, registryEnv)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// If state is not found, we still return a response but with default/empty values for state fields
	if err == gorm.ErrRecordNotFound {
		return &catalogPorts.CatalogSyncStatusResponse{
			SellerID:           seller.SellerID,
			Domain:             seller.Domain,
			RegistryEnv:        seller.RegistryEnv,
			Status:             string(catalogPorts.CatalogStatusNotSynced), // Default status
			RegistryLastSeenAt: seller.LastSeenInReg,
		}, nil
	}

	return &catalogPorts.CatalogSyncStatusResponse{
		SellerID:           state.SellerID,
		Domain:             state.Domain,
		RegistryEnv:        state.RegistryEnv,
		Status:             string(state.Status),
		LastPullAt:         state.LastPullAt,
		LastSuccessAt:      state.LastSuccessAt,
		LastError:          state.LastError,
		SyncVersion:        state.SyncVersion,
		RegistryLastSeenAt: seller.LastSeenInReg,
	}, nil
}
