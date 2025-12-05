package utils

import (
	"strings"
)

// ApiResponse is the standardized JSON response structure for all APIs.
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SplitAndTrim splits a comma-separated string into a slice of strings and trims whitespace.
func SplitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// JoinConditions joins a slice of SQL conditions with a given separator.
func JoinConditions(conditions []string, separator string) string {
	return strings.Join(conditions, separator)
}
