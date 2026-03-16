package ai

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/binit2-1/hackersquare/apps/api/internal/domain"
	"github.com/ollama/ollama/api"
)

type AuthTransport struct {
	Transport http.RoundTripper
	APIKey    string
}

func (t *AuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Bearer "+t.APIKey)
	return t.Transport.RoundTrip(r)
}

type OllamaService struct {
	client *api.Client
	model  string
}

func ensureHomeForRuntime() error {
	if os.Getenv("HOME") != "" {
		return nil
	}

	fallbackHome := filepath.Join(os.TempDir(), "hackersquare-home")
	if err := os.MkdirAll(fallbackHome, 0o700); err != nil {
		return fmt.Errorf("failed to create fallback HOME: %w", err)
	}

	if err := os.Setenv("HOME", fallbackHome); err != nil {
		return fmt.Errorf("failed to set fallback HOME: %w", err)
	}

	return nil
}

func NewOllamaService(apiKey string, model string) (domain.AIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OLLAMA_API_KEY environment variable not set")
	}

	if err := ensureHomeForRuntime(); err != nil {
		return nil, err
	}

	cloudURL, err := url.Parse(`https://ollama.com`)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ollama URL: %w", err)
	}

	httpClient := &http.Client{
		Transport: &AuthTransport{
			Transport: http.DefaultTransport,
			APIKey:    apiKey,
		},
	}

	return &OllamaService{
		client: api.NewClient(cloudURL, httpClient),
		model:  model,
	}, nil
}

func (s *OllamaService) GenerateProfileReadme(ctx context.Context, githubData string) (string, error) {
	stream := false
	req := &api.ChatRequest{
		Model: s.model,
		Messages: []api.Message{
			{Role: "system", Content: GenerateProfileReadmePrompt},
			{Role: "user", Content: fmt.Sprintf("Analyze this GitHub repository data:\n%s", githubData)},
		},
		Stream: &stream,
	}

	var finalSummary string

	err := s.client.Chat(ctx, req, func(response api.ChatResponse) error {
		finalSummary += response.Message.Content
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

func (s *OllamaService) GenerateSearchInsights(ctx context.Context, profileReadme, searchQuery, hackathonsContext string) (string, error) {
	userMessage := fmt.Sprintf("User Profile:\n%s\n\nSearch Query: %s\n\nTop Search Results:\n%s", profileReadme, searchQuery, hackathonsContext)

	stream := false
	req := &api.ChatRequest{
		Model: s.model,
		Messages: []api.Message{
			{Role: "system", Content: SearchInsightsPrompt},
			{Role: "user", Content: userMessage},
		},
		Stream: &stream,
	}

	var insights string

	err := s.client.Chat(ctx, req, func(response api.ChatResponse) error {
		insights += response.Message.Content
		return nil
	})

	if err != nil {
		return "", err
	}

	return insights, nil
}
