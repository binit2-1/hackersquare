package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)



func main(){

	err := godotenv.Load()
	if err != nil {
		log.Println("WARN: No .env file found. Relying on system environment variables.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	mux := http.NewServeMux()


	fmt.Printf("Starting server on port %s\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("FATAL: Server crashed: %v", err)
	} 
	

}