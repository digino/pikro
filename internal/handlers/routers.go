package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/digino/pikro/internal/config"
	"github.com/digino/pikro/internal/router"
	"github.com/google/uuid"
)

func ListRouters(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Never send passwords to the frontend
	type safeProfile struct {
		ID              string                  `json:"id"`
		Name            string                  `json:"name"`
		Host            string                  `json:"host"`
		Port            int                     `json:"port"`
		Username        string                  `json:"username"`
		UseTLS          bool                    `json:"useTls"`
		HotspotSettings *config.HotspotSettings `json:"hotspotSettings,omitempty"`
	}
	safe := make([]safeProfile, len(cfg.Routers))
	for i, p := range cfg.Routers {
		safe[i] = safeProfile{p.ID, p.Name, p.Host, p.Port, p.Username, p.UseTLS, p.HotspotSettings}
	}
	jsonOK(w, safe)
}

func AddRouter(w http.ResponseWriter, r *http.Request) {
	var input config.RouterProfile
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if input.Host == "" || input.Username == "" {
		jsonError(w, "host and username are required", http.StatusBadRequest)
		return
	}
	if input.Port == 0 {
		input.Port = 8728
	}
	input.ID = uuid.NewString()

	cfg, err := config.Load()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cfg.Routers = append(cfg.Routers, input)
	if err := config.Save(cfg); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, map[string]string{"id": input.ID})
}

func UpdateRouter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input struct {
		Name            string `json:"name"`
		Host            string `json:"host"`
		Port            int    `json:"port"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		UseTLS          bool   `json:"useTls"`
		HotspotSettings *struct {
			HotspotName string `json:"hotspotName"`
			DNSName     string `json:"dnsName"`
			Currency    string `json:"currency"`
		} `json:"hotspotSettings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	cfg, err := config.Load()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	found := false
	for i := range cfg.Routers {
		if cfg.Routers[i].ID != id {
			continue
		}
		found = true
		if input.Name != "" {
			cfg.Routers[i].Name = input.Name
		}
		if input.Host != "" {
			cfg.Routers[i].Host = input.Host
		}
		if input.Port != 0 {
			cfg.Routers[i].Port = input.Port
		}
		if input.Username != "" {
			cfg.Routers[i].Username = input.Username
		}
		if input.Password != "" {
			cfg.Routers[i].Password = input.Password
		}
		cfg.Routers[i].UseTLS = input.UseTLS
		// Merge only the fields the router dialog owns — ProfileMetas, Voucher,
		// and LoginPage live on the same struct but are managed elsewhere
		// (Profiles page, Settings tabs) and must not be clobbered.
		if input.HotspotSettings != nil {
			hs := cfg.Routers[i].HotspotSettings
			if hs == nil {
				hs = &config.HotspotSettings{}
				cfg.Routers[i].HotspotSettings = hs
			}
			hs.HotspotName = input.HotspotSettings.HotspotName
			hs.DNSName = input.HotspotSettings.DNSName
			hs.Currency = input.HotspotSettings.Currency
		}
		break
	}
	if !found {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	if err := config.Save(cfg); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteRouter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cfg, err := config.Load()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var filtered []config.RouterProfile
	for _, p := range cfg.Routers {
		if p.ID != id {
			filtered = append(filtered, p)
		}
	}
	if filtered == nil {
		filtered = []config.RouterProfile{}
	}
	cfg.Routers = filtered
	if err := config.Save(cfg); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func TestRouter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profile, err := findProfile(id)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	client := router.NewClient(profile.Host, profile.Port, profile.Username, profile.Password, profile.UseTLS)
	if err := client.Ping(); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]bool{"ok": true})
}

func findProfile(id string) (*config.RouterProfile, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	for i := range cfg.Routers {
		if cfg.Routers[i].ID == id {
			return &cfg.Routers[i], nil
		}
	}
	return nil, http.ErrNoCookie // reuse as "not found" sentinel
}
