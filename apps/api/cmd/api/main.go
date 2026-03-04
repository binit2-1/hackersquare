package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/binit2-1/hackersquare/apps/api/internal/repository/pg"
	scraper "github.com/binit2-1/hackersquare/apps/api/internal/scaper"
	"github.com/binit2-1/hackersquare/apps/api/internal/server"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)



func main(){

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
		port = ":8080"
	}


	db, err := sql.Open("pgx", dbConnectionURL)
	if err != nil {
		log.Fatalf("FATAL: Failed to parse Postgres configuration: %v", err)
	}

	defer db.Close()


	pgRepo := pg.NewPostgreEventRepo(db)
	authRepo := pg.NewPostgreUserRepo(db)
	hackathonHandler := server.NewHackathonHandler(pgRepo)
	authHandler := server.NewAuthHandler(authRepo)
	

	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /v1/auth/register", authHandler.Register)
	mux.HandleFunc("GET /v1/search", hackathonHandler.SearchHackathons)


	fmt.Printf("Starting server on port %s\n", port)

	//scrappers
	//DEVFOLIO
	go func() {
		if err := scraper.RunDevfolioScraper(db); err != nil {
			fmt.Printf("Scraper Error: %v\n", err)
		}
	}()
	//MLH
	go func() {
		if err := scraper.RunMLHScraper(db); err != nil {
			fmt.Printf("❌ MLH Scraper Error: %v\n", err)
		}
	}()
	//UNSTOP
	go func(){
		if err := scraper.RunUnstopScraper(db); err != nil{
			fmt.Printf("❌ Unstop Scraper Error: %v\n", err)
		}
	}()

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("FATAL: Server crashed: %v", err)
	} 
	

}