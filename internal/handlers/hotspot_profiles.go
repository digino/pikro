package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/digino/pikro/internal/config"
	"github.com/digino/pikro/internal/router"
)

func CreateHotspotProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	var body profileBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Name == "" {
		jsonError(w, "name is required", http.StatusBadRequest)
		return
	}
	result, err := client.CreateHotspotProfile(body.toParams())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	if err := saveProfileMeta(id, body.Name, body.Validity, body.Price); err != nil {
		// Non-fatal: profile created on router, just metadata save failed
		fmt.Printf("warn: failed to save profile meta for %s: %v\n", body.Name, err)
	}
	jsonCreated(w, result)
}

func UpdateHotspotProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	profileID := r.PathValue("profileID")
	var body profileBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := client.UpdateHotspotProfile(profileID, body.toParams()); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	if err := saveProfileMeta(id, body.Name, body.Validity, body.Price); err != nil {
		fmt.Printf("warn: failed to save profile meta for %s: %v\n", body.Name, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteHotspotProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	profileID := r.PathValue("profileID")
	// We need the profile name to delete its metadata — read it from body if provided,
	// otherwise skip meta cleanup (leftover metadata is harmless).
	var body struct{ Name string `json:"name"` }
	json.NewDecoder(r.Body).Decode(&body)

	if err := client.DeleteHotspotProfile(profileID); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	if body.Name != "" {
		deleteProfileMeta(id, body.Name)
	}
	w.WriteHeader(http.StatusNoContent)
}

// saveProfileMeta persists validity and price for a profile name into local config.
func saveProfileMeta(routerID, profileName, validity, price string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	for i := range cfg.Routers {
		if cfg.Routers[i].ID != routerID {
			continue
		}
		if cfg.Routers[i].HotspotSettings == nil {
			cfg.Routers[i].HotspotSettings = &config.HotspotSettings{}
		}
		s := cfg.Routers[i].HotspotSettings
		if s.ProfileMetas == nil {
			s.ProfileMetas = make(map[string]config.ProfileMeta)
		}
		s.ProfileMetas[profileName] = config.ProfileMeta{Validity: validity, Price: price}
		return config.Save(cfg)
	}
	return fmt.Errorf("router not found")
}

// deleteProfileMeta removes the metadata entry for a deleted profile.
func deleteProfileMeta(routerID, profileName string) {
	cfg, err := config.Load()
	if err != nil {
		return
	}
	for i := range cfg.Routers {
		if cfg.Routers[i].ID != routerID {
			continue
		}
		if cfg.Routers[i].HotspotSettings == nil || cfg.Routers[i].HotspotSettings.ProfileMetas == nil {
			return
		}
		delete(cfg.Routers[i].HotspotSettings.ProfileMetas, profileName)
		config.Save(cfg)
		return
	}
}

type profileBody struct {
	Name        string      `json:"name"`
	AddressPool string      `json:"addressPool"`
	SharedUsers interface{} `json:"sharedUsers"`
	RateLimit   string      `json:"rateLimit"`
	// Local metadata — stored in config, NOT sent to RouterOS
	Validity string `json:"validity"`
	Price    string `json:"price"`
}

func (b profileBody) toParams() router.HotspotProfileParams {
	su := ""
	switch v := b.SharedUsers.(type) {
	case string:
		su = v
	case float64:
		if v > 0 {
			su = fmt.Sprintf("%d", int(v))
		}
	}
	return router.HotspotProfileParams{
		Name:         b.Name,
		AddressPool:  b.AddressPool,
		SharedUsers:  su,
		RateLimit:    b.RateLimit,
		ValiditySecs: validityToSecs(b.Validity),
	}
}

// validityToSecs converts a shorthand validity string (e.g. "1d", "2h30m", "1w")
// to seconds. Returns 0 if blank or unparseable.
func validityToSecs(s string) int64 {
	if s == "" {
		return 0
	}
	var total int64
	re := regexp.MustCompile(`(\d+)([wdhm])`)
	for _, m := range re.FindAllStringSubmatch(s, -1) {
		n, _ := strconv.ParseInt(m[1], 10, 64)
		switch m[2] {
		case "w":
			total += n * 7 * 86400
		case "d":
			total += n * 86400
		case "h":
			total += n * 3600
		case "m":
			total += n * 60
		}
	}
	return total
}
