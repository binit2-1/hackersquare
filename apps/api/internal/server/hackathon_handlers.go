package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
)

type HackathonHandler struct {
	Repo domain.HackathonRepository
}

func NewHackathonHandler(repo domain.HackathonRepository) *HackathonHandler {
	return &HackathonHandler{
		Repo: repo,
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
		fmt.Printf("❌ Database Search Error: %v\n", err)

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
		fmt.Printf("❌ Nearby Query Error: %v\n", err)
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
