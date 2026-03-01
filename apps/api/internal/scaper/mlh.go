package scraper

import (
	"context"
	"database/sql"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/utils"
	"github.com/gocolly/colly/v2"
)

// RunMLHScraper fetches hackathons from the MLH season page and saves them
func RunMLHScraper(db *sql.DB) error {
	utils.Info("[MLH] Starting crawl")
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Select event items
	c.OnHTML("a[itemtype='https://schema.org/Event']", func(e *colly.HTMLElement) {
		name := e.ChildText("[itemprop='name']")
		applyURL := e.Attr("href")
		location := e.ChildText("p.font-bold.text-gray-700")
		dateText := e.ChildText("p.text-gray-600")

		utils.Debug("Scraped: %s | Loc: %s | URL: %s", name, location, applyURL)

		// Parse dates
		start, end := ParseMLHDates(dateText)

		// Directly pass to the Upsert function (No description needed anymore)
		SaveToDB(db, name, location, applyURL, start, end)
	})

	c.Limit(&colly.LimitRule{DomainGlob: "*mlh.io*", Delay: 2 * time.Second})

	return c.Visit("https://mlh.io/seasons/2026/events")
}

// ParseMLHDates converts MLH's string dates (e.g., "FEB 13 - 15") to time.Time
func ParseMLHDates(dateStr string) (time.Time, time.Time) {
	year := 2026 // Target season

	months := map[string]time.Month{
		"JAN": time.January, "FEB": time.February, "MAR": time.March,
		"APR": time.April, "MAY": time.May, "JUN": time.June,
		"JUL": time.July, "AUG": time.August, "SEP": time.September,
		"OCT": time.October, "NOV": time.November, "DEC": time.December,
	}

	re := regexp.MustCompile(`([A-Z]{3})\s+(\d+)(?:\s*-\s*(\d+))?`)
	matches := re.FindStringSubmatch(strings.ToUpper(dateStr))

	if len(matches) < 3 {
		return time.Now(), time.Now().Add(48 * time.Hour)
	}

	monthStr := matches[1]
	startDay, _ := strconv.Atoi(matches[2])
	month := months[monthStr]

	startDate := time.Date(year, month, startDay, 0, 0, 0, 0, time.UTC)

	var endDate time.Time
	if matches[3] != "" {
		endDay, _ := strconv.Atoi(matches[3])
		endDate = time.Date(year, month, endDay, 23, 59, 59, 0, time.UTC)
	} else {
		endDate = startDate.Add(24 * time.Hour)
	}

	return startDate, endDate
}

// SaveToDB uses PostgreSQL ON CONFLICT to efficiently Upsert the hackathon
func SaveToDB(db *sql.DB, title, loc, applyURL string, start, end time.Time) {
	query := `
		INSERT INTO hackathons (title, host, location, prize_usd, start_date, end_date, apply_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (apply_url) 
		DO UPDATE SET 
			title = EXCLUDED.title,
			location = EXCLUDED.location,
			start_date = EXCLUDED.start_date,
			end_date = EXCLUDED.end_date,
			updated_at = NOW();
	`
	
	// MLH doesn't reliably list exact prize amounts on the index page, so we default to 0.0
	_, err := db.ExecContext(context.Background(), query, title, "MLH", loc, 0.0, start, end, applyURL)
	if err != nil {
		utils.Error("Database upsert failed for %s: %v", title, err)
	} else {
		utils.Debug("Saved/Updated: %s", title)
	}
}