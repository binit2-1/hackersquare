package hackathon

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/binit2-1/hackersquare/apps/api/internal/database"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
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
	query := `SELECT id, title, host, location, prize, "startDate", "endDate", "applyUrl", tags FROM hackathons`
	rows, err := h.DB.Pool.Query(r.Context(), query)
	if err != nil{
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	//make() ensures if db is empty, we return an empty array instead of null in JSON response
	hackathons := make([]Hackathon, 0)

	for rows.Next(){
		var hack Hackathon
		
		err := rows.Scan(
			&hack.ID,
			&hack.Title,
			&hack.Host,
			&hack.Location,
			&hack.Prize,
			&hack.StartDate,
			&hack.EndDate,
			&hack.ApplyURL,
			&hack.Tags,
		)
		if err != nil{
			http.Error(w, "Failed to scan database row", http.StatusInternalServerError)
			return
		}

		hackathons = append(hackathons, hack)
	}

	if rows.Err() != nil{
		http.Error(w, "Error iterating over the rows", http.StatusInternalServerError)
		return
	}


	//res.json() equivalent in Go set Header to json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)


	err = json.NewEncoder(w).Encode(hackathons)
	if err!= nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

//GET /hackathons/{id} endpoint handler 
func(h *Handler) GetHackathonByID(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r) // Extracts path variables from the request, in this case, the 'id' from the URL
	id := vars["id"]
	
	
	query := `SELECT id, title, host, location, prize, "startDate", "endDate", "applyUrl", tags FROM hackathons WHERE id = $1`

	var hack Hackathon

	err:= h.DB.Pool.QueryRow(r.Context(), query, id).Scan(
		&hack.ID,
		&hack.Title,
		&hack.Host,
		&hack.Location,
		&hack.Prize,
		&hack.StartDate,
		&hack.EndDate,
		&hack.ApplyURL,
		&hack.Tags,
	)

	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			http.Error(w, "Hackathon not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Failed to fetch the data", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(hack)
	if err!= nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}