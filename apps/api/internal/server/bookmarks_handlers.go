package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
)


type BookmarkHandler struct{
	Repo domain.BookmarkRepository
}

func NewBookmarkHandler(repo domain.BookmarkRepository) *BookmarkHandler {
	return &BookmarkHandler{
		Repo: repo,
	}
}

type AddBookmarkRequest struct {
	HackathonID string `json:"hackathon_id"`
}

type RemoveBookmarkRequest struct {
	HackathonID string `json:"hackathon_id"`
}

func getUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value("userID").(string)
	if ok && userID != "" {
		return userID, true
	}
	userID, ok = r.Context().Value("user_id").(string)
	if ok && userID != "" {
		return userID, true
	}
	return "", false
}

func getHackathonID(r *http.Request) (string, error) {
	if id := r.URL.Query().Get("hackathon_id"); id != "" {
		return id, nil
	}

	var req AddBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return "", err
	}
	return req.HackathonID, nil
}



func(h *BookmarkHandler) AddBookmark(w http.ResponseWriter, r *http.Request){
	userID, ok := getUserID(r)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hackathonID, err := getHackathonID(r)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if hackathonID == "" {
		http.Error(w, "Missing hackathon_id", http.StatusBadRequest)
		return
	}

	err  = h.Repo.AddBookmark(userID, hackathonID)
	if err != nil {
		fmt.Printf("AddBookmark DB error: %v\n", err)
		http.Error(w, "Failed to add bookmark", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Added bookmark"})
}	

func(h *BookmarkHandler) RemoveBookmark(w http.ResponseWriter, r *http.Request){
	userID, ok := getUserID(r)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hackathonID, err := getHackathonID(r)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if hackathonID == "" {
		http.Error(w, "Missing hackathon_id", http.StatusBadRequest)
		return
	}

	err  = h.Repo.RemoveBookmark(userID, hackathonID)
	if err != nil {
		fmt.Printf("RemoveBookmark DB error: %v\n", err)
		http.Error(w, "Failed to remove bookmark", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Removed bookmark"})
}

func(h *BookmarkHandler) GetBookmarksByUser(w http.ResponseWriter, r *http.Request){

	userID, ok := getUserID(r)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bookmarks, err := h.Repo.GetBookmarksByUser(userID)
	if err != nil {
		fmt.Printf("GetBookmarksByUser DB error: %v\n", err)
		http.Error(w, "Failed to get bookmarks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookmarks)
}
