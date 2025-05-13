package models

import "time"

type User struct {
	ID                string     `json:"id"`
	Username          string     `json:"username"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	EmailVerified     bool       `json:"email_verified"`
	UserType          string     `json:"user_type"`
	Active            bool       `json:"active"`
	VerificationToken string     `json:"-"`
	TokenExpiry       *time.Time `json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
