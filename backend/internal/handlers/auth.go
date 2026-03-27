package handlers

import (
	"encoding/json"
	"fmt"
	"maven/internal/models"
	"maven/internal/store"
	"net/http"
)

type AuthHandler struct {
	userStore store.UserStore
}

func NewAuthHandler(us store.UserStore) *AuthHandler {
	fmt.Println("[AUTH] Initializing AuthHandler with thread-safe UserStore (sync.RWMutex)")
	return &AuthHandler{userStore: us}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[AUTH] Received Signup request for email: %s\n", r.URL.Path)
	fmt.Println("[CONCURRENCY] Handling Signup request in a new goroutine (standard net/http behavior)")
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	id := newUUID()
	user, err := h.userStore.Register(id, req.Email, req.Password)
	if err != nil {
		errorResponse(w, http.StatusConflict, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	user, err := h.userStore.Authenticate(req.Email, req.Password)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, user)
}
