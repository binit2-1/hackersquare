package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/database"
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
	fmt.Println("🚀 [Unstop Scraper] Booting up...")

	indexURL := "https://unstop.com/api/public/opportunity/search-result?opportunity=hackathons&page=1"
	indexData, err := fetchUnstopIndex(indexURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Unstop index: %w", err)
	}

	hackathons := indexData.Data.Data
	if len(hackathons) == 0 {
		return fmt.Errorf("no hackathons found in Unstop index")
	}

	fmt.Printf("📋 [Scraper] Found %d open hackathons. Beginning extraction...\n", len(hackathons))

	//Pagination
	for indexURL != "" {
		fmt.Printf("🚀 [Unstop] Fetching page: %s\n", indexURL)

		indexData, err := fetchUnstopIndex(indexURL)
		if err != nil {
			return err
		}

		pageHackathons := indexData.Data.Data
		if len(pageHackathons) == 0 {
			break
		}

		for i, hack := range pageHackathons {
			fmt.Printf("⏳ [%d/%d] Extracting data for: %s...\n", i+1, len(pageHackathons), hack.SeoURL)
			title, host, loc, prize, start, end, applyURL, err := UnstopAdapter(&hack)

			if err != nil {
				fmt.Printf("   ⚠️ Adapter failed for %s: %v\n", hack.SeoURL, err)
				continue
			}

			fmt.Printf("   🔄 Transformed: %s | Host: %s | Loc: %s | Prize: %s | Starts: %s | Ends: %s | URL: %s\n",
				title, host, loc, prize, start.Format("Jan 02, 2006"), end.Format("Jan 02, 2006"), applyURL,
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

			//Insert new hackathon into db
			newID := uuid.New().String()
			tags := []string{"Unstop", "Upcoming"}
			insertQuery := `INSERT INTO hackathons (id, title, host, location, prize, "startDate", "endDate", "applyUrl", tags, "updatedAt")
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

			_, err = db.Pool.Exec(context.Background(), insertQuery,
				newID,
				title,
				host,
				loc,
				prize,
				start,
				end,
				applyURL,
				tags,
				time.Now(),
			)

			if err != nil {
				fmt.Printf("   ❌ Database Insert Failed for %s: %v\n", title, err)
				continue
			}

			fmt.Printf("   💾 Successfully inserted '%s' into database.\n", title)

			//Rate limit to avoid hitting Unstop's servers too hard
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

func UnstopAdapter(raw *UnstopRawEvents) (title, host, location, prize string, startDate, endDate time.Time, applyURL string, err error) {

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

	applyURL = fmt.Sprintf("%s", raw.SeoURL)

	return title, host, location, prize, raw.RegnRequirements.StartRegnDt, raw.EndDate, applyURL, nil

}
