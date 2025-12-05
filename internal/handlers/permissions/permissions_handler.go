package handlers

import (
	"adapter/internal/domain/permissions"
	permissionPorts "adapter/internal/ports/permissions"
	"adapter/internal/shared/constants"
	"adapter/internal/shared/utils"
	"github.com/gofiber/fiber/v2"
)

type PermissionsHandler struct {
	permissionsService *permissions.PermissionsService
}

func NewPermissionsHandler(permissionsService *permissions.PermissionsService) *PermissionsHandler {
	return &PermissionsHandler{permissionsService: permissionsService}
}

func (h *PermissionsHandler) UpdatePermissions(c *fiber.Ctx) error {
	var req struct {
		Updates []permissionPorts.PermissionsUpdateRequest `json:"updates"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrInvalidRequestBody,
		})
	}

	if len(req.Updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrUpdatesArrayEmpty,
		})
	}

	results, err := h.permissionsService.UpdatePermissions(req.Updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrFailedToUpdatePermissions,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ApiResponse{
		Success: true,
		Message: "Permissions updated successfully",
		Data:    fiber.Map{"results": results},
	})
}

func (h *PermissionsHandler) QueryPermissions(c *fiber.Ctx) error {
	var req permissionPorts.PermissionsQueryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrInvalidRequestBody,
		})
	}

	if req.BapID == "" || req.Domain == "" || req.RegistryEnv == "" || len(req.SellerIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrRequiredPermissionsFields,
		})
	}

	response, err := h.permissionsService.QueryPermissions(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ApiResponse{
			Success: false,
			Message: constants.ErrFailedToQueryPermissions,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.ApiResponse{
		Success: true,
		Message: "Permissions queried successfully",
		Data:    response,
	})
}