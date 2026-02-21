package hackathon

import "time"

type Hackathon struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Host      string    `json:"host"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Location  string    `json:"location"`
	Prize     string    `json:"prize"`
	Tags      []string  `json:"tags"`
	ApplyURL  string    `json:"apply_url"`
}
