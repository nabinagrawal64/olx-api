package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/nabinagrawal64/olx-api/internal/middleware"
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
	logger *slog.Logger
}

func NewListingHandlers(db *sql.DB, logger *slog.Logger) *ListingHandlers {
	return &ListingHandlers{
		db: db,
		logger: logger,
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
		lh.logger.Error("Query Listing Error", slog.Any("error", err))
		http.Error(w, "Failed to query listings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate and store in slice
	var ls []listings
	for rows.Next() {
		var l listings
		if err := rows.Scan(&l.ID, &l.Title, &l.Price, &l.Description, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("Scan Error", slog.Any("error", err))
			http.Error(w, "Failed to scan listings", http.StatusInternalServerError)
			return
		}
		ls = append(ls, l)
	}
	if err := rows.Err(); err != nil {
		lh.logger.Error("Rows Error", slog.Any("error", err))
		http.Error(w, "Failed to scan listings", http.StatusInternalServerError)
		return
	}

	lh.logger.Info("Fetched listings", slog.Int("count", len(ls)))

	// JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(ls)
	if err != nil {
		lh.logger.Error("Encode Error", slog.Any("error", err))
		http.Error(w, "Failed to encode listings", http.StatusInternalServerError)
		return
	}
}

func (lh *ListingHandlers) DeleteListing(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()
	requestId := middleware.GetRequestIdFromContext(ctx)
	
	// Delete listing
	result, err := lh.db.ExecContext(ctx, `DELETE FROM listing WHERE id = $1`, id)
	if err != nil {
		lh.logger.Error("Failed to delete listing", slog.Any("error", err), "listing_id", id, "request_id", requestId)
		http.Error(w, "Failed to delete listing", http.StatusInternalServerError)
		return
	} 

	// Get affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		lh.logger.Error("Failed to get affected rows", slog.Any("error", err), "request_id", requestId)
		http.Error(w, "Failed to get affected rows", http.StatusInternalServerError)
		return
	} 
	if rowsAffected == 0 {
		lh.logger.Error("Listing not found")
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
