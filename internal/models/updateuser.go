// models/user.go
package models

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
