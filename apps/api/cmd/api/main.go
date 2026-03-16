package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/repository/pg"
	scraper "github.com/binit2-1/hackersquare/apps/api/internal/scaper"
	"github.com/binit2-1/hackersquare/apps/api/internal/server"
	"github.com/binit2-1/hackersquare/apps/api/internal/service/ai"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("WARN: No .env file found. Relying on system environment variables.")
	}

	dbConnectionURL := os.Getenv("DATABASE_URL")
	if dbConnectionURL == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sql.Open("pgx", dbConnectionURL)
	if err != nil {
		log.Fatalf("FATAL: Failed to parse Postgres configuration: %v", err)
	}

	defer db.Close()

	aiService, err := ai.NewOllamaService(os.Getenv("OLLAMA_API_KEY"), "minimax-m2.5:cloud")
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}

	// Initialize repositories
	pgRepo := pg.NewPostgreEventRepo(db)
	authRepo := pg.NewPostgreUserRepo(db)
	bookmarkRepo := pg.NewPostgresBookmarkRepo(db)

	// Initialize handlers
	hackathonHandler := server.NewHackathonHandler(pgRepo, authRepo, aiService)
	bookmarkHandler := server.NewBookmarkHandler(bookmarkRepo)
	authHandler := server.NewAuthHandler(authRepo, aiService)

	mux := http.NewServeMux()

	//public routes
	mux.HandleFunc("POST /v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /v1/auth/logout", server.AuthMiddleware(authHandler.Logout))
	mux.HandleFunc("GET /v1/search", hackathonHandler.SearchHackathons)
	mux.HandleFunc("GET /v1/hackathons/nearby", hackathonHandler.NearbyHackathons)

	//protected routes
	mux.HandleFunc("POST /v1/bookmarks", server.AuthMiddleware(bookmarkHandler.AddBookmark))
	mux.HandleFunc("DELETE /v1/bookmarks", server.AuthMiddleware(bookmarkHandler.RemoveBookmark))
	mux.HandleFunc("GET /v1/bookmarks", server.AuthMiddleware(bookmarkHandler.GetBookmarksByUser))
	mux.HandleFunc("GET /v1/search/overview", server.AuthMiddleware(hackathonHandler.GetSearchOverview))

	//me
	mux.HandleFunc("GET /v1/auth/me", server.AuthMiddleware(authHandler.GetMe))

	//profile
	mux.HandleFunc("PUT /v1/users/profile", server.AuthMiddleware(authHandler.UpdateProfile))
	mux.HandleFunc("PUT /v1/users/profile/readme", server.AuthMiddleware(authHandler.UpdateProfileReadme))

	//oAuth
	mux.HandleFunc("GET /v1/auth/github/login", authHandler.GithubLogin)
	mux.HandleFunc("GET /v1/auth/github/login/callback", authHandler.GithubLoginCallback)
	mux.HandleFunc("GET /v1/auth/github/connect", authHandler.ConnectGithub)
	mux.HandleFunc("GET /v1/auth/github/callback", server.AuthMiddleware(authHandler.GithubCallback))

	//AI
	mux.HandleFunc("POST /v1/users/profile/generate-summary", server.AuthMiddleware(authHandler.GenerateProfileSummary))

	fmt.Printf("Starting server on port %s\n", port)

	//cleanup
	go func() {
		if count, err := pgRepo.DeleteExpiredHackathons(); err != nil {
			fmt.Printf("Error deleting expired hackathons: %v\n", err)
		} else if count > 0 {
			fmt.Printf("Startup Cleanup: Swept away %d expired hackathons\n", count)
		}

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if count, err := pgRepo.DeleteExpiredHackathons(); err != nil {
				fmt.Printf("Error deleting expired hackathons: %v\n", err)
			} else {
				fmt.Printf("Scheduled Cleanup: Swept away %d expired hackathons\n", count)
			}
		}
	}()

	// Scrapers: run once on startup, then every 12 hours
	go func() {
		runAllScrapers := func() {
			if err := scraper.RunDevfolioScraper(db); err != nil {
				fmt.Printf("Devfolio Scraper Error: %v\n", err)
			}
			if err := scraper.RunMLHScraper(db); err != nil {
				fmt.Printf("MLH Scraper Error: %v\n", err)
			}
			if err := scraper.RunUnstopScraper(db); err != nil {
				fmt.Printf("Unstop Scraper Error: %v\n", err)
			}
		}

		runAllScrapers()

		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("Scheduled scraper cycle starting...")
			runAllScrapers()
		}
	}()

	if err := http.ListenAndServe(":"+port, server.CORSMiddleware(mux)); err != nil {
		log.Fatalf("FATAL: Server crashed: %v", err)
	}

}
