package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"

    "github.com/binit2-1/hackersquare/apps/api/internal/domain"
)

// We define the simple JSON structures ourselves, bypassing the SDK entirely.
type chatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type chatRequest struct {
    Model    string        `json:"model"`
    Messages []chatMessage `json:"messages"`
    Stream   bool          `json:"stream"`
}

type chatResponse struct {
    Message chatMessage `json:"message"`
    Error   string      `json:"error,omitempty"`
}

type OllamaService struct {
    apiKey     string
    model      string
    httpClient *http.Client
}

// NewOllamaService now requires zero filesystem hacks or SDK initializations.
func NewOllamaService(apiKey string, model string) (domain.AIService, error) {
    if apiKey == "" {
        return nil, fmt.Errorf("OLLAMA_API_KEY environment variable not set")
    }

    return &OllamaService{
        apiKey:     apiKey,
        model:      model,
        httpClient: &http.Client{},
    }, nil
}

// makeChatRequest is a DRY helper function to handle the raw HTTP call
func (s *OllamaService) makeChatRequest(ctx context.Context, systemPrompt, userMessage string) (string, error) {
    reqBody := chatRequest{
        Model:    s.model,
        Stream:   false,
        Messages: []chatMessage{
            {Role: "system", Content: systemPrompt},
            {Role: "user", Content: userMessage},
        },
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }

    // Standard REST call directly to the API endpoint
    req, err := http.NewRequestWithContext(ctx, "POST", "https://ollama.com/api/chat", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+s.apiKey)

    resp, err := s.httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("http request failed: %w", err)
    }
    defer resp.Body.Close()

    bodyBytes, _ := io.ReadAll(resp.Body)

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
    }

    var chatResp chatResponse
    if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
        return "", fmt.Errorf("failed to decode response: %w", err)
    }

    if chatResp.Error != "" {
        return "", fmt.Errorf("API error: %s", chatResp.Error)
    }

    if chatResp.Message.Content == "" {
        return "", fmt.Errorf("received empty response from AI")
    }

    return chatResp.Message.Content, nil
}

func (s *OllamaService) GenerateProfileReadme(ctx context.Context, githubData string) (string, error) {
    userMessage := fmt.Sprintf("Analyze this GitHub repository data:\n%s", githubData)
    return s.makeChatRequest(ctx, GenerateProfileReadmePrompt, userMessage)
}

func (s *OllamaService) GenerateSearchInsights(ctx context.Context, profileReadme, searchQuery, hackathonsContext string) (string, error) {
    userMessage := fmt.Sprintf("User Profile:\n%s\n\nSearch Query: %s\n\nTop Search Results:\n%s", profileReadme, searchQuery, hackathonsContext)
    return s.makeChatRequest(ctx, SearchInsightsPrompt, userMessage)
}