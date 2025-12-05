package handlers

import (
	ondc "adapter/internal/domain/registry_sync"
	ports "adapter/internal/ports/registry_sync"
	"adapter/internal/shared/constants"
	"adapter/internal/shared/utils"
	"github.com/gofiber/fiber/v2"
)

type RegistrySyncHandler struct {
	ondcService *ondc.ONDCService
}

func NewRegistrySyncHandler(ondcService *ondc.ONDCService) *RegistrySyncHandler {
	return &RegistrySyncHandler{ondcService: ondcService}
}

func (h *RegistrySyncHandler) SyncRegistry(c *fiber.Ctx) error {
	var req ports.SyncRegistryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrInvalidRequestBody,
		})
	}

	if req.RegistryEnv == "" || len(req.Domains) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false,
			Message: "registry_env and domains are required",
		})
	}

	response, err := h.ondcService.SyncRegistry(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrFailedToStartRegistrySync,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ApiResponse{
		Success: true,
		Message: "Registry sync completed successfully",
		Data:    response,
	})
}