package models

import (
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // never serialised to JSON
	CreatedAt    string `json:"created_at"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignupRequest) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(r.Email, "@") {
		return fmt.Errorf("email is invalid")
	}
	if len(r.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(id, email, passwordHash string) *User {
	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	}
}
