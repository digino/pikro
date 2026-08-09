package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/digino/pikro/internal/router"
)

func HotspotPreflight(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	result, err := client.HotspotPreflight()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, result)
}

func TeardownHotspot(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	result, err := client.TeardownHotspot()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, result)
}

func CheckExistingDHCP(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	iface := r.URL.Query().Get("iface")
	if iface == "" {
		jsonError(w, "iface is required", http.StatusBadRequest)
		return
	}
	result, err := client.CheckExistingDHCP(iface)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonOK(w, result)
}

func SetupHotspot(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	var req router.SetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.LANIface == "" || req.WANIface == "" || req.Subnet == "" || req.HotspotName == "" {
		jsonError(w, "lanIface, wanIface, subnet and hotspotName are required", http.StatusBadRequest)
		return
	}
	result, err := client.SetupHotspot(req)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonOK(w, result)
}
