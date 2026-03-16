package domain

import "context"

type AIService interface {
	GenerateProfileReadme(ctx context.Context, githubData string) (string, error)
	GenerateSearchInsights(ctx context.Context, profileReadme, searchQuery, hackathonsContext string) (string, error)
}
