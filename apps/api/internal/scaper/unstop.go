package scraper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/utils"
)

type UnstopAPIResponse struct {
	Data struct {
		Data        []UnstopRawEvents `json:"data"`
		NextPageURL string            `json:"next_page_url"`
		LastPage    int               `json:"last_page"`
	} `json:"data"`
}

type UnstopRawEvents struct {
	Title        string `json:"title,omitempty"`
	SeoURL       string `json:"seo_url"`
	Organisation struct {
		Name string `json:"name"`
	} `json:"organisation"`
	AddressWithCountryLogo struct {
		City    string `json:"city"`
		Country any `json:"country"`
	} `json:"address_with_country_logo"`
	EndDate          time.Time `json:"end_date"`
	RegnRequirements struct {
		StartRegnDt time.Time `json:"start_regn_dt"`
	} `json:"regnRequirements"`
	Prizes []UnstopPrize `json:"prizes"`
}

type UnstopPrize struct {
	Cash         float64 `json:"cash"`
	CurrencyCode string  `json:"currencyCode"`
	Others       string  `json:"others"`
}

func RunUnstopScraper(db *sql.DB) error {
	utils.Info("[Unstop] Starting scraper...")
	rate := utils.GetExchangeRate()
	indexURL := "https://unstop.com/api/public/opportunity/search-result?opportunity=hackathons&page=1"

	// Simplified pagination loop to prevent double-fetching Page 1
	for indexURL != "" {
		utils.Debug("Fetching page: %s", indexURL)

		indexData, err := fetchUnstopIndex(indexURL)
		if err != nil {
			return fmt.Errorf("failed to fetch Unstop page: %w", err)
		}

		pageHackathons := indexData.Data.Data
		if len(pageHackathons) == 0 {
			break
		}

		for i, hack := range pageHackathons {
			utils.Debug("[%d/%d] Extracting: %s", i+1, len(pageHackathons), hack.SeoURL)
			
			title, host, loc, prizeUSD, start, end, applyURL, err := UnstopAdapter(&hack, rate)
			if err != nil {
				utils.Error("Adapter failed for %s: %v", hack.SeoURL, err)
				continue
			}

			utils.Debug("Transformed: %s | Host: %s | Loc: %s | USD: %.2f", title, host, loc, prizeUSD)

			// The V2 PostgreSQL Upsert Query
			upsertQuery := `
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

			_, err = db.ExecContext(context.Background(), upsertQuery, title, host, loc, prizeUSD, start, end, applyURL)
			if err != nil {
				utils.Error("Database upsert failed for %s: %v", title, err)
				continue
			}

			utils.Debug("Successfully saved/updated '%s'", title)
		}

		// Grab the next page URL (API returns empty string if no more pages)
		indexURL = indexData.Data.NextPageURL
		
		// Rate limit between pages to avoid IP bans
		time.Sleep(2 * time.Second)
	}
	
	return nil
}

// Helpers
func fetchUnstopIndex(url string) (*UnstopAPIResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP: %d", resp.StatusCode)
	}

	var data UnstopAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

// V2 Adapter - Drops the string prize and returns pure floats for Postgres
func UnstopAdapter(raw *UnstopRawEvents, rate float64) (title, host, location string, prizeUSD float64, startDate, endDate time.Time, applyURL string, err error) {
	// Map basic fields
	title = raw.Title
	host = "Unstop"
	if raw.Organisation.Name != "" {
		host = raw.Organisation.Name
	}

	location = "Remote"
	city := raw.AddressWithCountryLogo.City
	countryName := ""
	if cStr, ok := raw.AddressWithCountryLogo.Country.(string); ok {
		countryName = cStr
	} else if cMap, ok := raw.AddressWithCountryLogo.Country.(map[string]any); ok {
		if n, ok := cMap["name"].(string); ok {
			countryName = n
		}
	}

	if city != "" && countryName != "" {
		location = fmt.Sprintf("%s, %s", city, countryName)
	} else if city != "" {
		location = city
	} else if countryName != "" {
		location = countryName
	}

	var totalCash float64
	currency := "INR"

	for _, p := range raw.Prizes {
		totalCash += p.Cash
		if p.CurrencyCode != "" {
			currency = p.CurrencyCode
		}
	}

	prizeUSD = 0.0
	switch currency {
	case "INR":
		prizeUSD = totalCash / rate
	case "USD":
		prizeUSD = totalCash
	}

	applyURL = raw.SeoURL

	return title, host, location, prizeUSD, raw.RegnRequirements.StartRegnDt, raw.EndDate, applyURL, nil
}