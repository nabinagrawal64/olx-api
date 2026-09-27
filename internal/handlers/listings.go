package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type listings struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandlers struct {
	db *sql.DB
}

func NewListingHandlers(db *sql.DB) *ListingHandlers {
	return &ListingHandlers{
		db: db,
	}
}

func (lh *ListingHandlers) GetListings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Query to fetch all listings
	rows, err := lh.db.QueryContext(ctx, `
			SELECT id, title, price, description, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100;
		`)
	if err != nil {
		log.Printf("Query: %v", err)
		http.Error(w, "Failed to query listings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate and store in slice
	var ls []listings
	for rows.Next() {
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

	// JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(ls)
	if err != nil {
		log.Printf("Encode Error: %v", err)
		http.Error(w, "Failed to encode listings", http.StatusInternalServerError)
		return
	}
}

func (lh *ListingHandlers) DeleteListing(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()
	
	// Delete listing
	result, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		log.Printf("Delete: %v", err)
		http.Error(w, "Failed to delete listing", http.StatusInternalServerError)
		return
	} 

	// Get affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Rows Affected Error: %v", err)
		http.Error(w, "Failed to get affected rows", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Listing not found", http.StatusNotFound)
		return
	}

	// JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Listing deleted successfully",
	})
}
