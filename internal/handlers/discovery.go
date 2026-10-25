package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/digino/pikro/internal/discovery"
)

func DiscoverRouters(w http.ResponseWriter, r *http.Request) {
	log.Println("[discovery] scan started")

	devices, err := discovery.Scan(3 * time.Second)
	if errors.Is(err, discovery.ErrLocalNetworkDenied) {
		log.Printf("[discovery] scan error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "discovery failed: " + err.Error(),
			"code":  "local_network_denied",
		})
		return
	}
	if err != nil {
		log.Printf("[discovery] scan error: %v", err)
		jsonError(w, "discovery failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[discovery] scan complete — %d device(s) found", len(devices))
	for i, d := range devices {
		log.Printf("[discovery]   [%d] ip=%s mac=%s identity=%q board=%q version=%q",
			i, d.IP, d.MAC, d.Identity, d.Board, d.Version)
	}

	if devices == nil {
		devices = []discovery.Device{}
	}
	jsonOK(w, devices)
}
