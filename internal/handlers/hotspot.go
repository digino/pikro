package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/digino/pikro/internal/config"
	"github.com/digino/pikro/internal/router"
	"github.com/digino/pikro/internal/sales"
)

const maxBatchUsers = 100

func ListHotspotUsers(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	users, err := client.HotspotUsers()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	if err := sales.SyncActivations(r.PathValue("id"), users, profileValidities(r.PathValue("id")), time.Now()); err != nil {
		log.Printf("sales: sync activations: %v", err)
	}
	jsonOK(w, users)
}

// profileValidities returns each profile's validity in seconds from local config.
func profileValidities(routerID string) map[string]int64 {
	out := map[string]int64{}
	cfg, err := config.Load()
	if err != nil {
		return out
	}
	for _, rp := range cfg.Routers {
		if rp.ID != routerID || rp.HotspotSettings == nil {
			continue
		}
		for name, meta := range rp.HotspotSettings.ProfileMetas {
			out[name] = config.ValiditySecs(meta.Validity)
		}
	}
	return out
}

func ListHotspotActive(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	active, err := client.HotspotActive()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, active)
}

func ListHotspotHosts(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	hosts, err := client.HotspotHosts()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, hosts)
}

func ListHotspotProfiles(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	profiles, err := client.HotspotProfiles()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, profiles)
}

func ListAddressPools(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	pools, err := client.AddressPools()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, pools)
}

func CreateHotspotUser(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	var body struct {
		Name            string `json:"name"`
		Password        string `json:"password"`
		Profile         string `json:"profile"`
		LimitUptime     string `json:"limitUptime"`
		LimitBytesTotal string `json:"limitBytesTotal"`
		RateLimit       string `json:"rateLimit"`
		Comment         string `json:"comment"`
		// ExpiryComment is pre-formatted by the frontend as "exp:YYYY-MM-DD HH:MM:SS"
		// using the router's current time + profile validity duration.
		ExpiryComment string `json:"expiryComment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Name == "" {
		jsonError(w, "name is required", http.StatusBadRequest)
		return
	}
	// Prepend expiry tag to comment so the cleanup script can parse it.
	comment := body.Comment
	if body.ExpiryComment != "" {
		if comment != "" {
			comment = body.ExpiryComment + " " + comment
		} else {
			comment = body.ExpiryComment
		}
	}
	user, err := client.CreateHotspotUser(router.HotspotUserParams{
		Name:            body.Name,
		Password:        body.Password,
		Profile:         body.Profile,
		LimitUptime:     body.LimitUptime,
		LimitBytesTotal: body.LimitBytesTotal,
		RateLimit:       body.RateLimit,
		Comment:         comment,
	})
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonOK(w, user)
}

func ToggleHotspotUser(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	userID := r.PathValue("userID")
	var body struct {
		Disabled bool `json:"disabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := client.ToggleHotspotUser(userID, body.Disabled); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateHotspotUser(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	userID := r.PathValue("userID")
	var body struct {
		Password        string `json:"password"`
		Profile         string `json:"profile"`
		LimitUptime     string `json:"limitUptime"`
		LimitBytesTotal string `json:"limitBytesTotal"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := client.UpdateHotspotUser(userID, router.UpdateHotspotUserParams{
		Password:        body.Password,
		Profile:         body.Profile,
		LimitUptime:     body.LimitUptime,
		LimitBytesTotal: body.LimitBytesTotal,
	}); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteHotspotUser(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	userID := r.PathValue("userID")
	if err := client.DeleteHotspotUser(userID); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DisconnectHotspotActive kicks a live session (by its /ip/hotspot/active
// .id) without touching the underlying /ip/hotspot/user record.
func DisconnectHotspotActive(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	sessionID := r.PathValue("sessionID")
	if err := client.DisconnectHotspotActive(sessionID); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func clientFromRequest(r *http.Request) (*router.Client, error) {
	routerID := r.PathValue("id")
	profile, err := findProfile(routerID)
	if err != nil {
		return nil, err
	}
	return router.NewClient(profile.Host, profile.Port, profile.Username, profile.Password, profile.UseTLS), nil
}

// CreateHotspotUsersBatch creates one chunk of a voucher batch and records the
// successfully created names in the sales ledger.
// POST /api/routers/{id}/hotspot/users/batch
func CreateHotspotUsersBatch(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	var body struct {
		BatchID         string `json:"batchId"`
		Profile         string `json:"profile"`
		LimitUptime     string `json:"limitUptime"`
		LimitBytesTotal string `json:"limitBytesTotal"`
		Comment         string `json:"comment"`
		Price           string `json:"price"`
		Currency        string `json:"currency"`
		Users           []struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		} `json:"users"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.BatchID == "" || body.Profile == "" {
		jsonError(w, "batchId and profile are required", http.StatusBadRequest)
		return
	}
	if len(body.Users) < 1 || len(body.Users) > maxBatchUsers {
		jsonError(w, "users must contain between 1 and 100 entries", http.StatusBadRequest)
		return
	}
	params := make([]router.HotspotUserParams, len(body.Users))
	for i, u := range body.Users {
		params[i] = router.HotspotUserParams{
			Name:            u.Name,
			Password:        u.Password,
			Profile:         body.Profile,
			LimitUptime:     body.LimitUptime,
			LimitBytesTotal: body.LimitBytesTotal,
			Comment:         body.Comment,
		}
	}
	results, err := client.CreateHotspotUsers(params)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	var created []string
	for _, res := range results {
		if res.OK {
			created = append(created, res.Name)
		}
	}
	if len(created) > 0 {
		if err := sales.AppendBatch(r.PathValue("id"), body.BatchID, time.Now(), body.Profile, body.Price, body.Currency, created); err != nil {
			log.Printf("sales: append batch: %v", err)
		}
	}
	jsonOK(w, map[string]any{"results": results})
}
