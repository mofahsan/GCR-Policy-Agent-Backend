package appError

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// CustomError defines a standard application error.
type CustomError struct {
	HTTPStatusCode int    `json:"-"`
	Code           string `json:"code"`
	Message        string `json:"message"`
}

// Error returns the string representation of the error.
func (e *CustomError) Error() string {
	return fmt.Sprintf("error: code=%s, message=%s", e.Code, e.Message)
}

// NewCustomError creates a new CustomError instance.
func NewCustomError(httpStatusCode int, code, message string) *CustomError {
	return &CustomError{
		HTTPStatusCode: httpStatusCode,
		Code:           code,
		Message:        message,
	}
}

// Pre-defined application errors
var (
	// Generic Errors
	ErrInvalidRequestBody = NewCustomError(400, "GEN_4001", "Request body is not valid")
	ErrInternalServer     = NewCustomError(500, "GEN_5001", "An unexpected internal server error occurred")
	ErrDatabaseConnectionFailed = NewCustomError(500, "DB_5001", "Failed to connect to the database")

	// Catalog Sync Errors
	ErrDomainRequired     = NewCustomError(400, "CAT_4001", "domain query parameter is required")
	ErrSellerAndDomain    = NewCustomError(400, "CAT_4002", "seller_id path parameter and domain query parameter are required")
	ErrRecordNotFound     = NewCustomError(404, "CAT_4004", "Record not found for the specified parameters")
	ErrGetPendingSellers  = NewCustomError(500, "CAT_5001", "Failed to get pending catalog sync sellers")
	ErrGetSyncStatus      = NewCustomError(500, "CAT_5002", "Failed to get sync status")
	
	// Permissions Errors
	ErrProcessingPermissions = NewCustomError(500, "PERM_5001", "Error processing permissions request")

	// Registry Sync Errors
	ErrRegistrySync = NewCustomError(500, "REG_5001", "Failed to start registry sync")
)


// ErrorHandler is a custom Fiber error handler.
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default to a 500 internal server error
		customErr := ErrInternalServer

		// Check if the error is a CustomError
		if e, ok := err.(*CustomError); ok {
			customErr = e
		}

		return c.Status(customErr.HTTPStatusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    customErr.Code,
				"message": customErr.Message,
			},
		})
	}
}