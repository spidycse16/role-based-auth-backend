package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

func PromoteToAdmin(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	currentUserType := middleware.GetUserType(r)

	if currentUserType != "system_admin" {
		http.Error(w, "Tumi to system admin na vai", http.StatusForbidden)
		return
	}

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Databse connect hocchena vai", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Check if the user is already an Admin or SystemAdmin
	var currentRole string
	err = db.QueryRow("SELECT user_type FROM users WHERE id = $1", userID).Scan(&currentRole)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch user role", http.StatusInternalServerError)
		return
	}

	if currentRole == "admin" || currentRole == "system_admin" {
		http.Error(w, "User is already an Admin or System Admin", http.StatusBadRequest)
		return
	}

	// Promote the user to Admin
	_, err = db.Exec("UPDATE users SET user_type = 'admin', updated_at = NOW() WHERE id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to promote user to Admin", http.StatusInternalServerError)
		return
	}

	// Update user_roles table to reflect the new role
	var adminRoleID string
	err = db.QueryRow("SELECT id FROM roles WHERE name = 'admin'").Scan(&adminRoleID)
	if err == sql.ErrNoRows {
		http.Error(w, "Admin role not found", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch Admin role ID", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE user_roles SET role_id = $1 WHERE user_id = $2", adminRoleID, userID)
	if err != nil {
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully promoted to Admin"))
}
