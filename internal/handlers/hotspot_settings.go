package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/digino/pikro/internal/config"
	"github.com/digino/pikro/internal/router"
)

func GetHotspotSettings(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profile, err := findProfile(id)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	if profile.HotspotSettings == nil {
		jsonOK(w, config.HotspotSettings{})
		return
	}
	jsonOK(w, profile.HotspotSettings)
}

// GetProfileMetas returns the locally-stored metadata (validity, price) for all profiles.
func GetProfileMetas(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profile, err := findProfile(id)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	if profile.HotspotSettings == nil || profile.HotspotSettings.ProfileMetas == nil {
		jsonOK(w, map[string]config.ProfileMeta{})
		return
	}
	jsonOK(w, profile.HotspotSettings.ProfileMetas)
}

func UploadLoginPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input struct {
		HTML string `json:"html"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}

	profile, err := findProfile(id)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}

	client := router.NewClient(profile.Host, profile.Port, profile.Username, profile.Password, profile.UseTLS)

	if err := client.UploadLoginPage("pikro-profile", input.HTML); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]bool{"ok": true})
}

// GetLoginPageHTML returns the actual contents of hotspot/login.html as
// currently stored on the router — what's really served to devices, which
// may differ from the template selected in Pikro's local settings if the
// router was edited outside Pikro or never had a template uploaded yet.
func GetLoginPageHTML(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profile, err := findProfile(id)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}

	client := router.NewClient(profile.Host, profile.Port, profile.Username, profile.Password, profile.UseTLS)

	html, err := client.GetLoginPageHTML()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]string{"html": html})
}

func PutHotspotSettings(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input config.HotspotSettings
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
		if cfg.Routers[i].ID == id {
			cfg.Routers[i].HotspotSettings = &input
			found = true
			break
		}
	}
	if !found {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	if err := config.Save(cfg); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, input)
}
