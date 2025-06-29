package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

func getUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

// fetchSystemAdminIDs returns a slice of user IDs (as strings) who have the system admin role
func fetchSystemAdminIDs(db *sql.DB) ([]string, error) {
	var systemAdminIDs []string

	// Adjust the query based on your schema
	// For example:
	// select user_id from user_roles join roles on user_roles.role_id = roles.id where roles.name = 'system_admin'
	query := `
		SELECT u.id
		FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		JOIN roles r ON ur.role_id = r.id
		WHERE r.name = 'system_admin'
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		systemAdminIDs = append(systemAdminIDs, id)
	}

	return systemAdminIDs, nil
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	deleteID := mux.Vars(r)["user_id"]
	if deleteID == "" {
		// http.Error(w, "Invalid user ID", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		// http.Error(w, "Unauthorized", http.StatusUnauthorized)
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	db := database.Connect()

	systemAdminIDs, err := fetchSystemAdminIDs(db)
	if err != nil {
		// http.Error(w, "Failed to fetch system admin info", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch system admin info")
		return
	}

	// Check if deleteID is current user or system admin
	if deleteID == currentUserID {
		// http.Error(w, "You cannot delete your own account", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "You cannot delete your own account")
		return
	}

	// Check if deleteID is in systemAdminIDs slice
	for _, adminID := range systemAdminIDs {
		if deleteID == adminID {
			// http.Error(w, "You cannot delete a system admin account", http.StatusForbidden)
			utils.ErrorResponse(w, http.StatusForbidden, "You cannot delete a system admin account")
			return
		}
	}

	query := `DELETE FROM users WHERE id = $1 AND deletion_requested = true`

	res, err := db.Exec(query, deleteID)
	if err != nil {
		// http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to execute query")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		// http.Error(w, "User not found or deletion not requested", http.StatusNotFound)
		utils.ErrorResponse(w, http.StatusNotFound, "User not found or deletion not requested")
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// response := map[string]interface{}{
	// 	"status":  strconv.Itoa(http.StatusOK),
	// 	"message": "User account deleted successfully",
	// 	"data":    nil,
	// }

	// json.NewEncoder(w).Encode(response)
	utils.SuccessResponse(w, http.StatusOK, "User account deleted successfully", nil)
}
