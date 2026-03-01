package domain

import "time"

type Hackathon struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Host      string   `json:"host"`
	Location  string   `json:"location"`
	Prize     string   `json:"prize"`
	PrizeUSD  *float64   `json:"prizeUSD"`
	StartDate time.Time   `json:"startDate"`
	EndDate   time.Time   `json:"endDate"`
	ApplyURL  string   `json:"applyURL"`
	Tags      []string `json:"tags"`
}

type HackathonRepository interface {

}