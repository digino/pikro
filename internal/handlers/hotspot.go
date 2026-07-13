package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/digino/pikro/internal/router"
)

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
	jsonOK(w, users)
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
		Comment         string `json:"comment"`
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
		Comment:         body.Comment,
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

func clientFromRequest(r *http.Request) (*router.Client, error) {
	routerID := r.PathValue("id")
	profile, err := findProfile(routerID)
	if err != nil {
		return nil, err
	}
	return router.NewClient(profile.Host, profile.Port, profile.Username, profile.Password, profile.UseTLS), nil
}
