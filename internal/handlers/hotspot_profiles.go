package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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
	if err := body.validate(); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
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
	if err := body.validate(); err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
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
	var body struct {
		Name string `json:"name"`
	}
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

var (
	profilePriceRe     = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	profileValidityRe  = regexp.MustCompile(`^(\d+[wdhm])+$`)
	profileRateLimitRe = regexp.MustCompile(`^\d+[kKMG]?(/\d+[kKMG]?)?( \S+)*$`) // rx/tx, then RouterOS burst fields
)

// validate checks the form before anything reaches the router or the config.
func (b profileBody) validate() error {
	name := strings.TrimSpace(b.Name)
	switch {
	case name == "":
		return errors.New("name is required")

	// ";", and a quote would break the RouterOS command creating the voucher.
	case strings.ContainsAny(name, `;,"`):
		return errors.New(`name can't contain ; , or "`)
	case b.Price != "" && (!profilePriceRe.MatchString(b.Price) || strings.Trim(b.Price, "0.") == ""):
		return errors.New("price must be a positive number, e.g. 500 or 1.50")
	case b.Validity != "" && (!profileValidityRe.MatchString(b.Validity) || config.ValiditySecs(b.Validity) == 0):
		return errors.New("validity must look like 1d, 4h, 2h30m or 1w")
	case b.RateLimit != "" && !profileRateLimitRe.MatchString(b.RateLimit):
		return errors.New("speed limit must look like 2M/4M")
	}
	switch v := b.SharedUsers.(type) {
	case nil:
	case float64:
		if v < 1 || v != float64(int(v)) {
			return errors.New("devices per voucher must be a whole number, 1 or more")
		}
	case string:
		if n, err := strconv.Atoi(v); v != "" && v != "unlimited" && (err != nil || n < 1) {
			return errors.New("devices per voucher must be a whole number, 1 or more")
		}
	default:
		return errors.New("devices per voucher must be a number")
	}
	return nil
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
		ValiditySecs: config.ValiditySecs(b.Validity),
	}
}
