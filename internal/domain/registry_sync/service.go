package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"adapter/internal/config"
	registryPorts "adapter/internal/ports/registry_sync"
	catalogPorts "adapter/internal/ports/catalog_sync"
	"adapter/internal/shared/crypto"
	"adapter/internal/shared/log"
	"github.com/go-resty/resty/v2"
)

type ONDCLookupRequest struct {
	Country string `json:"country"`
	Type    string `json:"type"`
	Domain  string `json:"domain"`
}

type Subscriber struct {
	SubscriberID  string `json:"subscriber_id"`
	UkID          string `json:"ukId"`
	BrID          string `json:"br_id"`
	Domain        string `json:"domain"`
	Country       string `json:"country"`
	City          string `json:"city"`
	SigningKey    string `json:"signing_public_key"`
	EncryptionKey string `json:"encr_public_key"`
	Status        string `json:"status"`
	ValidFrom     string `json:"valid_from"`
	ValidUntil    string `json:"valid_until"`
	Created       string `json:"created"`
	Updated       string `json:"updated"`
}

type ONDCLookupResponse []Subscriber

type ONDCService struct {
	client       *resty.Client
	crypto       *crypto.ONDCCrypto
	sellerRepo   catalogPorts.SellerRepository
	domains      []string
	registryURL  string
	privateKey   string
	subscriberID string
	uniqueKeyID  string
	registryEnv  string
}

func NewONDCService(sellerRepo catalogPorts.SellerRepository, cfg *config.Config) *ONDCService {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(5 * time.Second)

	return &ONDCService{
		client:       client,
		crypto:       crypto.NewONDCCrypto(),
		sellerRepo:   sellerRepo,
		domains:      cfg.Domains,
		registryURL:  cfg.RegistryURL,
		privateKey:   cfg.PrivateKey,
		subscriberID: cfg.SubscriberID,
		uniqueKeyID:  cfg.UniqueKeyID,
		registryEnv:  cfg.RegistryEnv,
	}
}

func (s *ONDCService) SyncRegistry(req registryPorts.SyncRegistryRequest) (*registryPorts.SyncRegistryResponse, error) {
	runAt := time.Now()
	response := &registryPorts.SyncRegistryResponse{
		RegistryEnv: req.RegistryEnv,
		RunAt:       runAt.Format(time.RFC3339),
		Domains:     []registryPorts.DomainSyncSummary{},
	}

	for _, domain := range req.Domains {
		summary := registryPorts.DomainSyncSummary{Domain: domain}

		registrySellers, err := s.FetchSellersFromRegistry(domain)
		if err != nil {
			log.Error(context.Background(), err, fmt.Sprintf("Failed to fetch sellers from registry for domain %s", domain))
			continue
		}
		summary.TotalSellersInRegistry = len(registrySellers)

		dbSellers, err := s.sellerRepo.GetSellersByDomainAndRegistry(domain, req.RegistryEnv)
		if err != nil {
			log.Error(context.Background(), err, fmt.Sprintf("Failed to fetch sellers from DB for domain %s", domain))
			continue
		}

		registrySellerMap := make(map[string]catalogPorts.Seller)
		now := time.Now()
		for _, sub := range registrySellers {
			validFrom, _ := time.Parse(time.RFC3339, sub.ValidFrom)
			validUntil, _ := time.Parse(time.RFC3339, sub.ValidUntil)
			raw, _ := json.Marshal(sub)

			seller := catalogPorts.Seller{
				SellerID: sub.SubscriberID, Domain: sub.Domain, RegistryEnv: req.RegistryEnv,
				Status: sub.Status, Type: "BPP", SubscriberURL: sub.SubscriberID,
				Country: sub.Country, City: sub.City, ValidFrom: validFrom, ValidUntil: validUntil,
				Active: true, LastSeenInReg: now, RegistryRaw: string(raw),
			}
			registrySellerMap[seller.SellerID] = seller
		}

		dbSellerMap := make(map[string]catalogPorts.Seller)
		for _, seller := range dbSellers {
			dbSellerMap[seller.SellerID] = seller
		}

		var sellersToInsert []catalogPorts.Seller
		var sellersToUpdate []catalogPorts.Seller
		var removedSellerIDs []string

		for id, seller := range registrySellerMap {
			if _, exists := dbSellerMap[id]; !exists {
				sellersToInsert = append(sellersToInsert, seller)
			} else {
				sellersToUpdate = append(sellersToUpdate, seller)
			}
		}

		for id := range dbSellerMap {
			if _, exists := registrySellerMap[id]; !exists {
				removedSellerIDs = append(removedSellerIDs, id)
			}
		}

		summary.NewSellers = len(sellersToInsert)
		summary.UpdatedSellers = len(sellersToUpdate)
		summary.DeactivatedSellers = len(removedSellerIDs)

		if len(sellersToInsert) > 0 {
			if err := s.sellerRepo.InsertSellers(sellersToInsert); err != nil {
				log.Error(context.Background(), err, "Failed to insert new sellers")
			}
		}
		if len(sellersToUpdate) > 0 {
			if err := s.sellerRepo.UpdateSellers(sellersToUpdate); err != nil {
				log.Error(context.Background(), err, "Failed to update existing sellers")
			}
		}
		for _, seller := range sellersToInsert {
			state := &catalogPorts.SellerCatalogState{
				SellerID: seller.SellerID, Domain: domain, RegistryEnv: req.RegistryEnv,
				Status: catalogPorts.CatalogStatusNotSynced,
			}
			if err := s.sellerRepo.UpsertCatalogState(state); err != nil {
				log.Error(context.Background(), err, "Failed to insert catalog state")
			}
		}
		if len(removedSellerIDs) > 0 {
			if err := s.sellerRepo.DeactivateSellers(removedSellerIDs, domain, req.RegistryEnv); err != nil {
				log.Error(context.Background(), err, "Failed to deactivate sellers")
			}
		}
		response.Domains = append(response.Domains, summary)
	}
	return response, nil
}

func (s *ONDCService) FetchSellersFromRegistry(domain string) (ONDCLookupResponse, error) {
	reqBody := ONDCLookupRequest{Country: "IND", Type: "BPP", Domain: domain}
	authHeader, err := s.generateAuthHeader(reqBody)
	if err != nil {
		return nil, err
	}
	var response ONDCLookupResponse
	resp, err := s.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authHeader).
		SetBody(reqBody).
		SetResult(&response).
		Post(s.registryURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode(), resp.String())
	}
	return response, nil
}

func (s *ONDCService) generateAuthHeader(body ONDCLookupRequest) (string, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	if s.privateKey == "" {
		return "", fmt.Errorf("PRIVATE_KEY not configured")
	}
	currentTime := int(time.Now().Unix())
	ttl := 30
	signature, err := s.crypto.SignRequest(s.privateKey, payload, currentTime, ttl)
	if err != nil {
		return "", err
	}
	authHeader := fmt.Sprintf(
		`Signature keyId="%s|%s|ed25519",algorithm="ed25519",created="%d",expires="%d",headers="(created) (expires) digest",signature="%s"`,
		s.subscriberID, s.uniqueKeyID, currentTime, currentTime+ttl, signature,
	)
	return authHeader, nil
}
