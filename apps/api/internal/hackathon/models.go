package hackathon

import "time"

type Hackathon struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Host      string    `json:"host"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Location  string    `json:"location"`
	Prize     string    `json:"prize"`
	Tags      []string  `json:"tags"`
	ApplyURL  string    `json:"applyUrl"`
}
