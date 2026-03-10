package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func GetExchangeRate() float64 {

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get("https://api.frankfurter.dev/v1/latest?base=USD&symbols=INR")
	if err != nil {
		return 90.92 // falback
	}
	defer resp.Body.Close()

	var response ExchangeRateResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil || response.Rates["INR"] == 0 {
		return 90.92 // fallback if decoding fails
	}

	return response.Rates["INR"]
}
