package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
	"github.com/sagorsarker04/Developer-Assignment/internal/utils"
)

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	UserType      string `json:"user_type"`
	EmailVerified bool   `json:"email_verified"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ListAllUsers handler called")
	userType := middleware.GetUserType(r)
	fmt.Println("User Type:", userType)

	db := database.Connect()

	rows, err := db.Query(`
		SELECT id, username, email, first_name, last_name, user_type, email_verified, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	fmt.Println(rows)
	if err != nil {
		log.Println("Query Error:", err)
		// http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email,
			&user.FirstName, &user.LastName, &user.UserType,
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			log.Println("Scan Error:", err)
			// http.Error(w, "Failed to read user data", http.StatusInternalServerError)
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to read user data")
			return
		}
		users = append(users, user)
	}

	// response := map[string]interface{}{
	// 	"status":  strconv.Itoa(http.StatusOK),
	// 	"message": "Users retrieved successfully",
	// 	"data":    users,
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// if err := json.NewEncoder(w).Encode(response); err != nil {
	// 	log.Println("JSON Encode Error:", err)
	// 	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	// }
	utils.SuccessResponse(w, http.StatusOK, "Users retrieved successfully", users)
}
