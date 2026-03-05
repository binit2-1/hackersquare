package pg

import (
	"database/sql"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
)

type PostgresBookmarkRepo struct {
	db *sql.DB
}

func NewPostgresBookmarkRepo(db *sql.DB) *PostgresBookmarkRepo {
	return &PostgresBookmarkRepo{
		db: db,
	}
}

func (h *PostgresBookmarkRepo) AddBookmark(userID, hackathonID string) error {
	query := `INSERT INTO bookmarks (user_id, hackathon_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := h.db.Exec(
		query,
		userID,
		hackathonID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (h *PostgresBookmarkRepo) RemoveBookmark(userID, hackathonID string) error {

	query := `DELETE FROM bookmarks WHERE user_id = $1 AND hackathon_id = $2`

	_, err := h.db.Exec(
		query,
		userID,
		hackathonID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (h *PostgresBookmarkRepo) GetBookmarksByUser(userID string) ([]domain.Bookmark, error) {

	query := `SELECT user_id, hackathon_id, created_at FROM bookmarks WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err  := h.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []domain.Bookmark
	for rows.Next(){
		var bookmark domain.Bookmark
		err := rows.Scan(
			&bookmark.UserID,
			&bookmark.HackathonID,
			&bookmark.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}
	return bookmarks, nil
}

