package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/digino/pikro/internal/cloud"
	"github.com/digino/pikro/internal/config"
	"github.com/digino/pikro/internal/router"
)

// Pikro Cloud setup, run by a certified technician signed in to Cloud.

type cloudSessionView struct {
	Configured bool       `json:"configured"` // PIKRO_CLOUD_URL is set
	SignedIn   bool       `json:"signedIn"`
	Email      string     `json:"email,omitempty"`
	Name       string     `json:"name,omitempty"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
	Revoked    bool       `json:"revoked,omitempty"` // signed in, but no longer certified
}

func sessionView(s *config.CloudSession) cloudSessionView {
	v := cloudSessionView{Configured: cloud.URL() != ""}
	if s != nil && time.Now().Before(s.ExpiresAt) {
		v.SignedIn, v.Email, v.Name, v.ExpiresAt = true, s.Email, s.Name, &s.ExpiresAt
	}
	return v
}

func saveCloudSession(s *config.CloudSession) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	cfg.Cloud = s
	return config.Save(cfg)
}

// cloudError answers with Cloud's own message and status, or 502 when Cloud
// couldn't be reached.
func cloudError(w http.ResponseWriter, err error) {
	var e *cloud.Error
	if errors.As(err, &e) && e.Status != 0 {
		jsonError(w, e.Message, e.Status)
		return
	}
	jsonError(w, err.Error(), http.StatusBadGateway)
}

// GetCloudSession reports the signed-in technician, checked with Cloud when
// it can be reached: an expired session is forgotten, a revoked one flagged.
func GetCloudSession(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	v := sessionView(cfg.Cloud)
	if v.SignedIn {
		p, err := cloud.New().Me(r.Context(), cfg.Cloud.Token)
		var e *cloud.Error
		switch {
		case errors.As(err, &e) && e.Status == http.StatusUnauthorized:
			saveCloudSession(nil)
			v = sessionView(nil)
		case err == nil && p.Technician == nil:
			v.Revoked = true
		}
		// Any other error: Cloud unreachable, keep what we know.
	}
	jsonOK(w, v)
}

func PostCloudCode(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := cloud.New().SendCode(r.Context(), body.Email); err != nil {
		cloudError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PostCloudSession(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	s, err := cloud.New().CreateSession(r.Context(), body.Email, body.Code)
	if err != nil {
		cloudError(w, err)
		return
	}
	saved := &config.CloudSession{Token: s.Token, Email: s.Person.Email, ExpiresAt: s.ExpiresAt}
	if t := s.Person.Technician; t != nil {
		saved.Name = t.FirstName + " " + t.LastName
	}
	if err := saveCloudSession(saved); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	v := sessionView(saved)
	v.Revoked = s.Person.Technician == nil
	jsonOK(w, v)
}

// DeleteCloudSession signs the technician out: on Cloud when it can be
// reached (best effort), and always on this computer.
func DeleteCloudSession(w http.ResponseWriter, r *http.Request) {
	if cfg, err := config.Load(); err == nil && cfg.Cloud != nil {
		cloud.New().SignOut(r.Context(), cfg.Cloud.Token)
	}
	if err := saveCloudSession(nil); err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonOK(w, sessionView(nil))
}

type cloudPlanOption struct {
	Profile string `json:"profile"`
	Price   int64  `json:"price"` // 0: no whole-unit price set in Pikro, can't be sold online
}

// cloudPlanOptions lists the router's hotspot profiles with the price Pikro
// stores for each. Cloud takes whole currency units only.
func cloudPlanOptions(profiles []string, metas map[string]config.ProfileMeta) []cloudPlanOption {
	out := make([]cloudPlanOption, len(profiles))
	for i, name := range profiles {
		out[i].Profile = name
		if p, err := strconv.ParseInt(metas[name].Price, 10, 64); err == nil && p > 0 {
			out[i].Price = p
		}
	}
	return out
}

// routerCloudContext is what both the status and the connect handlers need.
func routerCloudContext(r *http.Request) (*router.Client, *config.HotspotSettings, []cloudPlanOption, error) {
	profile, err := findProfile(r.PathValue("id"))
	if err != nil {
		return nil, nil, nil, err
	}
	client := router.NewClient(profile.Host, profile.Port, profile.Username, profile.Password, profile.UseTLS)
	settings := profile.HotspotSettings
	if settings == nil {
		settings = &config.HotspotSettings{}
	}
	profiles, err := client.HotspotProfiles()
	if err != nil {
		return nil, nil, nil, err
	}
	names := make([]string, len(profiles))
	for i, p := range profiles {
		names[i], _ = p["name"].(string)
	}
	return client, settings, cloudPlanOptions(names, settings.ProfileMetas), nil
}

// GetRouterCloud reports whether the router is connected to Cloud and the
// plans it could sell online.
func GetRouterCloud(w http.ResponseWriter, r *http.Request) {
	client, settings, plans, err := routerCloudContext(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	status, err := client.CloudStatus(cloud.URL())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]any{"status": status, "plans": plans, "currency": settings.Currency})
}

// PostRouterCloud connects the router to Cloud: automatic time, Let's Encrypt
// root, registration, sync script + scheduler + walled garden, buy button on
// the login page, first sync.
// Stops at the first failed step; every step is reported.
func PostRouterCloud(w http.ResponseWriter, r *http.Request) {
	var body struct {
		OwnerPhone string   `json:"ownerPhone"`
		Profiles   []string `json:"profiles"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid body", http.StatusBadRequest)
		return
	}
	cfg, err := config.Load()
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session := sessionView(cfg.Cloud)
	if !session.SignedIn {
		jsonError(w, "sign in to Pikro Cloud first", http.StatusUnauthorized)
		return
	}
	client, settings, options, err := routerCloudContext(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	reg := cloud.Registration{OwnerPhone: body.OwnerPhone, Name: settings.HotspotName}
	for _, o := range options {
		if !slices.Contains(body.Profiles, o.Profile) {
			continue
		}
		if o.Price == 0 {
			jsonError(w, "profile "+o.Profile+" has no price: set one in Profiles first", http.StatusBadRequest)
			return
		}
		reg.Plans = append(reg.Plans, cloud.Plan{Label: o.Profile, Price: o.Price, Currency: settings.Currency, Profile: o.Profile})
	}
	if len(reg.Plans) != len(body.Profiles) {
		jsonError(w, "unknown hotspot profile", http.StatusBadRequest)
		return
	}

	var steps []router.SetupStepResult
	done := func(vouchers int, cloudRouterID string) {
		jsonOK(w, map[string]any{"steps": steps, "vouchers": vouchers, "cloudRouterId": cloudRouterID})
	}
	step := func(s router.SetupStepResult) bool {
		steps = append(steps, s)
		return s.Error == ""
	}

	clock := router.SetupStepResult{Name: "Turn on automatic time", OK: true}
	if err := client.EnableAutoTime(time.Now().UTC()); err != nil {
		clock = router.SetupStepResult{Name: clock.Name, Error: err.Error()}
	}
	if !step(clock) || !step(client.CloudImportCert()) {
		done(0, "")
		return
	}

	status, err := client.CloudStatus(cloud.URL())
	if err == nil && status.HardwareID == "" {
		err = errors.New("could not read the router's serial number")
	}
	if err == nil {
		reg.HardwareID, reg.RouterOSVersion = status.HardwareID, status.RouterOSVersion
		reg.RouterToken, err = client.CloudToken()
	}
	var registered *cloud.Registered
	if err == nil {
		registered, err = cloud.New().RegisterRouter(r.Context(), cfg.Cloud.Token, reg)
	}
	regStep := router.SetupStepResult{Name: "Register with Pikro Cloud", OK: err == nil}
	if err != nil {
		regStep.Error = err.Error()
		var e *cloud.Error
		if errors.As(err, &e) {
			regStep.Error = e.Message
		}
	}
	if !step(regStep) {
		done(0, "")
		return
	}

	cloudURL, _ := url.Parse(cloud.URL())
	// A warning on the buy button doesn't stop the connect: the page stays as it is.
	if !step(client.CloudInstall(registered.Script, cloudURL.Hostname())) || !step(client.CloudLoginPage(cloud.URL())) {
		done(0, registered.RouterID)
		return
	}
	first, vouchers := client.CloudFirstSync()
	step(first)
	done(vouchers, registered.RouterID)
}

