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

	systemPrompt := `You are an expert Developer Advocate and Technical Profiler.
Your task is to analyze a developer's raw GitHub repository data and produce a structured, high-signal Markdown profile README optimized for their portfolio.

Critical objective:
- Capture true development behavior based on repository names, descriptions, and languages.
- Never hallucinate skills, languages, or projects that do not exist in the data.

Required Output Format (Strictly follow this Markdown structure):

## Hi, I'm a [Infer Role, e.g., Fullstack Developer, Systems Engineer, Frontend Engineer].

[Write a concise, 2-3 sentence bio summarizing what they build, their primary ecosystem, and their developer focus based on the data.]

### 🛠 Tech Stack

* **Core:** [List primary languages and major frameworks found in the data]
* **Tools:** [List inferred tools, e.g., Git, Docker, etc., based on project types]

### 🚀 Featured Projects

[Select the top 2 or 3 most impressive original repositories based on stars, topics, and descriptions. For each, use this exact format:]

**[Project Name]**

* [1-2 bullet points explaining what it is and its features based on the description]
* **Tech:** [List the languages and topics used in this specific project]\

---`


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


func GenerateSearchInsights(profileReadme string, searchQuery string) (string, error) {
	apiKey := os.Getenv("OLLAMA_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OLLAMA_API_KEY environment variable not set")
	} 


		baseURL, _ := url.Parse("https://ollama.com")
		httpClient := &http.Client{
			Transport: &authTransport{
				Transport: http.DefaultTransport,
				APIKey:    apiKey,
			},
		}

		client := api.NewClient(baseURL, httpClient)

		systemPrompt := `You are an expert Developer Advocate. 
Your task is to provide a highly concise, 2-sentence insight on how a user's search query aligns with their developer profile.

Strict rules:
- Keep it under 5 to 6 sentences. Be punchy and direct.
- Identify 1 specific strength from their profile that gives them an edge for this type of hackathon.
- Do not use formatting like headers or code blocks. Simple text with occasional bolding is fine.
- If the search query is vague (like "near me"), focus on their general tech stack's versatility.
- Add some tips based on the hackathons title which hackathons the user should participate based their profile his location, etc.`

	userMessage := fmt.Sprintf("User Profile:\n%s\n\nSearch Query: %s", profileReadme, searchQuery)

	stream := false 
	req := &api.ChatRequest{
		Model: "minimax-m2.5:cloud", 
		Messages: []api.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		Stream: &stream,
	}

	var insight string
	err := client.Chat(context.Background(), req, func(resp api.ChatResponse) error {
		insight += resp.Message.Content
		return nil
	})

	if err != nil {
		return "", err
	}

	return insight, nil
}