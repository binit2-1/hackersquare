package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
	"github.com/binit2-1/hackersquare/apps/api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo  domain.UserRepository
	AiService domain.AIService
}

func NewAuthHandler(repo domain.UserRepository, aiService domain.AIService) *AuthHandler {
	return &AuthHandler{
		UserRepo:  repo,
		AiService: aiService,
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

type UpdateProfileReadmeRequest struct {
	Readme string `json:"readme"`
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

	if err := h.setAuthCookie(w, user); err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

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
	if err != nil || user == nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := h.setAuthCookie(w, user); err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})

}

func (h *AuthHandler) setAuthCookie(w http.ResponseWriter, user *domain.User) error {

	tokenString, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return err
	}

	baseURL := strings.Trim(strings.ToLower(os.Getenv("NEXT_APP_BASE_URL")), "'\"")
	useSecureCookie := true
	if strings.HasPrefix(baseURL, "http://localhost") || strings.HasPrefix(baseURL, "http://127.0.0.1") {
		useSecureCookie = false
	}
	sameSiteMode := http.SameSiteLaxMode
	if useSecureCookie {
		sameSiteMode = http.SameSiteNoneMode
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   useSecureCookie,
		SameSite: sameSiteMode,
	}

	http.SetCookie(w, cookie)
	return nil
}

