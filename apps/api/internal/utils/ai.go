package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/ollama/ollama/api"
)


type authTransport struct {
	Transport http.RoundTripper
	APIKey    string
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+t.APIKey)
	return t.Transport.RoundTrip(r)
}

func GenerateAIOverview(githubData string) (string, error) {
	apiKey := os.Getenv("OLLAMA_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OLLAMA_API_KEY environment variable not set")
	}

	cloudURL, _ := url.Parse("https://ollama.com")
	httpClient := &http.Client{
		Transport: &authTransport{
			Transport: http.DefaultTransport,
			APIKey:    apiKey,
		},
	}
	client := api.NewClient(cloudURL, httpClient)

	systemPrompt := `You are a Senior Technical Profiler.
Your task is to analyze a developer's raw GitHub repository data and produce a structured, high-signal README summary optimized for a professional portfolio.

Critical objective:
- Capture true development behavior based on repository names, languages, and descriptions.
- Never hallucinate skills, languages, or projects that do not exist in the data.

Required output format (Markdown):
## 🧑‍💻 Developer Overview
[A strict 2-3 sentence executive summary of their profile based purely on the data]

### 🛠️ Runtime & Tech Stack
- [List primary languages explicitly found in the data]
- [List primary frameworks or tools discovered]

### 🏗️ Project Archetypes
- [Identify the types of tools they build: e.g., Low-level Systems, Web Scrapers, CLI Tools, UI Components]

### ⚙️ Developer Philosophy
- [Infer their focus based on project scopes: e.g., Open-source contributor, heavy automation focus, UI/UX centric]

Strict rules:
- Every bullet point must be backed by evidence in the provided data.
- If evidence is weak for a category, omit the category.
- Do not use generic filler words ("passionate", "ninja"). Be objective and precise.
- Output raw Markdown only. Do not wrap the response in a code block.`


	stream := false
	req := &api.ChatRequest{
		Model: "minimax-m2.5:cloud", 
		Messages: []api.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fmt.Sprintf("Analyze this GitHub repository data:\n%s", githubData)},
		},
		Stream: &stream,
	}

	var finalSummary string
	err := client.Chat(context.Background(), req, func(resp api.ChatResponse) error {
		finalSummary += resp.Message.Content
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("agent.project-summariser failed: %w", err)
	}

	if finalSummary == "" {
		return "", fmt.Errorf("Project summariser returned an empty summary")
	}

	return finalSummary, nil
}