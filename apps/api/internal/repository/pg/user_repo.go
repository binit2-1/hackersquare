package pg

import (
	"database/sql"
	"fmt"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)


type PostgresUserRepo struct{
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
	
	
	
	err =  h.db.QueryRow(
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

	query := `SELECT id, name, email from users WHERE id = $1`

	err := h.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}


	return &user, nil

}