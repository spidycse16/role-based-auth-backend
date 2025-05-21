package middleware

import (
	"context"
	// "fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sagorsarker04/Developer-Assignment/internal/config"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "No valid authentication token", http.StatusUnauthorized)
			return
		}

		// Load the config
		// cfg, err := config.LoadConfig()
		// if err != nil {
		// 	http.Error(w, "Failed to load config", http.StatusInternalServerError)
		// 	return
		// }
		//new singletone class
		cfg:=config.GetConfig()

		// Parse the token
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, _ := claims["user_id"].(string)
			username, _ := claims["username"].(string)
			userType, _ := claims["user_type"].(string)

			// Set values in the context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UsernameKey, username)
			ctx = context.WithValue(ctx, UserTypeKey, userType)
			//fmt.Println(ctx)

			// Pass the request to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}
