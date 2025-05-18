// middleware/permissions.go
package middleware

import (
	"fmt"
	"net/http"

	"github.com/sagorsarker04/Developer-Assignment/internal/database"
)

// FetchAllUserPermissions returns all permissions of a given user
func FetchAllUserPermissions(userID string) ([]string, error) {
	db, err := database.Connect()
	if err != nil {
		fmt.Println("[ERROR] Failed to connect to database:", err)
		return nil, err
	}
	defer database.Close(db)

	query := `
		SELECT p.name
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE ur.user_id = $1
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Println("[ERROR] Failed to execute query:", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			fmt.Println("[ERROR] Failed to read permission:", err)
			continue
		}
		permissions = append(permissions, perm)
	}

	// Debug print: Show all fetched permissions
	fmt.Println("[DEBUG] Fetched Permissions for User", userID, ":", permissions)

	return permissions, nil
}


// CheckPermission checks if any required permission is in the available permissions list
func CheckPermission(required []string, available []string) bool {
	for _, reqPerm := range required {
		for _, availPerm := range available {
			if reqPerm == availPerm {
				return true
			}
		}
	}
	return false
}

// RequireAnyPermission checks if the user has at least one of the required permissions
func RequireAnyPermission(required []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		if userID == "" {
			http.Error(w, "No Valid user", http.StatusForbidden)
			return
		}

		availablePermissions, err := FetchAllUserPermissions(userID)
		if err != nil || len(availablePermissions) == 0 {
			http.Error(w, "No valiable permissions", http.StatusForbidden)
			return
		}

		if !CheckPermission(required, availablePermissions) {
			http.Error(w, "No matching Permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
