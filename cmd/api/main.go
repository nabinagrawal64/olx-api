package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nabinagrawal64/olx-api/internal/config"
	"github.com/nabinagrawal64/olx-api/internal/db"
	"github.com/nabinagrawal64/olx-api/internal/handlers"
)

func main() {
	// Load configurations
	cfg := config.MustLoad()

	// Connect the database
	_, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} 

	// Create router
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)

	// Initialize HTTP Server
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
  
	log.Printf("Starting server on %s", srv.Addr)
	
	if err := srv.ListenAndServe(); err != nil { 
		log.Fatalf("Server Failed %v",err)
	}   
}


// https://olx-api-58x8.onrender.com/healthz