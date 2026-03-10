package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
	"github.com/binit2-1/hackersquare/apps/api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo domain.UserRepository
}

func NewAuthHandler(repo domain.UserRepository) *AuthHandler {
	return &AuthHandler{
		UserRepo: repo,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user := &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	if err := h.UserRepo.CreateUser(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	h.setAuthCookie(w, user)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetUserByEmail(req.Email)
	if err != nil{
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	h.setAuthCookie(w, user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})

}

func (h *AuthHandler) setAuthCookie(w http.ResponseWriter, user *domain.User) {
	
	tokenString, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate auth token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		Expires:   time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, 
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func (h *AuthHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found in context", http.StatusUnauthorized)
		return
	}
	
	user, err:= h.UserRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}






func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found in context", http.StatusUnauthorized)
		return
	}

	user, err := h.UserRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}
	
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

}