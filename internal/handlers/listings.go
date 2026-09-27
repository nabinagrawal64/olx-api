package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type listings struct{
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

func Listings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		rows, err := db.Query(`
			SELECT id, title, price, description, city, created_at FROM listings
			ORDER BY created_at DESC
			LIMIT 100;
		`)
		if err != nil{
			log.Printf("Query: %v", err)
			http.Error(w, "Failed to query listings", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var ls []listings
		for rows.Next(){
			var l listings
			if err := rows.Scan(&l.ID, &l.Title, &l.Price, &l.Description, &l.City, &l.CreatedAt); err != nil {
				log.Printf("Scan Error: %v", err)
				http.Error(w, "Failed to scan listings", http.StatusInternalServerError)
				return
			} 
			ls = append(ls, l)
		}
		if err := rows.Err(); err != nil {
			log.Printf("Rows Error: %v", err)
			http.Error(w, "Failed to scan listings", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		err = json.NewEncoder(w).Encode(ls)
		if err != nil {
			log.Printf("Encode Error: %v", err)
			http.Error(w, "Failed to encode listings", http.StatusInternalServerError)
			return
		}
	}
}  