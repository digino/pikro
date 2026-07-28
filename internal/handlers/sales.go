package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/digino/pikro/internal/sales"
)

// GetSales returns the generation ledger for a router and year.
// GET /api/routers/{id}/sales?year=2025
func GetSales(w http.ResponseWriter, r *http.Request) {
	routerID := r.PathValue("id")
	yearStr := r.URL.Query().Get("year")
	year := time.Now().Year()
	if y, err := strconv.Atoi(yearStr); err == nil && y > 2000 {
		year = y
	}
	entries, err := sales.Load(routerID, year)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, entries)
}
