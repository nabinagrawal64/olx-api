package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/nabinagrawal64/olx-api/internal/config"
	"github.com/nabinagrawal64/olx-api/internal/db"
	"github.com/nabinagrawal64/olx-api/internal/handlers"
	"github.com/nabinagrawal64/olx-api/internal/middleware"
)

func main() {
	// Load configurations
	cfg := config.MustLoad()

	// Connect the database
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} 
	
	// Error Log Format
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelInfo,
	}) 
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Create router
	lh := handlers.NewListingHandlers(db, logger)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", lh.GetListings)
	mux.HandleFunc("DELETE /listings/{id}", lh.DeleteListing)

	// Initialize HTTP Server
	handler := middleware.RequestId(mux)
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
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