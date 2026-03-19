package pg

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
	"github.com/lib/pq"
)

type PostgresEventRepo struct {
	db *sql.DB
}

var searchStopWords = map[string]struct{}{
	"a":            {},
	"an":           {},
	"the":          {},
	"for":          {},
	"to":           {},
	"of":           {},
	"in":           {},
	"on":           {},
	"at":           {},
	"by":           {},
	"with":         {},
	"about":        {},
	"and":          {},
	"or":           {},
	"is":           {},
	"are":          {},
	"nearby":       {},
	"me":           {},
	"hackathon":    {},
	"hackathons":   {},
	"hack":         {},
	"hacks":        {},
	"event":        {},
	"events":       {},
	"competition":  {},
	"competitions": {},
	"challenge":    {},
	"challenges":   {},
}

func NormalizeSearchQuery(rawQuery string) string {
	trimmed := strings.TrimSpace(rawQuery)
	if trimmed == "" {
		return ""
	}

	words := strings.Fields(strings.ToLower(trimmed))
	filtered := make([]string, 0, len(words))
	for _, word := range words {
		cleaned := strings.Trim(word, ".,!?:;()[]{}\"'`")
		if cleaned == "" {
			continue
		}

		if _, isStopWord := searchStopWords[cleaned]; isStopWord {
			continue
		}

		filtered = append(filtered, cleaned)
	}

	if len(filtered) == 0 {
		return ""
	}

	return strings.Join(filtered, " ")
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
		normalizedQuery := NormalizeSearchQuery(filters.Query)
		if normalizedQuery != "" {
			conditions = append(conditions, fmt.Sprintf("search_vector @@ websearch_to_tsquery('english', $%d)", argID))
			args = append(args, normalizedQuery)
			argID++
		}
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

func (h *PostgresEventRepo) NearbyHackathons(city, country string, page, limit int) ([]domain.Hackathon, int, error) {
	query := `SELECT id, title, COALESCE(host, 'Unknown Host'), COALESCE(location, 'TBA'), COALESCE(prize_usd, 0.0), start_date, end_date, COALESCE(apply_url, ''), COUNT(*) OVER() FROM hackathons WHERE end_date >= CURRENT_DATE`

	var geoConditions []string
	var args []any
	argID := 1

	if city != "" {
		geoConditions = append(geoConditions, fmt.Sprintf("location ILIKE $%d", argID))
		args = append(args, "%"+city+"%")
		argID++
	}

	if country != "" {
		geoConditions = append(geoConditions, fmt.Sprintf("location ILIKE $%d", argID))
		args = append(args, "%"+country+"%")
		argID++
	}

	if len(geoConditions) > 0 {
		query += " AND (" + strings.Join(geoConditions, " OR ") + ")"
	}

	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	query += fmt.Sprintf(" ORDER BY start_date DESC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("nearby query failed: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan nearby row: %w", err)
		}

		hackathons = append(hackathons, event)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("nearby row iteration error: %w", err)
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

func (h *PostgresEventRepo) GetUserRecommendations(tags []string, city, state, country string, limit int) ([]domain.Hackathon, error) {
	tagQuery := strings.Join(tags, " ")
	if tagQuery == "" {
		tagQuery = "hackathon"
	}

	query := `
    SELECT id, title, COALESCE(host, 'Unknown Host'), COALESCE(location, 'TBA'), COALESCE(prize_usd, 0.0), start_date, end_date, COALESCE(apply_url, '')
    FROM hackathons
    WHERE start_date >= CURRENT_DATE + INTERVAL '5 days'
    `

	var args []any
	argID := 1
	var scoreCases []string

	// Tier 4: Exact City Match (Highest Priority)
	if city != "" {
		scoreCases = append(scoreCases, fmt.Sprintf("WHEN location ILIKE $%d THEN 4", argID))
		args = append(args, "%"+city+"%")
		argID++
	}
	// Tier 3: State Match
	if state != "" {
		scoreCases = append(scoreCases, fmt.Sprintf("WHEN location ILIKE $%d THEN 3", argID))
		args = append(args, "%"+state+"%")
		argID++
	}
	// Tier 2: Country Match
	if country != "" {
		scoreCases = append(scoreCases, fmt.Sprintf("WHEN location ILIKE $%d THEN 2", argID))
		args = append(args, "%"+country+"%")
		argID++
	}

	// Tier 1: Online/Remote (Safe fallback before showing random global events)
	scoreCases = append(scoreCases, "WHEN location ILIKE '%online%' OR location ILIKE '%remote%' THEN 1")

	// Compile the scoring logic
	locationScoreExpr := "CASE " + strings.Join(scoreCases, " ") + " ELSE 0 END"

	// Add the tags for the ts_rank
	args = append(args, tagQuery)
	rankArg := argID
	argID++

	// The Magic: Strictly order by Geographic Tier FIRST, then Tech Stack relevance SECOND
	query += fmt.Sprintf(" ORDER BY (%s) DESC, ts_rank(search_vector, websearch_to_tsquery('english', $%d)) DESC LIMIT $%d", locationScoreExpr, rankArg, argID)
	args = append(args, limit)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("recommendations query failed: %w", err)
	}
	defer rows.Close()

	var hackathons []domain.Hackathon
	for rows.Next() {
		var event domain.Hackathon
		err := rows.Scan(&event.ID, &event.Title, &event.Host, &event.Location, &event.PrizeUSD, &event.StartDate, &event.EndDate, &event.ApplyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recommendation row: %w", err)
		}
		hackathons = append(hackathons, event)
	}

	return hackathons, nil
}

func (r *PostgresEventRepo) UpsertSubscription(ctx context.Context, platform, chatID string, tags []string, country string) error {
	query := `
        INSERT INTO channel_subscriptions (platform, chat_id, tech_tags, country, is_active)
        VALUES ($1, $2, $3, $4, true)
        ON CONFLICT (chat_id) 
        DO UPDATE SET 
            tech_tags = EXCLUDED.tech_tags,
            country = EXCLUDED.country,
            is_active = true,
            platform = EXCLUDED.platform;
    `

	_, err := r.db.ExecContext(ctx, query, platform, chatID, pq.Array(tags), country)
	if err != nil {
		return fmt.Errorf("failed to upsert subscription: %w", err)
	}
	return nil
}

func (r *PostgresEventRepo) GetMatchingChats(ctx context.Context, hackLocation string, hackTags []string) ([]string, error) {
	if len(hackTags) == 0 {
		hackTags = []string{"hackathon"}
	}

	query := `
		SELECT chat_id 
		FROM channel_subscriptions 
		WHERE is_active = true
		AND (country = '' OR $1 ILIKE '%' || country || '%')
		AND (cardinality(tech_tags) = 0 OR tech_tags && $2)
	`

	rows, err := r.db.QueryContext(ctx, query, hackLocation, pq.Array(hackTags))
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	var chatIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			chatIDs = append(chatIDs, id)
		}
	}
	return chatIDs, nil

}

func (r *PostgresEventRepo) GetNewHackathonsSince(ctx context.Context, since time.Time) ([]domain.Hackathon, error) {
	query := `
		SELECT id, title, COALESCE(host, 'Unknown Host'), COALESCE(location, 'TBA'), COALESCE(prize_usd, 0.0), start_date, end_date, COALESCE(apply_url, '')
		FROM hackathons
		WHERE created_at > $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, since)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch new hackathons: %w", err)
	}
	defer rows.Close()

	var hackathons []domain.Hackathon
	for rows.Next() {
		var event domain.Hackathon
		err := rows.Scan(&event.ID, &event.Title, &event.Host, &event.Location, &event.PrizeUSD, &event.StartDate, &event.EndDate, &event.ApplyURL)
		if err == nil {
			hackathons = append(hackathons, event)
		}
	}
	return hackathons, nil

}
