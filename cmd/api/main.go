package main

import (
	// "fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux();

	mux.HandleFunc("GET /healthz", 
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))// using backtick to create raw string literal
		},
	)

	srv := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
  
	log.Printf("Starting server on :%s", srv.Addr)
	err := srv.ListenAndServe();
 
	if err != nil {
		log.Fatalf("Server Failed %v",err)
	}   
}