// PostRouterCloudButton (re-)adds the buy button to a connected router's login page.
func PostRouterCloudButton(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonOK(w, client.CloudLoginPage(cloud.URL()))
}

// DeleteRouterCloud takes the router off Cloud: first tells Cloud to stop
// selling for it (with the router's token, read from its script before the
// script goes), then cleans the router. If Cloud can't be reached the router
// is cleaned anyway and cloudNotified is false.
func DeleteRouterCloud(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	token, err := client.CloudToken()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	notified := true
	if token != "" {
		err := cloud.New().Disconnect(r.Context(), token)
		var e *cloud.Error
		// 401: Cloud already doesn't accept this token (disconnected before, or
		// set up again elsewhere) — nothing left to stop.
		if err != nil && !(errors.As(err, &e) && e.Status == http.StatusUnauthorized) {
			log.Printf("cloud: disconnect: %v", err)
			notified = false
		}
	}
	if err := client.CloudRemove(); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	if s := client.CloudLoginPage(""); s.Error != "" {
		jsonError(w, s.Error, http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]bool{"cloudNotified": notified})
}

// withoutCloudVouchers drops unsold Pikro Cloud vouchers from a user list,
// so they aren't counted, edited or printed in Pikro. Once sold and used, the on-login script prefixes the
// comment, so the voucher shows like any other. Returns how many it hid.
func withoutCloudVouchers(users []map[string]any) ([]map[string]any, int) {
	shown := make([]map[string]any, 0, len(users))
	for _, u := range users {
		if c, _ := u["comment"].(string); c == router.CloudComment {
			continue
		}
		shown = append(shown, u)
	}
	return shown, len(users) - len(shown)
}

// refuseCloudVoucher answers 409 when userID is an unsold Pikro Cloud voucher.
func refuseCloudVoucher(w http.ResponseWriter, client *router.Client, userID string) bool {
	pool, err := client.IsCloudVoucher(userID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return true
	}
	if pool {
		jsonError(w, "this voucher is sold online by Pikro Cloud and can't be changed here", http.StatusConflict)
		return true
	}
	return false
}
