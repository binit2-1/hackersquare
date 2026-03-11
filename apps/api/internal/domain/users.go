package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Headline     string    `json:"headline"`
	Location     string    `json:"location"`
	GithubHandle string    `json:"github_handle"`
	WebsiteURL   string    `json:"website_url"`
	LinkedinURL  string    `json:"linkedin_url"`
	TwitterURL   string    `json:"twitter_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProfileUpdateRequest struct {
	Headline    string `json:"headline"`
	Location    string `json:"location"`
	WebsiteURL  string `json:"website_url"`
	LinkedinURL string `json:"linkedin_url"`
	TwitterURL  string `json:"twitter_url"`
}

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	UpdateUserProfile(userID string, data ProfileUpdateRequest) error
}
