package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

func PromoteToModerator(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]

	if userID == "" {
		// http.Error(w, "User ID is required", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Connect to the database
	db := database.Connect()

	// Get the role ID for "Moderator"
	var roleID string
	err := db.QueryRow("SELECT id FROM roles WHERE name = 'moderator'").Scan(&roleID)
	if err == sql.ErrNoRows {
		// http.Error(w, "Moderator role not found", http.StatusNotFound)
		utils.ErrorResponse(w, http.StatusNotFound, "Moderator role not found")
		return
	} else if err != nil {
		// http.Error(w, "Failed to fetch moderator role", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch moderator role")
		return
	}
	var currentRole string
	db.QueryRow("SELECT user_type FROM users WHERE id = $1", userID).Scan(&currentRole)
	fmt.Println(currentRole)

	//user chara kaoke moderator korte parbona
	if currentRole!="user"{
		// http.Error(w, "User not found", http.StatusNotFound)
		utils.ErrorResponse(w, http.StatusNotFound, "Admins cannot be promoted to Moderator!")
		return
	}
	// Update the user's main role in the users table
	_, err = db.Exec("UPDATE users SET user_type = 'moderator', updated_at = NOW() WHERE id = $1", userID)
	if err != nil {
		// http.Error(w, "Failed to update user type", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user type")
		return
	}

	// Update the user's role in the user_roles table
	_, err = db.Exec("UPDATE user_roles SET role_id = $1 WHERE user_id = $2", roleID, userID)
	if err != nil {
		// http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	// Return a success response
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// response := map[string]interface{}{
	// 	"status":  strconv.Itoa(http.StatusOK),
	// 	"message": "User promoted to Moderator successfully",
	// 	"data":    nil,
	// }

	// json.NewEncoder(w).Encode(response)
	utils.SuccessResponse(w, http.StatusOK, "User promoted to Moderator successfully", nil)

}
