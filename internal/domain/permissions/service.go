package permissions

import (
	ports "adapter/internal/ports/permissions"

	"gorm.io/gorm"
	"time"
)

type PermissionsService struct {
	repo ports.PermissionsRepository
}

func NewPermissionsService(repo ports.PermissionsRepository) *PermissionsService {
	return &PermissionsService{repo: repo}
}

func (s *PermissionsService) UpdatePermissions(updates []ports.PermissionsUpdateRequest) ([]ports.PermissionsUpdateResponse, error) {
	var results []ports.PermissionsUpdateResponse
	var policiesToUpsert []ports.BapAccessPolicy
	bapsToUpsert := make(map[string]ports.Bap)

	for _, update := range updates {
		// Prepare BapAccessPolicy for upsert
		policy := ports.BapAccessPolicy{
			SellerID:       update.SellerID,
			Domain:         update.Domain,
			RegistryEnv:    update.RegistryEnv,
			BapID:          update.BapID,
			Decision:       ports.AccessDecision(update.Decision),
			DecisionSource: ports.DecisionSource(update.DecisionSource),
			DecidedAt:      time.Now(),
			ExpiresAt:      update.ExpiresAt,
			Reason:         update.Reason,
		}
		policiesToUpsert = append(policiesToUpsert, policy)

		// Collect unique BAPs to ensure they exist in the `baps` table
		if _, exists := bapsToUpsert[update.BapID]; !exists {
			bapsToUpsert[update.BapID] = ports.Bap{BapID: update.BapID}
		}

		results = append(results, ports.PermissionsUpdateResponse{
			SellerID:    update.SellerID,
			Domain:      update.Domain,
			RegistryEnv: update.RegistryEnv,
			BapID:       update.BapID,
			Decision:    update.Decision,
			Stored:      false, // Will be set to true after successful DB operation
		})
	}

	// Upsert BAPs first to satisfy foreign key constraints
	if err := s.repo.UpsertBaps(bapsToUpsert); err != nil {
		return results, err
	}

	// Upsert the access policies
	if err := s.repo.UpsertBapAccessPolicies(policiesToUpsert); err != nil {
		return results, err
	}

	// Mark all as stored if we reach here without an error
	for i := range results {
		results[i].Stored = true
	}

	return results, nil
}

func (s *PermissionsService) QueryPermissions(req ports.PermissionsQueryRequest) (*ports.PermissionsQueryResponse, error) {
	bapStatus := ""
	bap, err := s.repo.FindBapByID(req.BapID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			bapStatus = "NEW_BAP"
			// Create the BAP
			bapsToUpsert := map[string]ports.Bap{req.BapID: {BapID: req.BapID}}
			if err := s.repo.UpsertBaps(bapsToUpsert); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		bapStatus = "EXISTING_BAP"
		// Update last_seen_at
		bap.LastSeenAt = time.Now()
		bapsToUpsert := map[string]ports.Bap{req.BapID: *bap}
		if err := s.repo.UpsertBaps(bapsToUpsert); err != nil {
			return nil, err
		}
	}

	policies, err := s.repo.QueryBapAccessPolicies(req.BapID, req.Domain, req.RegistryEnv, req.SellerIDs)
	if err != nil {
		return nil, err
	}

	policyMap := make(map[string]ports.BapAccessPolicy)
	for _, p := range policies {
		policyMap[p.SellerID] = p
	}

	var permissions []ports.PermissionDetail
	for _, sellerID := range req.SellerIDs {
		if policy, ok := policyMap[sellerID]; ok {
			permissions = append(permissions, ports.PermissionDetail{
				SellerID:       policy.SellerID,
				Domain:         policy.Domain,
				RegistryEnv:    policy.RegistryEnv,
				BapID:          policy.BapID,
				Decision:       string(policy.Decision),
				DecisionSource: (*string)(&policy.DecisionSource),
				DecidedAt:      &policy.DecidedAt,
				ExpiresAt:      policy.ExpiresAt,
			})
		} else if req.IncludeNoPolicy {
			permissions = append(permissions, ports.PermissionDetail{
				SellerID:    sellerID,
				Domain:      req.Domain,
				RegistryEnv: req.RegistryEnv,
				BapID:       req.BapID,
				Decision:    "NO_POLICY",
			})
		}
	}

	return &ports.PermissionsQueryResponse{
		BapStatus:   bapStatus,
		Domain:      req.Domain,
		RegistryEnv: req.RegistryEnv,
		Permissions: permissions,
	}, nil
}
