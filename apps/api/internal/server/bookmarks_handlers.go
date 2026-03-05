package server

import (
	"encoding/json"
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



func(h *BookmarkHandler) AddBookmark(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req AddBookmarkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err  := h.Repo.AddBookmark(userID, req.HackathonID)
	if err != nil {
		http.Error(w, "Failed to add bookmark", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Added bookmark"})
}	

func(h *BookmarkHandler) RemoveBookmark(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}


	var req RemoveBookmarkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err  := h.Repo.RemoveBookmark(userID, req.HackathonID)
	if err != nil {
		http.Error(w, "Failed to remove bookmark", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Removed bookmark"})
}

func(h *BookmarkHandler) GetBookmarksByUser(w http.ResponseWriter, r *http.Request){

	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bookmarks, err := h.Repo.GetBookmarksByUser(userID)
	if err != nil {
		http.Error(w, "Failed to get bookmarks", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookmarks)
}