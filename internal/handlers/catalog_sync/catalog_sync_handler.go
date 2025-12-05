package handlers

import (
	catalogSyncPorts "adapter/internal/ports/catalog_sync"
	"adapter/internal/shared/constants"
	"adapter/internal/shared/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CatalogSyncHandler struct {
	service catalogSyncPorts.Service
}

func NewCatalogSyncHandler(service catalogSyncPorts.Service) *CatalogSyncHandler {
	return &CatalogSyncHandler{service: service}
}

func (h *CatalogSyncHandler) GetPendingCatalogSyncSellers(c *fiber.Ctx) error {
	domain := c.Query("domain")
	if domain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrDomainRequired,
		})
	}
	registryEnv := c.Query("registry_env", "preprod")
	status := c.Query("status")
	limit := c.QueryInt("limit", 100)
	page := c.QueryInt("page", 1)
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	response, err := h.service.GetPendingCatalogSyncSellers(domain, registryEnv, status, limit, page, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrGetPendingSellers,
		})
	}
	return c.Status(fiber.StatusOK).JSON(utils.ApiResponse{
		Success: true,
		Message: "Pending catalog sync sellers retrieved successfully",
		Data:    response,
	})
}

func (h *CatalogSyncHandler) GetSyncStatus(c *fiber.Ctx) error {
	sellerID := c.Params("seller_id")
	domain := c.Query("domain")
	registryEnv := c.Query("registry_env", "preprod") // Default to "preprod"

	if sellerID == "" || domain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrSellerIDAndDomainRequired,
		})
	}

	response, err := h.service.GetSyncStatus(sellerID, domain, registryEnv)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(utils.ApiResponse{
				Success: false,
				Message: constants.ErrRecordNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrGetSyncStatus,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ApiResponse{
		Success: true,
		Message: "Sync status retrieved successfully",
		Data:    response,
	})
}