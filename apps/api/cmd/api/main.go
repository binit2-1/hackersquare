package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/binit2-1/hackersquare/apps/api/internal/database"
	"github.com/binit2-1/hackersquare/apps/api/internal/hackathon"
	"github.com/gorilla/mux"
)

func main(){

	//Initialize db
	dbService, err := database.New()
	if err != nil{
		log.Fatalf("Failed to initialize database service: %v", err)
	}

	defer dbService.Close() //'defer' guarantees this runs right before main() exits

	// Initialize Hackathon Handler and inject the DB!
    hackathonHandler := hackathon.NewHandler(dbService)

	//create router
	mux := mux.NewRouter()

	//health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"messages": "HackerSquare API is running smoothly!"}`))
	})

	//ROUTES
	mux.HandleFunc("/api/hackathons", hackathonHandler.GetHackathons).Methods("GET")

	//start server
	port := os.Getenv("PORT")
	if port == ""{
		port = ":8080" //default port if not set in env
	}
	fmt.Printf("Starting server on port %s...\n", port)

	//http.ListenAndServe blocks the main thread to keep the sever alive 
	err = http.ListenAndServe(port, mux)
	if err!= nil{
		log.Fatalf("Failed to start server: %v", err)
	} 
}

