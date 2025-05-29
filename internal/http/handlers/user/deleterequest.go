package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

func DeleteRequest(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	deleteID := mux.Vars(r)["user_id"]
	userType := middleware.GetUserType(r)
	fmt.Println(userType)
	if userType == "system_admin" {
		// http.Error(w, "System admin cant be deleted", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "System admin cant be deleted")
		return
	}
	if userID == "" || deleteID == "" {
		// http.Error(w, "You cannot access this page!", http.StatusBadRequest)
		utils.ErrorResponse(w, http.StatusBadRequest, "You cannot access this page!")
		return
	}

	if userID != deleteID {
		// http.Error(w, "You can only delete your own account", http.StatusForbidden)
		utils.ErrorResponse(w, http.StatusForbidden, "You can only delete your own account")
		return
	}

	query := `UPDATE users SET deletion_requested = true, updated_at = NOW() WHERE id = $1`

	// Connect to the database
	db := database.Connect()

	_, err := db.Exec(query, userID)
	if err != nil {
		// http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to execute query")
		return
	}

	// Return a success response
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// response := map[string]interface{}{
	// 	"status":  strconv.Itoa(http.StatusOK),
	// 	"message": "Your account deletion request has been submitted successfully",
	// 	"data":    nil,
	// }

	// json.NewEncoder(w).Encode(response)

	utils.SuccessResponse(w, http.StatusOK, "Your account deletion request has been submitted successfully", nil)

}
