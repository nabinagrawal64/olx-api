package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/nabinagrawal64/olx-api/internal/httpx"
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
	requestId := middleware.GetRequestIdFromContext(ctx)
	// Query to fetch all listings
	rows, err := lh.db.QueryContext(ctx, `
			SELECT id, title, price, description, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100;
		`)
	if err != nil {
		lh.logger.Error("Query Listing Error", slog.Any("error", err), "request_id", requestId)
		httpx.Error(w, http.StatusInternalServerError, "Failed to query listings", httpx.CodeInternalError)
		return
	}
	defer rows.Close()

	// Iterate and store in slice
	var ls []listings
	for rows.Next() {
		var l listings
		if err := rows.Scan(&l.ID, &l.Title, &l.Price, &l.Description, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("Scan Error", slog.Any("error", err), "request_id", requestId)
			httpx.Error(w, http.StatusInternalServerError, "Failed to scan listings", httpx.CodeInternalError)
			return
		}
		ls = append(ls, l)
	}
	if err := rows.Err(); err != nil {
		lh.logger.Error("Rows Error", slog.Any("error", err), "request_id", requestId)
		httpx.Error(w, http.StatusInternalServerError, "Failed to scan listings", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("Fetched listings", slog.Int("count", len(ls)))

	// JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(ls)
	if err != nil {
		lh.logger.Error("Encode Error", slog.Any("error", err), "request_id", requestId)
		httpx.Error(w, http.StatusInternalServerError, "Failed to encode listings", httpx.CodeInternalError)
		return
	}
}

func (lh *ListingHandlers) DeleteListing(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()
	requestId := middleware.GetRequestIdFromContext(ctx)
	
	// Delete listing
	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		lh.logger.Error("Failed to delete listing", slog.Any("error", err), "listing_id", id, "request_id", requestId)
		// http.Error(w, "Failed to delete listing", http.StatusInternalServerError)
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete listing", httpx.CodeInternalError)
		return
	} 

	// Get affected rows
	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	lh.logger.Error("Failed to get affected rows", slog.Any("error", err), "request_id", requestId)
	// 	http.Error(w, "Failed to get affected rows", http.StatusInternalServerError)
	// 	return
	// } 
	// if rowsAffected == 0 {
	// 	lh.logger.Error("Listing not found")
	// 	http.Error(w, "Listing not found", http.StatusNotFound)
	// 	return
	// }

	// JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Listing deleted successfully",
	})
}

func (lh *ListingHandlers) CreateListing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.GetRequestIdFromContext(ctx)

	// decode the request body into a listing struct
	var req CreateListingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		lh.logger.Error("Failed to decode request body", slog.Any("error", err), "request_id", requestId)
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", httpx.CodeBadRequest)
		return
	}

	// validate the request body
	if err := req.Validate(); err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		lh.logger.Error("Invalid request body", slog.Any("request", req), "request_id", requestId)
		httpx.ValidationError(w, http.StatusBadRequest, err.Error(), httpx.CodeBadRequest, verr.Field)
		return 
	}

	// insert the listing into the database
	var out CreateListingResponse
	err = lh.db.QueryRowContext(ctx, 
		`INSERT INTO listings (title, price, description, city) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, title, created_at`,
		req.Title, req.Price, req.Description, req.City,
	).Scan(&out.ID, &out.Title, &out.CreatedAt)
	if err != nil {
		lh.logger.Error("Failed to create listing", slog.Any("error", err), "request_id", requestId)
		httpx.Error(w, http.StatusInternalServerError, "Failed to create listing", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("Listing created", slog.Any("listing", req), "request_id", requestId, slog.String("listing_id", out.ID))

	// JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Listing created successfully",
		"listing": out,
	})
}