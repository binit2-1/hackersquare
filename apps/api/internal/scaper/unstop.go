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
		Data []UnstopRawEvents `json:"data"`
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
}

func RunUnstopScraper(db *database.Service) error {
	fmt.Println("🚀 [Unstop Scraper] Booting up...")

	indexURL := "https://unstop.com/api/public/opportunity/search-result?opportunity=hackathons"
	indexData, err := fetchUnstopIndex(indexURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Unstop index: %w", err)
	}

	hackathons := indexData.Data.Data
	if len(hackathons) == 0 {
		return fmt.Errorf("no hackathons found in Unstop index")
	}

	fmt.Printf("📋 [Scraper] Found %d open hackathons. Beginning extraction...\n", len(hackathons))

	for i, hack := range hackathons {
		fmt.Printf("⏳ [%d/%d] Extracting data for: %s...\n", i+1, len(hackathons), hack.SeoURL)
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

	return nil
}


//Helpers
func fetchUnstopIndex(url string) (*UnstopAPIResponse, error) {
	resp, err := http.Get(url)
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

func UnstopAdapter(raw *UnstopRawEvents) (title, host, location, prize string, startDate, endDate time.Time, applyURL string, err error){
	
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

	prize = "TBA"

	applyURL = fmt.Sprintf("%s", raw.SeoURL)

	return title, host, location, prize, raw.RegnRequirements.StartRegnDt, raw.EndDate, applyURL, nil

}