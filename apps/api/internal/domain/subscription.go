package domain

import "time"

type ChannelSubscription struct {
	ID        string    `json:"id"`
	Platform  string    `json:"platform"`
	ChatID    string    `json:"chat_id"`
	TechTags  []string  `json:"tech_tags"`
	Country   string    `json:"country"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
