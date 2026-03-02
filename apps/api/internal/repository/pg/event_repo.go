package pg

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
)

type PostgresEventRepo struct {
	db *sql.DB
}

func NewPostgreEventRepo(db *sql.DB) domain.HackathonRepository {
	return &PostgresEventRepo{
		db: db,
	}
}

func (h *PostgresEventRepo) SearchHackathons(filters domain.SearchFilters) ([]domain.Hackathon, int, error) {
	query := `SELECT id, title, host, location, prize_usd, start_date, end_date, apply_url, COUNT(*) OVER() FROM hackathons WHERE 1=1`
	var conditions []string
	var args []any
	argID := 1

	if filters.Query != "" {
		conditions = append(conditions, fmt.Sprintf("search_vector @@ websearch_to_tsquery('english', $%d)", argID))
		args = append(args, filters.Query)
		argID++
	}

	if filters.Location != "" {
		conditions = append(conditions, fmt.Sprintf("location ILIKE $%d", argID))
		args = append(args, "%"+filters.Location+"%")
		argID++
	}

	if filters.MinPrize > 0 {
		conditions = append(conditions, fmt.Sprintf("prize_usd >= $%d", argID))
		args = append(args, filters.MinPrize)
		argID++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	limit := filters.Limit
	if limit <= 0 {
		limit = 20
	}

	page := filters.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	//append ORDER BY, LIMIT and OFFSET
	query += fmt.Sprintf(" ORDER BY start_date DESC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var hackathons []domain.Hackathon
	totalCount := 0

	for rows.Next() {
		var event domain.Hackathon
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Host,
			&event.Location,
			&event.PrizeUSD,
			&event.StartDate,
			&event.EndDate,
			&event.ApplyURL,
			&totalCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}

		hackathons = append(hackathons, event)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return hackathons, totalCount, nil
}
