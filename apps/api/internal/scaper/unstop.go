package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/database"
	"github.com/binit2-1/hackersquare/apps/api/internal/utils"
	"github.com/google/uuid"
)

type UnstopAPIResponse struct {
	Data struct {
		Data        []UnstopRawEvents `json:"data"`
		NextPageURL string            `json:"next_page_url"`
		LastPage    int               `json:"last_page"`
	} `json:"data"`
}

type UnstopRawEvents struct {
	Title        string `json:"title"`
	SeoURL       string `json:"seo_url"`
	Organisation struct {
		Name string `json:"name"`
	} `json:"organisation"`
	AddressWithCountryLogo struct {
		City    string `json:"city"`
		Country struct {
			Name string `json:"name"`
		} `json:"country"`
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

func RunUnstopScraper(db *database.Service) error {
	utils.Info("[Unstop] Starting scraper...")
	rate := utils.GetExchangeRate()
	indexURL := "https://unstop.com/api/public/opportunity/search-result?opportunity=hackathons&page=1"
	indexData, err := fetchUnstopIndex(indexURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Unstop index: %w", err)
	}

	hackathons := indexData.Data.Data
	if len(hackathons) == 0 {
		return fmt.Errorf("no hackathons found in Unstop index")
	}

	utils.Info("Found %d open hackathons", len(hackathons))

	// Pagination
	for indexURL != "" {
		utils.Debug("Fetching page: %s", indexURL)

		indexData, err := fetchUnstopIndex(indexURL)
		if err != nil {
			return err
		}

		pageHackathons := indexData.Data.Data
		if len(pageHackathons) == 0 {
			break
		}

		for i, hack := range pageHackathons {
			utils.Debug("[%d/%d] Extracting: %s", i+1, len(pageHackathons), hack.SeoURL)
			title, host, loc, prize, prizeUSD, start, end, applyURL, err := UnstopAdapter(&hack, rate)

			if err != nil {
				utils.Error("Adapter failed for %s: %v", hack.SeoURL, err)
				continue
			}

			utils.Debug("Transformed: %s | Host: %s | Loc: %s | Prize: %s | USD: %.2f", title, host, loc, prize, prizeUSD)

			var exists bool
			checkQuery := `SELECT EXISTS(SELECT 1 FROM hackathons WHERE "applyUrl" = $1)`
			err = db.Pool.QueryRow(context.Background(), checkQuery, applyURL).Scan(&exists)
			if err != nil {
				utils.Error("Database check failed for %s: %v", title, err)
				continue
			}

			// Skip if already in database
			if exists {
				utils.Debug("Skipping '%s': Already exists", title)
				continue
			}

			// Insert new hackathon into db
			newID := uuid.New().String()
			tags := []string{"Unstop", "Upcoming"}
			insertQuery := `INSERT INTO hackathons (id, title, host, location, prize, "prizeUSD", "startDate", "endDate", "applyUrl", tags, "updatedAt")
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
			utils.Debug("Inserting: %s | Prize: %s | USD: %.2f", title, prize, prizeUSD)
			_, err = db.Pool.Exec(context.Background(), insertQuery,
				newID,
				title,
				host,
				loc,
				prize,
				prizeUSD,
				start,
				end,
				applyURL,
				tags,
				time.Now(),
			)

			if err != nil {
				utils.Error("Database insert failed for %s: %v", title, err)
				continue
			}

			utils.Debug("Successfully inserted '%s'", title)

			// Rate limit
			time.Sleep(1 * time.Second)
		}

		indexURL = indexData.Data.NextPageURL
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

func UnstopAdapter(raw *UnstopRawEvents, rate float64) (title, host, location, prize string, prizeUSD float64, startDate, endDate time.Time, applyURL string, err error) {

	//map basic fields
	title = raw.Title
	host = "Unstop"
	if raw.Organisation.Name != "" {
		host = raw.Organisation.Name
	}

	location = "Remote"
	if raw.AddressWithCountryLogo.City != "" {
		location = fmt.Sprintf("%s, %s", raw.AddressWithCountryLogo.City, raw.AddressWithCountryLogo.Country.Name)
	}

	var totalCash float64
	currency := "INR"

	for _, p := range raw.Prizes {
		totalCash += p.Cash
		if p.CurrencyCode != "" {
			currency = p.CurrencyCode
		}
	}

	if totalCash > 0 {
		prize = fmt.Sprintf("%.2f %s", totalCash, currency)
	} else if len(raw.Prizes) > 0 && raw.Prizes[0].Others != "" {
		prize = raw.Prizes[0].Others + " and more"
	} else {
		prize = "TBA"
	}

	prizeUSD = 0.0
	switch currency {
	case "INR":
		prizeUSD = totalCash / rate
	case "USD":
		prizeUSD = totalCash
	}

	applyURL = fmt.Sprintf("%s", raw.SeoURL)

	return title, host, location, prize, prizeUSD, raw.RegnRequirements.StartRegnDt, raw.EndDate, applyURL, nil

}
