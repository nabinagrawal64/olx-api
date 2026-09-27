package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nabinagrawal64/olx-api/internal/config"
)

func main() {
	cfg := config.MustLoad()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", 
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))// using backtick to create raw string literal
		},
	)

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
