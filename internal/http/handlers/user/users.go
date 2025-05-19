package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sagorsarker04/Developer-Assignment/internal/database"
	"github.com/sagorsarker04/Developer-Assignment/internal/http/middleware"
)

func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	// Extract the user type from the context
	userType := middleware.GetUserType(r)
	fmt.Println(r)
	fmt.Println("User Type:", userType)

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close(db)

	// Fetch all users
	rows, err := db.Query(`
		SELECT id, username, email, first_name, last_name, user_type, email_verified, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect users
	var users []map[string]interface{}
	for rows.Next() {
		var id, username, email, firstName, lastName, userType string
		var emailVerified bool
		var createdAt, updatedAt string

		if err := rows.Scan(&id, &username, &email, &firstName, &lastName, &userType, &emailVerified, &createdAt, &updatedAt); err != nil {
			http.Error(w, "Failed to read user data", http.StatusInternalServerError)
			return
		}

		users = append(users, map[string]interface{}{
			"id":             id,
			"username":       username,
			"email":          email,
			"first_name":     firstName,
			"last_name":      lastName,
			"user_type":      userType,
			"email_verified": emailVerified,
			"created_at":     createdAt,
			"updated_at":     updatedAt,
		})
	}

	// Return the users as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  strconv.Itoa(http.StatusOK),
		"message": "Users retrieved successfully",
		"data":    users,
	}

	json.NewEncoder(w).Encode(response)

}
