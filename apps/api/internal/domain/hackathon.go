package domain

import (
	"context"
	"time"
)

type Hackathon struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Host      string    `json:"host"`
	Location  string    `json:"location"`
	PrizeUSD  *float64  `json:"prize_usd"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	ApplyURL  string    `json:"apply_url"`
	Tags      []string  `json:"tags"`
}

type SearchFilters struct {
	Query      string
	Location   string
	Status     string
	PrizeRange string
	Page       int
	Limit      int
}

type HackathonRepository interface {
	SearchHackathons(filters SearchFilters) ([]Hackathon, int, error)
	NearbyHackathons(city, country string, page, limit int) ([]Hackathon, int, error)
	DeleteExpiredHackathons() (int64, error)
	GetUserRecommendations(tags []string, city, state, country string, limit int) ([]Hackathon, error)
	GetMatchingChats(ctx context.Context, hackLocation string, hackTags []string) ([]string, error)
}
