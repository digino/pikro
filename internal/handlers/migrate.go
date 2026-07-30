package handlers

import "net/http"

// PostMigrateMikhmon converts Mikhmon-style hotspot user comments and
// profile on-login scripts to Pikro's own exp:<epoch> convention. Safe to
// call multiple times — already-migrated users/profiles are left untouched.
func PostMigrateMikhmon(w http.ResponseWriter, r *http.Request) {
	client, err := clientFromRequest(r)
	if err != nil {
		jsonError(w, "router not found", http.StatusNotFound)
		return
	}
	result, err := client.MigrateFromMikhmon()
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	jsonOK(w, result)
}
