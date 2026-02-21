package hackathon

import (
	"encoding/json"
	"net/http"
	"time"
)

//GET /hackathons endpoint handler
func GetHackathons(w http.ResponseWriter, r *http.Request){
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