package models

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
