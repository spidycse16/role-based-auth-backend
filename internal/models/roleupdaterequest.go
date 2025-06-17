package models

// RoleUpdateRequest represents the expected request body for updating a role
type RoleUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}