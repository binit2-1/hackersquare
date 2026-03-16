package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type GitHubRepo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Stargazers  int       `json:"stargazers_count"`
	Fork        bool      `json:"fork"`
	Topics      []string  `json:"topics"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FetchGithubData(username string) (string, error) {

	apiURL := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated&per_page=15", username)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("github api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var repos []GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return "", fmt.Errorf("failed to decode github response: %w", err)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Target Developer GitHub Handle: %s\n\n", username)
	sb.WriteString("Recent Original Repositories:\n")

	validReposCount := 0

	for _, repo := range repos {
		// Skip forks to ensure the AI only analyzes code the user actually wrote
		if repo.Fork {
			continue
		}

		validReposCount++
		fmt.Fprintf(&sb, "- Repository: %s\n", repo.Name)

		if repo.Description != "" {
			fmt.Fprintf(&sb, "  Description: %s\n", repo.Description)
		}
		if repo.Language != "" {
			fmt.Fprintf(&sb, "  Primary Language: %s\n", repo.Language)
		}
		if len(repo.Topics) > 0 {
			fmt.Fprintf(&sb, "  Topics: %s\n", strings.Join(repo.Topics, ", "))
		}

		fmt.Fprintf(&sb, "  Stars: %d\n\n", repo.Stargazers)
	}

	if validReposCount == 0 {
		return "", fmt.Errorf("no original public repositories found for user %s", username)
	}

	return sb.String(), nil
}
