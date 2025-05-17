// handlers/user.go
package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/models"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

// UpdateUser updates the details of a specific user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL path
	userID := mux.Vars(r)["user_id"]

	// Verify the authenticated user is updating their own profile
	authUserID := middleware.GetUserID(r)
	if userID != authUserID {
		http.Error(w, "Forbidden: You can only update your own profile", http.StatusForbidden)
		return
	}

	// Parse the request body
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and Email are required", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Update the user in the database
	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, email, created_at, updated_at
	`
	var updatedUser struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	err = db.QueryRow(query, req.Name, req.Email, userID).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Return the updated user details as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}
