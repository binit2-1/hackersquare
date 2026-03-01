package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/binit2-1/hackersquare/apps/api/internal/repository/pg"
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
	hackathonHandler := server.NewHackathonHandler(pgRepo)
	

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/search", hackathonHandler.SearchHackathons)


	fmt.Printf("Starting server on port %s\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("FATAL: Server crashed: %v", err)
	} 
	

}