// Package sales manages the local generation ledger.
// One file per router per year: ~/.pikro-sales-<routerID>-<year>.json
// Each file is a JSON array of SaleEntry appended on every batch creation.
// Files are append-only and never rewritten — old years are naturally archivable.
package sales

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SaleEntry records a single voucher generation event.
type SaleEntry struct {
	At       time.Time `json:"at"`       // UTC creation time
	Profile  string    `json:"profile"`
	Price    string    `json:"price"`    // as stored in ProfileMeta, e.g. "500"
	Currency string    `json:"currency"` // e.g. "XOF"
	Count    int       `json:"count"`    // number of vouchers in this batch
}

var mu sync.Mutex

func ledgerPath(routerID string, year int) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	name := filepath.Join(home, ".pikro-sales-"+routerID+"-"+itoa(year)+".json")
	return name, nil
}

// Append adds entry to the ledger for the entry's year.
// The file is created if it does not exist. Safe for concurrent use.
func Append(routerID string, entry SaleEntry) error {
	mu.Lock()
	defer mu.Unlock()

	path, err := ledgerPath(routerID, entry.At.Year())
	if err != nil {
		return err
	}

	// Read existing entries (file may not exist yet).
	var entries []SaleEntry
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &entries) // ignore parse errors on corrupt file
	}

	entries = append(entries, entry)

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// Load returns all entries for the given router and year.
// Returns an empty slice (not an error) if the file does not exist yet.
func Load(routerID string, year int) ([]SaleEntry, error) {
	path, err := ledgerPath(routerID, year)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []SaleEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []SaleEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func itoa(n int) string {
	return time.Date(n, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006")
}