func (h *AuthHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	baseURL := strings.Trim(strings.ToLower(os.Getenv("NEXT_APP_BASE_URL")), "'\"")
	useSecureCookie := true
	if strings.HasPrefix(baseURL, "http://localhost") || strings.HasPrefix(baseURL, "http://127.0.0.1") {
		useSecureCookie = false
	}
	sameSiteMode := http.SameSiteLaxMode
	if useSecureCookie {
		sameSiteMode = http.SameSiteNoneMode
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   useSecureCookie,
		SameSite: sameSiteMode,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found in context", http.StatusUnauthorized)
		return
	}

	var req domain.ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.UserRepo.UpdateUserProfile(userID, req); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func (h *AuthHandler) ConnectGithub(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	redirectURI := os.Getenv("NEXT_APP_BASE_URL") + "/api/v1/auth/github/callback"
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=read:user&redirect_uri=%s", clientID, redirectURI)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: Session lost during GitHub redirect", http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Failed to get authorization code from github", http.StatusBadRequest)
		return
	}

	tokenReqBody, _ := json.Marshal(map[string]string{
		"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
		"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		"code":          code,
	})

	tokenReq, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(tokenReqBody))
	tokenReq.Header.Set("Content-Type", "application/json")
	tokenReq.Header.Set("Accept", "application/json")

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil || tokenResp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to exchange code for access token", http.StatusInternalServerError)
		return
	}
	defer tokenResp.Body.Close()

	var tokenData struct {
		AccessToken      string `json:"access_token"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err != nil {
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
	}

	if tokenData.Error != "" {
		http.Error(w, fmt.Sprintf("GitHub Error: %s", tokenData.ErrorDescription), http.StatusUnauthorized)
		return
	}

	if tokenData.AccessToken == "" {
		http.Error(w, "Failed to retrieve access token", http.StatusUnauthorized)
		return
	}

	userReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)

	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		http.Error(w, "Failed to fetch GitHub profile", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	var githubUser struct {
		Login string `json:"login"`
	}

	if err := json.NewDecoder(userResp.Body).Decode(&githubUser); err != nil {
		http.Error(w, "Failed to parse GitHub user response", http.StatusInternalServerError)
		return
	}

	// 3. Save the handle to your PostgreSQL Database
	err = h.UserRepo.LinkGithubHandle(userID, githubUser.Login)
	if err != nil {
		http.Error(w, "Failed to save GitHub handle to database", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GithubLogin(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	redirectURI := os.Getenv("NEXT_APP_BASE_URL") + "/api/v1/auth/github/login/callback"
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=read:user,user:email&redirect_uri=%s", clientID, redirectURI)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GithubLoginCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/?error=missing_code", http.StatusTemporaryRedirect)
		return
	}

	tokenReqBody, _ := json.Marshal(map[string]string{
		"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
		"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		"code":          code,
	})

	tokenReq, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(tokenReqBody))
	tokenReq.Header.Set("Content-Type", "application/json")
	tokenReq.Header.Set("Accept", "application/json")

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		http.Redirect(w, r, "/?error=github_unreachable", http.StatusTemporaryRedirect)
		return
	}
	defer tokenResp.Body.Close()

	var tokenData struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(tokenResp.Body).Decode(&tokenData)

	if tokenData.AccessToken == "" {
		http.Redirect(w, r, "/?error=no_token", http.StatusTemporaryRedirect)
		return
	}

	userReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil || userResp == nil {
		http.Redirect(w, r, "/?error=failed_to_fetch_user", http.StatusTemporaryRedirect)
		return
	}
	defer userResp.Body.Close()

	var githubUser struct {
		Login string `json:"login"`
		Name  string `json:"name"`
	}
	json.NewDecoder(userResp.Body).Decode(&githubUser)

	emailReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	emailReq.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
	emailResp, err := http.DefaultClient.Do(emailReq)
	if err != nil {
		http.Redirect(w, r, "/?error=failed_to_fetch_emails", http.StatusTemporaryRedirect)
		return
	}
	defer emailResp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	json.NewDecoder(emailResp.Body).Decode(&emails)

	var primaryEmail string
	for _, e := range emails {
		if e.Primary && e.Verified {
			primaryEmail = e.Email
			break
		}
	}

	if primaryEmail == "" {
		http.Redirect(w, r, "/?error=no_verified_email", http.StatusTemporaryRedirect)
		return
	}

	user, err := h.UserRepo.GetUserByEmail(primaryEmail)
	if user == nil {
		name := githubUser.Name
		if name == "" {
			name = githubUser.Login
		}

		randomBytes := make([]byte, 16)
		rand.Read(randomBytes)
		dummyPassword := hex.EncodeToString(randomBytes)

		newUser := &domain.User{
			Name:         name,
			Email:        primaryEmail,
			PasswordHash: dummyPassword,
		}

		err := h.UserRepo.CreateUser(newUser)
		if err != nil {
			http.Redirect(w, r, "/?error=account_creation_failed", http.StatusTemporaryRedirect)
			return
		}

		_ = h.UserRepo.LinkGithubHandle(newUser.ID, githubUser.Login)
		user = newUser
	}

	if err := h.setAuthCookie(w, user); err != nil {
		http.Redirect(w, r, "/?error=session_creation_failed", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func (h *AuthHandler) GenerateProfileSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found in context", http.StatusUnauthorized)
		return
	}

	user, err := h.UserRepo.GetUserByID(userID)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	profileSnapshot := fmt.Sprintf(
		"name: %s\nheadline: %s\nlocation: %s\ngithub_handle: %s\nwebsite: %s\nlinkedin: %s\ntwitter: %s",
		user.Name,
		user.Headline,
		user.Location,
		user.GithubHandle,
		user.WebsiteURL,
		user.LinkedinURL,
		user.TwitterURL,
	)

	githubData := profileSnapshot
	if user.GithubHandle != "" {
		if fetchedData, fetchErr := utils.FetchGithubData(user.GithubHandle); fetchErr == nil {
			githubData = fetchedData
		} else {
			fmt.Printf("GitHub fetch failed for %s, using profile fallback: %v\n", user.GithubHandle, fetchErr)
		}
	}

	aiCtx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	summary, err := h.AiService.GenerateProfileReadme(aiCtx, githubData)
	if err != nil {
		http.Error(w, fmt.Sprintf("AI generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	err = h.UserRepo.UpdateProfileReadme(userID, summary)
	if err != nil {
		http.Error(w, "Failed to save summary to database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "AI profile generated successfully",
		"summary": summary,
	})
}

func (h *AuthHandler) UpdateProfileReadme(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found in context", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileReadmeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.UserRepo.UpdateProfileReadme(userID, req.Readme); err != nil {
		http.Error(w, "Failed to update profile readme", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile readme updated successfully"})
}
