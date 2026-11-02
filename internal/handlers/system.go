package handlers

import (
	"net/http"
	"sync"

	"github.com/digino/pikro/internal/router"
	"github.com/digino/pikro/internal/update"
)

// AppVersion is set by main via handlers.AppVersion = Version at startup.
var AppVersion = "dev"

func GetAppVersion(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]string{"version": AppVersion})
}

// Created on first use, once main has set AppVersion.
var updateChecker = sync.OnceValue(func() *update.Checker { return update.NewChecker(AppVersion) })

// GetUpdate says whether a newer Pikro release is out. A failed check is not
// an error for the page: it answers "nothing available" plus the reason.
func GetUpdate(w http.ResponseWriter, r *http.Request) {
	info, err := updateChecker().Check(r.Context(), r.URL.Query().Get("refresh") == "1")
	resp := struct {
		update.Info
		Error string `json:"error,omitempty"`
	}{Info: info}
	if err != nil {
		resp.Error = err.Error()
	}
	jsonOK(w, resp)
}

func GetSystemResource(w http.ResponseWriter, r *http.Request) {
	profile, err := findProfile(r.PathValue("id"))
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	client := router.NewClient(profile.Host, profile.Port, profile.Username, profile.Password, profile.UseTLS)
	resource, err := client.SystemResource()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	rememberMAC(profile, client)
	jsonOK(w, resource)
}

func GetPollSnapshot(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	snap, err := client.Poll()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, snap)
}

func GetWanIP(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	ip, err := client.WanIP()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, map[string]string{"ip": ip})
}

func GetInterfaceTraffic(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	traffic, err := client.InterfaceTraffic()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, traffic)
}

func GetSystemClock(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	clock, err := client.SystemClock()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, clock)
}

func GetSystemLogs(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	logs, err := client.SystemLogs(router.IsHotspotLog)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, logs)
}

func GetDHCPLeases(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	leases, err := client.DHCPLeases()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, leases)
}

func GetDHCPServers(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	servers, err := client.DHCPServers()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, servers)
}
