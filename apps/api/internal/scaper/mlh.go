package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
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
	c.SetRequestTimeout(60 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	})

	// Select event items
	c.OnHTML("a[itemtype='https://schema.org/Event']", func(e *colly.HTMLElement) {
		name := e.ChildText("[itemprop='name']")
		applyURL := e.Attr("href")
		location := e.ChildText("p.font-bold.text-gray-700")
		dateText := e.ChildText("p.text-gray-600")

		start, end, err := ParseMLHDates(dateText)
		if err != nil {
			return
		}

		if end.Before(time.Now()) {
			return
		}

		// Fire the Deep Scraper to dig for the prize
		prizeUSD := DeepScrape(applyURL)

		// Pass the dynamic prizeUSD to the database
		SaveToDB(db, name, location, applyURL, prizeUSD, start, end)
	})

	c.Limit(&colly.LimitRule{DomainGlob: "*mlh.io*", Delay: 2 * time.Second})

	return c.Visit("https://mlh.io/seasons/2026/events")
}

func DeepScrape(applyURL string) float64 {
	// 1. THE URL CLEANER: Strip MLH's massive UTM tracking parameters
	cleanURL := strings.Split(applyURL, "?")[0]
	utils.Debug("[DeepScrape] Investigating clean URL: %s", cleanURL)

	// 2. THE GRACE PERIOD: Give student servers 15 seconds to wake up from cold starts
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", cleanURL, nil)
	if err != nil {
		return 0.0
	}

	// 3. THE DISGUISE: Pretend to be a modern Chrome browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")

	resp, err := client.Do(req)
	if err != nil {
		utils.Error("[DeepScrape] Network error fetching %s: %v", cleanURL, err)
		return 0.0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.Error("[DeepScrape] Blocked/Failed! HTTP Status %d for %s", resp.StatusCode, cleanURL)
		return 0.0
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0.0
	}
	htmlStr := string(bodyBytes)

	// 4. THE PROXIMITY REGEX: Find the prizes and filter out false positives
	re := regexp.MustCompile(`\$([0-9]{1,3}(?:,[0-9]{3})*(?:[kK])?)`)
	matches := re.FindAllStringSubmatchIndex(htmlStr, -1)

	var maxPrize float64
	for _, match := range matches {
		startIdx := match[0]
		endIdx := match[1]

		ctxStart := startIdx - 60
		if ctxStart < 0 {
			ctxStart = 0
		}
		ctxEnd := endIdx + 60
		if ctxEnd > len(htmlStr) {
			ctxEnd = len(htmlStr)
		}
		contextText := strings.ToLower(htmlStr[ctxStart:ctxEnd])

		// Ignore false positives like sponsors, fees, or travel stipends
		if strings.Contains(contextText, "sponsor") ||
			strings.Contains(contextText, "stipend") ||
			strings.Contains(contextText, "travel") ||
			strings.Contains(contextText, "fee") ||
			strings.Contains(contextText, "cost") {
			continue
		}

		valStr := htmlStr[match[2]:match[3]]
		valStr = strings.ReplaceAll(valStr, ",", "")
		valStr = strings.ToLower(valStr)

		multiplier := 1.0
		if strings.Contains(valStr, "k") {
			multiplier = 1000.0
			valStr = strings.ReplaceAll(valStr, "k", "")
		}

		if val, err := strconv.ParseFloat(valStr, 64); err == nil {
			calculatedPrize := val * multiplier
			if calculatedPrize > 250000 { // Ignore AWS/Cloud credit bundles
				continue
			}
			if calculatedPrize > maxPrize {
				maxPrize = calculatedPrize
			}
		}
	}

	return maxPrize
}

// ParseMLHDates converts MLH's string dates (e.g., "FEB 13 - 15") to time.Time
func ParseMLHDates(dateStr string) (time.Time, time.Time, error) {
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
		return time.Time{}, time.Time{}, fmt.Errorf("unrecognized date format: %s", dateStr)
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

	return startDate, endDate, nil
}

// SaveToDB uses PostgreSQL ON CONFLICT to efficiently Upsert the hackathon
func SaveToDB(db *sql.DB, title, loc, applyURL string, prizeUSD float64, start, end time.Time) {
	query := `
		INSERT INTO hackathons (title, host, location, prize_usd, start_date, end_date, apply_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (apply_url) 
		DO UPDATE SET 
			title = EXCLUDED.title,
			location = EXCLUDED.location,
			prize_usd = EXCLUDED.prize_usd,
			start_date = EXCLUDED.start_date,
			end_date = EXCLUDED.end_date,
			updated_at = NOW();
	`
	_, err := db.ExecContext(context.Background(), query, title, "MLH", loc, prizeUSD, start, end, applyURL)
	if err != nil {
		utils.Error("Database upsert failed for %s: %v", title, err)
	} else {
		utils.Debug("Saved/Updated: %s | Prize: $%.2f", title, prizeUSD)
	}
}
