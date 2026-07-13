package handlers

import (
	"encoding/json"
	"net/http"
)

func RunSpeedTest(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	var body struct {
		Target string `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Target == "" {
		body.Target = "8.8.8.8"
	}
	result, err := client.BandwidthTest(body.Target)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, result)
}
