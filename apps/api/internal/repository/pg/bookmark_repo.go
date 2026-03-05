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

func (h *PostgresBookmarkRepo) GetBookmarksByUser(userID string) ([]domain.Hackathon, error) {

	query := `
		SELECT 
			h.id, 
			h.title, 
			COALESCE(h.host, 'Unknown Host'), 
			COALESCE(h.location, 'TBA'), 
			COALESCE(h.prize_usd, 0.0), 
			h.start_date, 
			h.end_date, 
			COALESCE(h.apply_url, '')
		FROM hackathons h
		INNER JOIN bookmarks b ON h.id = b.hackathon_id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`

	rows, err  := h.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []domain.Hackathon
	for rows.Next(){
		var bookmark domain.Hackathon
		err := rows.Scan(
			&bookmark.ID,
			&bookmark.Title,
			&bookmark.Host,
			&bookmark.Location,
			&bookmark.PrizeUSD,
			&bookmark.StartDate,
			&bookmark.EndDate,
			&bookmark.ApplyURL,
		)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bookmark)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if bookmarks == nil {
		bookmarks = []domain.Hackathon{}
	}

	
	return bookmarks, nil
}

