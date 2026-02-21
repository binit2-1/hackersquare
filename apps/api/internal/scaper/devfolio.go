package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/binit2-1/hackersquare/apps/api/internal/database"
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
			Hackathon struct {
				Name     string `json:"name"`
				Slug     string `json:"slug"`
				Desc     string `json:"desc"`
				StartsAt string `json:"starts_at"`
				EndsAt   string `json:"ends_at"`
				Location string `json:"location"`
			} `json:"hackathon"`
		} `json:"pageProps"`
	} `json:"props"`
}

// RunDevfolioScraper fetches hackathons from Devfolio and saves them to the database
func RunDevfolioScraper(db *database.Service) error {
	fmt.Println("🚀 [Scraper] Booting up Devfolio ETL Pipeline...")

	// Extract the dynamic build ID required for Devfolio's API
	buildID, err := getBuildID()
	if err != nil {
		return fmt.Errorf("failed to get build ID: %w", err)
	}
	fmt.Printf("🔍 [Scraper] Secured Build ID: %s\n", buildID)

	// Fetch the Index API to get the list of slugs
	indexURL := fmt.Sprintf("https://devfolio.co/_next/data/%s/hackathons.json", buildID)
	indexData, err := fetchIndex(indexURL)
	if err != nil {
		return fmt.Errorf("failed to fetch index: %w", err)
	}

	if len(indexData.PageProps.DehydratedState.Queries) == 0 {
		return fmt.Errorf("no queries found in Devfolio index data")
	}

	hackathons := indexData.PageProps.DehydratedState.Queries[0].State.Data.OpenHackathons
	fmt.Printf("📋 [Scraper] Found %d open hackathons. Beginning extraction...\n", len(hackathons))

	// Loop through the index and fetch Details via subdomain HTML
	for i, hack := range hackathons {
		fmt.Printf("⏳ [%d/%d] Extracting data for: %s...\n", i+1, len(hackathons), hack.Slug)

		// Construct the subdomain URL
		subdomainURL := fmt.Sprintf("https://%s.devfolio.co/", hack.Slug)

		detailData, err := fetchDetail(subdomainURL)
		if err != nil {
			fmt.Printf("   ❌ Error fetching %s: %v\n", hack.Slug, err)
			continue // Skip failed hackathons, don't stop the entire scraper
		}

		descLength := len(detailData.Props.PageProps.Hackathon.Desc)
		hackName := detailData.Props.PageProps.Hackathon.Name
		fmt.Printf("   ✅ Success: Extracted %d bytes of markdown description for '%s'\n", descLength, hackName)

		title, host, loc, prize, start, end, applyURL, err := DevfolioAdapter(detailData)
		if err != nil {
			fmt.Printf("   ⚠️ Adapter failed for %s: %v\n", hack.Slug, err)
			continue
		}

		fmt.Printf("   🔄 Transformed: %s | Host: %s | Loc: %s | Prize: %s | Starts: %s | Ends: %s | URL: %s\n",
			title,
			host,
			loc,
			prize,
			start.Format("Jan 02, 2006"),
			end.Format("Jan 02, 2006"),
			applyURL,
		)

		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM hackathons WHERE "applyUrl" = $1)`
		
		err = db.Pool.QueryRow(context.Background(), checkQuery, applyURL).Scan(&exists)
		if err != nil {
			fmt.Printf("   ❌ Database Check Failed for %s: %v\n", title, err)
			continue
		}

		// Skip if already in database
		if exists {
			fmt.Printf("   ⏭️ Skipping '%s': Already exists in database.\n", title)
			continue
		}

		// Insert new hackathon
		newID := uuid.New().String()
		tags := []string{"Devfolio", "Upcoming"}

		insertQuery := `
			INSERT INTO hackathons (id, title, host, location, prize, "startDate", "endDate", "applyUrl", tags, "updatedAt")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

		_, err = db.Pool.Exec(context.Background(), insertQuery, newID, title, host, loc, prize, start, end, applyURL, tags, time.Now())
		if err != nil {
			fmt.Printf("   ❌ Database Insert Failed for %s: %v\n", title, err)
			continue
		}

		fmt.Printf("   💾 Successfully saved '%s' to PostgreSQL!\n", title)


		// Rate limiting to avoid being blocked
		time.Sleep(1 * time.Second)

		// TEST LIMITER: Break after 2 iterations during development so you don't wait 50 seconds
		// Remove this 'if' block when you are ready to scrape the entire site.
		
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

// DevfolioAdapter converts Devfolio hackathon data to database format
func DevfolioAdapter(raw *DevfolioDetailResponse) (title, host, location, prize string, startDate, endDate time.Time, applyURL string, err error) {
	h := raw.Props.PageProps.Hackathon

	// Map basic fields
	title = h.Name
	host = "Devfolio" // Hardcoded default since Devfolio doesn't explicitly name a host
	location = h.Location
	prize = "TBA" // You can add logic later to extract prize money from the Desc
	applyURL = fmt.Sprintf("https://%s.devfolio.co/", h.Slug)

	// Parse dates from RFC3339 format
	startDate, err = time.Parse(time.RFC3339, h.StartsAt)
	if err != nil {
		return "", "", "", "", time.Time{}, time.Time{}, "", fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err = time.Parse(time.RFC3339, h.EndsAt)
	if err != nil {
		return "", "", "", "", time.Time{}, time.Time{}, "", fmt.Errorf("invalid end date format: %v", err)
	}

	return title, host, location, prize, startDate, endDate, applyURL, nil
}
