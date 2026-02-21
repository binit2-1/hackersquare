package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main(){
	//create router
	mux := mux.NewRouter()

	//health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"messages": "HackerSquare API is running smoothly!"}`))
	})

	//start server
	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)

	//http.ListenAndServe blocks the main thread to keep the sever alive 
	err := http.ListenAndServe(port, mux)
	if err!= nil{
		log.Fatalf("Failed to start server: %v", err)
	} 
}

