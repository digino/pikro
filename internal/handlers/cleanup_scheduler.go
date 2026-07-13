package handlers

import (
	"encoding/json"
	"net/http"
)

func GetCleanupScheduler(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	installed, interval, err := client.CleanupSchedulerStatus()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]any{"installed": installed, "interval": interval})
}

func PutCleanupScheduler(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	var body struct {
		Enabled  bool   `json:"enabled"`
		Interval string `json:"interval"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if !body.Enabled {
		if err := client.RemoveCleanupScheduler(); err != nil {
			jsonError(w, err.Error(), http.StatusBadGateway)
			return
		}
		jsonOK(w, map[string]any{"installed": false, "interval": ""})
		return
	}
	if body.Interval == "" {
		jsonError(w, "interval is required", http.StatusBadRequest)
		return
	}
	if err := client.InstallCleanupScheduler(body.Interval); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]any{"installed": true, "interval": body.Interval})
}
