package ports

type PermissionsRepository interface {
	UpsertBaps(baps map[string]Bap) error
	UpsertBapAccessPolicies(policies []BapAccessPolicy) error
	FindBapByID(bapID string) (*Bap, error)
	QueryBapAccessPolicies(bapID, domain, registryEnv string, sellerIDs []string) ([]BapAccessPolicy, error)
}
