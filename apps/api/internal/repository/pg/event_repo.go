package pg

import (
	"database/sql"
	"fmt"
	"strconv"
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
	query := `SELECT id, title, COALESCE(host, 'Unknown Host'), COALESCE(location, 'TBA'), COALESCE(prize_usd, 0.0), start_date, end_date, COALESCE(apply_url, ''), COUNT(*) OVER() FROM hackathons WHERE 1=1`

	var conditions []string
	var args []any
	argID := 1

	if filters.Query != "" {
		conditions = append(conditions, fmt.Sprintf("search_vector @@ websearch_to_tsquery('english', $%d)", argID))
		args = append(args, filters.Query)
		argID++
	}

	if filters.Status != "" {
		switch filters.Status {
		case "upcoming":
			conditions = append(conditions, "start_date > NOW()")
		case "ongoing":
			conditions = append(conditions, "start_date <= NOW() AND end_date >= NOW()")
		case "past":
			conditions = append(conditions, "end_date < NOW()")
		}
	} else {
		// The Safety Net: never return expired hackathons unless explicitly requesting "past"
		conditions = append(conditions, "end_date >= CURRENT_DATE")
	}

	if filters.Location != "" {
		switch filters.Location {
		case "online":
			conditions = append(conditions, "(location ILIKE '%remote%' OR location ILIKE '%online%')")
		case "in-person":
			conditions = append(conditions, "(location NOT ILIKE '%remote%' AND location NOT ILIKE '%online%')")
		}
	}

	if filters.PrizeRange != "" {
		if strings.HasSuffix(filters.PrizeRange, "+") {
			minStr := strings.TrimSuffix(filters.PrizeRange, "+")
			if minPrize, err := strconv.ParseFloat(minStr, 64); err == nil {
				conditions = append(conditions, fmt.Sprintf("prize_usd >= $%d", argID))
				args = append(args, minPrize)
				argID++
			}
		} else {
			parts := strings.Split(filters.PrizeRange, "-")
			if len(parts) == 2 {
				minPrize, err1 := strconv.ParseFloat(parts[0], 64)
				maxPrize, err2 := strconv.ParseFloat(parts[1], 64)

				if err1 == nil && err2 == nil {
					conditions = append(conditions, fmt.Sprintf("prize_usd >= $%d AND prize_usd <= $%d", argID, argID+1))
					args = append(args, minPrize, maxPrize)
					argID += 2
				}
			}
		}
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

	offset := (filters.Page - 1) * filters.Limit

	//append ORDER BY, LIMIT and OFFSET
	query += fmt.Sprintf(" ORDER BY start_date DESC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, filters.Limit, offset)

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

func (h *PostgresEventRepo) DeleteExpiredHackathons() (int64, error) {
	query := `DELETE FROM hackathons WHERE end_date < CURRENT_DATE`

	res, err := h.db.Exec(query)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	fmt.Printf("Deleted %d expired hackathons\n", rowsAffected)
	return rowsAffected, nil
}
