package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/database"
	"github.com/binit2-1/hackersquare/apps/api/internal/utils"
	"github.com/google/uuid"
)

// DevfolioIndexResponse represents the JSON structure from the hackathons list API.
type DevfolioIndexResponse struct {
	PageProps struct {
		DehydratedState struct {
			Queries []struct {
				State struct {
					Data struct {
						OpenHackathons []struct {
							Slug string `json:"slug"`
						} `json:"open_hackathons"`
					} `json:"data"`
				} `json:"state"`
			} `json:"queries"`
		} `json:"dehydratedState"`
	} `json:"pageProps"`
}

// DevfolioDetailResponse represents the hackathon detail data embedded in Next.js __NEXT_DATA__
type DevfolioDetailResponse struct {
	Props struct {
		PageProps struct {
			AggregatePrizeValue    float64 `json:"aggregatePrizeValue"`
			AggregatePrizeCurrency string  `json:"aggregatePrizeCurrency"`
			Hackathon              struct {
				Name     string `json:"name"`
				Slug     string `json:"slug"`
				Desc     string `json:"desc"`
				StartsAt string `json:"starts_at"`
				EndsAt   string `json:"ends_at"`
				Location string `json:"location"`
				Prizes   []struct {
					Title string `json:"title"`
					Value int    `json:"amount"`
				}
			} `json:"hackathon"`
		} `json:"pageProps"`
	} `json:"props"`
}

// DevfolioPrizeJSON represents the structure of the prize data from the Devfolio prizes API endpoint
type DevfolioPrizeJSON struct {
	PageProps struct {
		AggregatePrizeValue    float64 `json:"aggregatePrizeValue"`
		AggregatePrizeCurrency string  `json:"aggregatePrizeCurrency"`
	} `json:"pageProps"`
}

// RunDevfolioScraper fetches hackathons from Devfolio and saves them to the database
func RunDevfolioScraper(db *database.Service) error {
	utils.Info("[Devfolio] Booting up scraper...")
	rate := utils.GetExchangeRate()

	// Extract build ID
	buildID, err := getBuildID()
	if err != nil {
		return fmt.Errorf("failed to get build ID: %w", err)
	}
	utils.Info("Secured Build ID: %s", buildID)

	// Fetch index of hackathons
	indexURL := fmt.Sprintf("https://devfolio.co/_next/data/%s/hackathons.json", buildID)
	indexData, err := fetchIndex(indexURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Devfolio index: %w", err)
	}

	if len(indexData.PageProps.DehydratedState.Queries) == 0 {
		return fmt.Errorf("no queries found in Devfolio index data")
	}

	hackathons := indexData.PageProps.DehydratedState.Queries[0].State.Data.OpenHackathons
	utils.Info("Found %d open hackathons", len(hackathons))

	// Loop through index and fetch details
	for i, hack := range hackathons {
		utils.Debug("[%d/%d] Extracting data for: %s", i+1, len(hackathons), hack.Slug)

		// Construct subdomain URL
		subdomainURL := fmt.Sprintf("https://%s.devfolio.co/", hack.Slug)

		detailData, err := fetchDetail(subdomainURL)
		if err != nil {
			utils.Error("Error fetching %s: %v", hack.Slug, err)
			continue
		}

		descLength := len(detailData.Props.PageProps.Hackathon.Desc)
		hackName := detailData.Props.PageProps.Hackathon.Name
		utils.Debug("Extracted %d bytes of markdown description for '%s'", descLength, hackName)

		actualPrize, err := fetchPrize(hack.Slug, buildID)
		if err != nil {
			utils.Debug("Prize fetch failed for %s, defaulting to TBA", hack.Slug)
		}

		title, host, loc, prize, prizeUSD, start, end, applyURL, err := DevfolioAdapter(detailData, actualPrize, rate)
		if err != nil {
			utils.Error("Adapter failed for %s: %v", hack.Slug, err)
			continue
		}

		utils.Debug("Transformed: %s | Host: %s | Loc: %s | Prize: %s | Prize USD: %.2f | Starts: %s | Ends: %s | URL: %s", title, host, loc, prize, prizeUSD, start.Format("Jan 02, 2006"), end.Format("Jan 02, 2006"), applyURL)

		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM hackathons WHERE "applyUrl" = $1)`

		err = db.Pool.QueryRow(context.Background(), checkQuery, applyURL).Scan(&exists)
		if err != nil {
			utils.Error("Database check failed for %s: %v", title, err)
			continue
		}

		// Skip if already in database
		if exists {
			utils.Debug("Skipping '%s': Already exists in database", title)
			continue
		}

		// Insert new hackathon
		newID := uuid.New().String()
		tags := []string{"Devfolio", "Upcoming"}

		insertQuery := `
			INSERT INTO hackathons (id, title, host, location, prize, "prizeUSD", "startDate", "endDate", "applyUrl", tags, "updatedAt")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

		_, err = db.Pool.Exec(context.Background(), insertQuery, newID, title, host, loc, prize, prizeUSD, start, end, applyURL, tags, time.Now())
		if err != nil {
			utils.Error("Database insert failed for %s: %v", title, err)
			continue
		}

		utils.Debug("Successfully saved '%s' to PostgreSQL", title)

		// Rate limiting
		time.Sleep(1 * time.Second)

	}

	return nil
}

