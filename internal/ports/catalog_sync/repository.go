package ports

type SellerRepository interface {
	InsertSellers(sellers []Seller) error
	UpdateSellers(sellers []Seller) error
	GetAllSellers() ([]Seller, error)
	GetSellerByID(sellerID, domain, registryEnv string) (*Seller, error)
	GetPendingSellers(domain, registryEnv, status string, limit, offset int) ([]SellerInfo, error)
	GetSellersByDomainAndRegistry(domain, registryEnv string) ([]Seller, error)
	DeactivateSellers(sellerIDs []string, domain, registryEnv string) error
	UpsertCatalogState(state *SellerCatalogState) error
	GetSellerCatalogState(sellerID, domain, registryEnv string) (*SellerCatalogState, error)
}
