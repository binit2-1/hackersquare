package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
)

type HackathonHandler struct {
	Repo      domain.HackathonRepository
	UserRepo  domain.UserRepository
	AIService domain.AIService
}

func NewHackathonHandler(repo domain.HackathonRepository, userRepo domain.UserRepository, aiService domain.AIService) *HackathonHandler {
	return &HackathonHandler{
		Repo:      repo,
		UserRepo:  userRepo,
		AIService: aiService,
	}
}

func (h *HackathonHandler) SearchHackathons(w http.ResponseWriter, r *http.Request) {

	queryValues := r.URL.Query()

	filters := domain.SearchFilters{
		Query:      queryValues.Get("q"),
		Location:   queryValues.Get("location"),
		PrizeRange: queryValues.Get("prizeRange"),
		Status:     queryValues.Get("status"),
		Page:       1,
		Limit:      20,
	}

	if pageStr := queryValues.Get("page"); pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil && val > 0 {
			filters.Page = val
		}
	}

	if limitStr := queryValues.Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = val
		}
	}

	hackathons, totalCount, err := h.Repo.SearchHackathons(filters)
	if err != nil {
		fmt.Printf("Database Search Error: %v\n", err)

		http.Error(w, "Failed to search hackathons", http.StatusInternalServerError)
		return
	}

	if hackathons == nil {
		hackathons = []domain.Hackathon{}
	}

	if filters.Limit <= 0 {
		filters.Limit = 20
	}

	totalPages := (totalCount + filters.Limit - 1) / filters.Limit

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]any{
		"data": hackathons,
		"metadata": map[string]any{
			"totalRecords": totalCount,
			"currentPage":  filters.Page,
			"limit":        filters.Limit,
			"totalPages":   totalPages,
		},
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *HackathonHandler) NearbyHackathons(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	city := strings.TrimSpace(queryValues.Get("city"))
	country := strings.TrimSpace(queryValues.Get("country"))

	if strings.EqualFold(city, "unknown") {
		city = ""
	}
	if strings.EqualFold(country, "unknown") {
		country = ""
	}

	// Raw ISO-2 country codes (e.g. IN, US) produce noisy ILIKE matches.
	if len([]rune(country)) == 2 {
		country = ""
	}

	page := 1
	limit := 20

	if pageStr := queryValues.Get("page"); pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil && val > 0 {
			page = val
		}
	}

	if limitStr := queryValues.Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	if city == "" && country == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]any{
			"data": []domain.Hackathon{},
			"metadata": map[string]any{
				"totalRecords": 0,
				"currentPage":  page,
				"limit":        limit,
				"totalPages":   0,
			},
		}

		if err := json.NewEncoder(w).Encode(&response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}

	hackathons, totalCount, err := h.Repo.NearbyHackathons(city, country, page, limit)
	if err != nil {
		fmt.Printf("Nearby Query Error: %v\n", err)
		http.Error(w, "Failed to fetch nearby hackathons", http.StatusInternalServerError)
		return
	}

	if hackathons == nil {
		hackathons = []domain.Hackathon{}
	}

	totalPages := (totalCount + limit - 1) / limit

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]any{
		"data": hackathons,
		"metadata": map[string]any{
			"totalRecords": totalCount,
			"currentPage":  page,
			"limit":        limit,
			"totalPages":   totalPages,
		},
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *HackathonHandler) GetSearchOverview(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	query := r.URL.Query().Get("q")

	if query == "" {
		query = "upcoming hackathons"
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userProfile, err := h.UserRepo.GetUserByID(userID)
	if err != nil || userProfile == nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	profileContext := strings.TrimSpace(userProfile.ProfileReadme)
	if profileContext == "" {
		profileContext = fmt.Sprintf(
			"Name: %s\nHeadline: %s\nLocation: %s\nGitHub: %s",
			userProfile.Name,
			userProfile.Headline,
			userProfile.Location,
			userProfile.GithubHandle,
		)
	}

	filters := domain.SearchFilters{
		Query:      query,
		Location:   queryValues.Get("location"),
		PrizeRange: queryValues.Get("prizeRange"),
		Status:     queryValues.Get("status"),
		Page:       1,
		Limit:      5,
	}

	hackathons, _, err := h.Repo.SearchHackathons(filters)
	if err != nil {
		http.Error(w, "Failed to fetch hackathon context", http.StatusInternalServerError)
		return
	}

	var contextBuilder strings.Builder
	for i, hack := range hackathons {
		if i >= 3 {
			break
		}
		contextBuilder.WriteString(fmt.Sprintf("- Title: %s\n  Location: %s\n  Tags/Tech: %s\n\n",
			hack.Title, hack.Location, strings.Join(hack.Tags, ", ")))
	}

	hackathonsContext := contextBuilder.String()
	if hackathonsContext == "" {
		hackathonsContext = "No specific hackathons found for this exact query."
	}

	aiCtx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	insight, err := h.AIService.GenerateSearchInsights(aiCtx, profileContext, query, hackathonsContext)
	if err != nil {
		fmt.Printf("AI overview generation error: %v\n", err)
		http.Error(w, "Failed to generate AI overview", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"overview": insight,
	})

}
