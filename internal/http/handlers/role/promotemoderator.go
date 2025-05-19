package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)

func PromoteToModerator(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Get the role ID for "Moderator"
	var roleID string
	err = db.QueryRow("SELECT id FROM roles WHERE name = 'moderator'").Scan(&roleID)
	if err == sql.ErrNoRows {
		http.Error(w, "Moderator role not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch moderator role", http.StatusInternalServerError)
		return
	}

	// Update the user's main role in the users table
	_, err = db.Exec("UPDATE users SET user_type = 'moderator', updated_at = NOW() WHERE id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to update user type", http.StatusInternalServerError)
		return
	}

	// Update the user's role in the user_roles table
	_, err = db.Exec("UPDATE user_roles SET role_id = $1 WHERE user_id = $2", roleID, userID)
	if err != nil {
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "User promoted to Moderator successfully",
		"data":    nil,
	}

	json.NewEncoder(w).Encode(response)

}
