package hackathon

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/database"
	"github.com/gorilla/mux"
)

// Handler holds the database connection so our routes can use it
type Handler struct{
	DB *database.Service
}

//Constructor for Handler
func NewHandler(db *database.Service) *Handler{
	return &Handler{DB: db}
}

//GET /hackathons endpoint handler
func(h *Handler) GetHackathons(w http.ResponseWriter, r *http.Request){
	dummyData := []Hackathon{
		{
			ID:        "1",
			Title:     "Global AI Hackathon",
			Host:      "TechCorp",
			Location:  "Online",
			Prize:     "$10,000",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 3), // Adds 3 days
			ApplyURL:  "https://example.com/apply",
			Tags:      []string{"AI", "Web3"},
		},
	}

	//res.json() equivalent in Go set Header to json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)


	err := json.NewEncoder(w).Encode(dummyData)
	if err!= nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

//GET /hackathons/{id} endpoint handler 
func(h *Handler) GetHackathonByID(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r) // Extracts path variables from the request, in this case, the 'id' from the URL
	id := vars["id"]
	
	// For demonstration, we return a dummy hackathon with the requested ID
	dummyHackathon := Hackathon{
		ID:        id,
		Title:     "Global AI Hackathon",
		Host:      "TechCorp",
		Location:  "Online",
		Prize:     "$10,000",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 3), // Adds 3 days
		ApplyURL:  "https://example.com/apply",
		Tags:      []string{"AI", "Web3"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(dummyHackathon)
	if err!= nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}