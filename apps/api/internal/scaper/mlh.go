package scraper

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/database"
	"github.com/binit2-1/hackersquare/apps/api/internal/utils"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

type MLHRawEvent struct {
	Name        string
	DetailURL   string
	Location    string
	StartDate   string
	EndDate     string
	Description string
}

func RunMLHScraper(db *database.Service) error {
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

		// Get raw description
		desc := ScrapeExternalDescription(applyURL)

		// Parse dates
		start, end := ParseMLHDates(dateText)

		// Save (deduplicated)
		SaveToDB(db, name, location, applyURL, desc, start, end)
	})

	c.Limit(&colly.LimitRule{DomainGlob: "*mlh.io*", Delay: 2 * time.Second})

	return c.Visit("https://mlh.io/seasons/2026/events")
}

// ScrapeExternalDescription visits the hackathon's own site to get raw text
func ScrapeExternalDescription(url string) string {
	var bodyText string
	d := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0..."),
	)
	// We set a 5-second timeout so one slow website doesn't hang your whole scraper
	d.SetRequestTimeout(5 * time.Second)

	d.OnHTML("body", func(e *colly.HTMLElement) {
		// Just grab the first 1000 characters of text to keep things light
		text := strings.TrimSpace(e.Text)
		if len(text) > 1000 {
			bodyText = text[:1000]
		} else {
			bodyText = text
		}
	})

	d.Visit(url)
	return bodyText
}

func ParseMLHDates(dateStr string) (time.Time, time.Time) {
	year := 2026 // Hardcoded based on the season you are scraping

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

func SaveToDB(db *database.Service, title, loc, url, desc string, start, end time.Time) {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM hackathons WHERE "applyUrl" = $1)`
	db.Pool.QueryRow(context.Background(), checkQuery, url).Scan(&exists)

	if exists {
		utils.Debug("Skipping %s: Already exists", title)
		return
	}

	query := `
		INSERT INTO hackathons (id, title, host, location, prize, "startDate", "endDate", "applyUrl", tags, "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	newID := uuid.New().String()
	_, err := db.Pool.Exec(context.Background(), query, newID, title, "MLH", loc, "TBA", start, end, url, []string{"MLH", "Student"}, time.Now())
	if err == nil {
		utils.Debug("Saved: %s", title)
	}
}
