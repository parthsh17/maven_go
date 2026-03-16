package store

import (
	"fmt"
	"maven/internal/models"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	mu    sync.RWMutex
	users map[string]*models.User // keyed by email
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*models.User),
	}
}

// Register hashes the password with bcrypt and stores the user.
// Returns an error if the email is already registered.
func (us *UserStore) Register(id, email, password string) (*models.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.users[email]; exists {
		return nil, fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	fmt.Printf("[AUTH] Generated Bcrypt Hash for %s: %s\n", email, string(hash))

	user := models.NewUser(id, email, string(hash))
	us.users[email] = user
	fmt.Printf("[STORE] User %s successfully stored in in-memory database\n", email)
	return user, nil
}

// Authenticate verifies the email and password.
// Returns the User if credentials are valid, or an error otherwise.
func (us *UserStore) Authenticate(email, password string) (*models.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	user, exists := us.users[email]
	if !exists {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// GetUser returns a user by email (read-only lookup).
func (us *UserStore) GetUser(email string) (*models.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	user, exists := us.users[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
