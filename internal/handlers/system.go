package handlers

import "net/http"

// AppVersion is set by main via handlers.AppVersion = Version at startup.
var AppVersion = "dev"

func GetAppVersion(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]string{"version": AppVersion})
}

func GetSystemResource(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}
	resource, err := client.SystemResource()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, resource)
}
