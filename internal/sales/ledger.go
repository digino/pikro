// Package sales manages the local sales ledger.
// One file per router per year: ~/.pikro-sales-<routerID>-<year>.jsonl
// Each line is a JSON event: a "batch" (vouchers generated) or an
// "activation" (a voucher of a batch was used for the first time).
// "Sold" means activated. Files are append-only and never rewritten.
// The pre-v2 ~/.pikro-sales-<routerID>-<year>.json arrays are no longer read:
// they lack voucher names, so their sales can never be confirmed.
package sales

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	typeBatch      = "batch"
	typeActivation = "activation"
)

// Batch is a group of generated vouchers sharing one id.
type Batch struct {
	ID       string    `json:"id"`
	At       time.Time `json:"at"`
	Profile  string    `json:"profile"`
	Price    string    `json:"price"`
	Currency string    `json:"currency"`
	Count    int       `json:"count"`
	Names    []string  `json:"-"`
}

// Activation records the first use of a voucher.
type Activation struct {
	Batch    string    `json:"batch"`
	Name     string    `json:"name"`
	At       time.Time `json:"at"`
	Profile  string    `json:"profile"`
	Price    string    `json:"price"`
	Currency string    `json:"currency"`
}

// Ledger is the result of Load.
type Ledger struct {
	Batches     []Batch      `json:"batches"`
	Activations []Activation `json:"activations"`
}

// line is the on-disk JSON Lines record for both event types.
type line struct {
	Type     string    `json:"type"`
	ID       string    `json:"id,omitempty"`
	Batch    string    `json:"batch,omitempty"`
	Name     string    `json:"name,omitempty"`
	At       time.Time `json:"at"`
	Profile  string    `json:"profile,omitempty"`
	Price    string    `json:"price,omitempty"`
	Currency string    `json:"currency,omitempty"`
	Names    []string  `json:"names,omitempty"`
}

var mu sync.Mutex

func basePath(routerID string, year int) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".pikro-sales-"+routerID+"-"+strconv.Itoa(year)), nil
}

// AppendBatch records one chunk of generated vouchers. Chunks sharing an id
// are merged on load. Safe for concurrent use.
func AppendBatch(routerID, id string, at time.Time, profile, price, currency string, names []string) error {
	at = at.UTC()
	mu.Lock()
	defer mu.Unlock()
	return appendLine(routerID, at.Year(), line{
		Type: typeBatch, ID: id, At: at,
		Profile: profile, Price: price, Currency: currency, Names: names,
	})
}

func appendLine(routerID string, year int, l line) error {
	base, err := basePath(routerID, year)
	if err != nil {
		return err
	}
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(base+".jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(data, '\n')); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// readYear returns the merged batches (with names) and the raw activations of one year.
// Missing files yield empty results. Unparseable lines are skipped.
func readYear(routerID string, year int) ([]Batch, []Activation, error) {
	base, err := basePath(routerID, year)
	if err != nil {
		return nil, nil, err
	}
	var batches []Batch
	index := map[string]int{}

	var acts []Activation
	data, err := os.ReadFile(base + ".jsonl")
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, err
	}
	for _, raw := range splitLines(data) {
		var l line
		if json.Unmarshal(raw, &l) != nil {
			continue
		}
		switch l.Type {
		case typeBatch:
			i, ok := index[l.ID]
			if !ok {
				index[l.ID] = len(batches)
				batches = append(batches, Batch{
					ID: l.ID, At: l.At.UTC(), Profile: l.Profile, Price: l.Price, Currency: l.Currency,
				})
				i = len(batches) - 1
			}
			b := &batches[i]
			if l.At.Before(b.At) {
				b.At = l.At.UTC()
			}
			b.Names = append(b.Names, l.Names...)
			b.Count = len(b.Names)
		case typeActivation:
			acts = append(acts, Activation{Batch: l.Batch, Name: l.Name, At: l.At.UTC()})
		}
	}
	return batches, acts, nil
}

func splitLines(data []byte) [][]byte {
	var out [][]byte
	start := 0
	for i, c := range data {
		if c == '\n' {
			if i > start {
				out = append(out, data[start:i])
			}
			start = i + 1
		}
	}
	if start < len(data) {
		out = append(out, data[start:])
	}
	return out
}

// Load returns the batches and activations recorded in the given year.
// Activation metadata (profile, price, currency) comes from its batch, which
// may live in the previous year. Missing files yield an empty ledger.
func Load(routerID string, year int) (Ledger, error) {
	mu.Lock()
	defer mu.Unlock()
	batches, acts, err := readYear(routerID, year)
	if err != nil {
		return Ledger{}, err
	}
	prev, _, err := readYear(routerID, year-1)
	if err != nil {
		return Ledger{}, err
	}
	byID := map[string]Batch{}
	for _, b := range append(prev, batches...) {
		byID[b.ID] = b
	}
	for i := range acts {
		b := byID[acts[i].Batch]
		acts[i].Profile, acts[i].Price, acts[i].Currency = b.Profile, b.Price, b.Currency
	}
	if batches == nil {
		batches = []Batch{}
	}
	if acts == nil {
		acts = []Activation{}
	}
	return Ledger{Batches: batches, Activations: acts}, nil
}
