package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

func DeleteRequest(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	deleteID := mux.Vars(r)["user_id"]

	if userID == "" || deleteID == "" {
		http.Error(w, "You cannot access this page!", http.StatusBadRequest)
		return
	}

	if userID != deleteID {
		http.Error(w, "You can only delete your own account", http.StatusForbidden)
		return
	}

	query := `UPDATE users SET deletion_requested = true, updated_at = NOW() WHERE id = $1`

	// Connect to the database
	db:=database.Connect()

	_, err := db.Exec(query, userID)
	if err != nil {
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Your account deletion request has been submitted successfully",
		"data":    nil,
	}

	json.NewEncoder(w).Encode(response)

}