func getBuildID() (string, error) {
	resp, err := http.Get("https://devfolio.co/hackathons/open")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"buildId":"([^"]+)"`)
	matches := re.FindSubmatch(htmlBytes)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find buildId")
	}
	return string(matches[1]), nil
}

func fetchIndex(url string) (*DevfolioIndexResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var data DevfolioIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func fetchDetail(url string) (*DevfolioDetailResponse, error) {
	// Fetch the hackathon page HTML
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract JSON from Next.js __NEXT_DATA__ script tag
	re := regexp.MustCompile(`(?s)<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)
	matches := re.FindSubmatch(htmlBytes)
	if len(matches) < 2 {
		return nil, fmt.Errorf("could not find __NEXT_DATA__ JSON in HTML")
	}

	rawJSON := matches[1]

	// Parse JSON into struct
	var data DevfolioDetailResponse
	if err := json.Unmarshal(rawJSON, &data); err != nil {
		return nil, fmt.Errorf("failed to decode extracted JSON: %w", err)
	}

	return &data, nil
}

func fetchPrize(slug string, buildID string) (string, error) {
	hostname := fmt.Sprintf("%s.devfolio.co", slug)

	url := fmt.Sprintf("https://%s/_next/data/%s/hackathon3/%s/prizes.json?slug=%s",
		hostname, buildID, hostname, hostname)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "TBA", err
	}

	var prizeData DevfolioPrizeJSON
	if err := json.NewDecoder(resp.Body).Decode(&prizeData); err != nil {
		return "TBA", err
	}

	value := prizeData.PageProps.AggregatePrizeValue
	currency := prizeData.PageProps.AggregatePrizeCurrency

	if value == 0 || currency == "" {
		return "TBA", nil
	}

	return fmt.Sprintf("%.2f %s", value, currency), nil
}

// DevfolioAdapter converts Devfolio hackathon data to database format
func DevfolioAdapter(raw *DevfolioDetailResponse, actualPrize string, rate float64) (title, host, location, prize string, prizeUSD float64, startDate, endDate time.Time, applyURL string, err error) {
	h := raw.Props.PageProps.Hackathon

	// Map basic fields
	title = h.Name
	host = "Devfolio" // Hardcoded default since Devfolio doesn't explicitly name a host
	location = h.Location
	prize = actualPrize

	totalCash := raw.Props.PageProps.AggregatePrizeValue
	currency := raw.Props.PageProps.AggregatePrizeCurrency

	prizeUSD = 0.0
	switch currency {
	case "INR":
		prizeUSD = totalCash / rate

	case "USD":
		prizeUSD = totalCash
	default:
		prizeUSD = totalCash
	}

	applyURL = fmt.Sprintf("https://%s.devfolio.co/", h.Slug)

	// Parse dates from RFC3339 format
	startDate, err = time.Parse(time.RFC3339, h.StartsAt)
	if err != nil {
		return "", "", "", "", 0, time.Time{}, time.Time{}, "", fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err = time.Parse(time.RFC3339, h.EndsAt)
	if err != nil {
		return "", "", "", "", 0, time.Time{}, time.Time{}, "", fmt.Errorf("invalid end date format: %v", err)
	}

	return title, host, location, prize, prizeUSD, startDate, endDate, applyURL, nil
}
