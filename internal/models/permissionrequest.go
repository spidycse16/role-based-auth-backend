package models

// PermissionRequest represents the JSON payload for creating a permission
type PermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
}