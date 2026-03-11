package pg

import (
	"database/sql"
	"fmt"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgreUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		db: db,
	}
}

func (h *PostgresUserRepo) CreateUser(user *domain.User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err = h.db.QueryRow(
		query,
		user.Name,
		user.Email,
		string(hashedPassword),
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.PasswordHash = string(hashedPassword)

	return nil
}

func (h *PostgresUserRepo) GetUserByEmail(email string) (*domain.User, error) {

	var user domain.User

	query := `SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := h.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (h *PostgresUserRepo) GetUserByID(id string) (*domain.User, error) {
	var user domain.User

	query := `SELECT id, name, email, COALESCE(headline, ''), COALESCE(location, ''), COALESCE(github_handle, ''), COALESCE(website_url, ''), COALESCE(linkedin_url, ''), COALESCE(twitter_url, '') FROM users WHERE id = $1`

	err := h.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Headline,
		&user.Location,
		&user.GithubHandle,
		&user.WebsiteURL,
		&user.LinkedinURL,
		&user.TwitterURL,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (h *PostgresUserRepo) UpdateUserProfile(userID string, data domain.ProfileUpdateRequest) error {
	query := `UPDATE users SET headline = $1, location = $2, website_url = $3, linkedin_url = $4, twitter_url = $5, updated_at = NOW() WHERE id = $6`

	result, err := h.db.Exec(query, data.Headline, data.Location, data.WebsiteURL, data.LinkedinURL, data.TwitterURL, userID)
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (h *PostgresUserRepo) LinkGithubHandle(userID string, githubHandle string) error {
	query := `UPDATE users SET github_handle = $1, updated_at = NOW() WHERE id = $2`
	_, err := h.db.Exec(query, githubHandle, userID)
	return err
}
