package tests

import (
	"maven/internal/models"
	"maven/internal/store"
	"testing"
)

func TestUserStore_RegisterAndAuthenticate(t *testing.T) {
	us := store.NewUserStore()

	user, err := us.Register("user-001", "alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email 'alice@example.com', got %q", user.Email)
	}
	if user.PasswordHash == "" {
		t.Error("expected non-empty PasswordHash")
	}

	// PasswordHash must not equal plain text password
	if user.PasswordHash == "secret123" {
		t.Error("PasswordHash should not equal plaintext password")
	}

	got, err := us.Authenticate("alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}
	if got.ID != user.ID {
		t.Errorf("expected ID %q, got %q", user.ID, got.ID)
	}
}

func TestUserStore_DuplicateEmail(t *testing.T) {
	us := store.NewUserStore()
	_, _ = us.Register("u1", "bob@example.com", "pass1234")
	_, err := us.Register("u2", "bob@example.com", "pass5678")
	if err == nil {
		t.Error("expected duplicate email error, got nil")
	}
}

func TestUserStore_WrongPassword(t *testing.T) {
	us := store.NewUserStore()
	_, _ = us.Register("u3", "carol@example.com", "rightpass")

	_, err := us.Authenticate("carol@example.com", "wrongpass")
	if err == nil {
		t.Error("expected authentication error for wrong password, got nil")
	}
}

func TestUserStore_UserNotFound(t *testing.T) {
	us := store.NewUserStore()
	_, err := us.Authenticate("nobody@example.com", "any")
	if err == nil {
		t.Error("expected not-found error, got nil")
	}
}

func TestSignupRequest_Validate_EmptyEmail(t *testing.T) {
	req := models.SignupRequest{Email: "", Password: "pass123"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for empty email, got nil")
	}
}

func TestSignupRequest_Validate_InvalidEmail(t *testing.T) {
	req := models.SignupRequest{Email: "notanemail", Password: "pass123"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for invalid email, got nil")
	}
}

func TestSignupRequest_Validate_ShortPassword(t *testing.T) {
	req := models.SignupRequest{Email: "test@test.com", Password: "abc"}
	if err := req.Validate(); err == nil {
		t.Error("expected error for short password, got nil")
	}
}

func TestSignupRequest_Validate_Valid(t *testing.T) {
	req := models.SignupRequest{Email: "valid@test.com", Password: "secure123"}
	if err := req.Validate(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}
