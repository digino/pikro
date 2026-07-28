package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/digino/pikro/internal/handlers"
)

func Start(port int, webFS fs.FS) error {
	mux := http.NewServeMux()
	registerAPI(mux)
	registerSPA(mux, webFS)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), requestLogger(mux))
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only log API requests — static assets are noise
		if len(r.URL.Path) >= 5 && r.URL.Path[:5] == "/api/" {
			rw := &statusWriter{ResponseWriter: w, status: 200}
			start := time.Now()
			next.ServeHTTP(rw, r)
			log.Printf("%s %s %d (%s)", r.Method, r.URL.Path, rw.status, time.Since(start).Round(time.Millisecond))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

func registerAPI(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/version", handlers.GetAppVersion)
	mux.HandleFunc("GET /api/discover", handlers.DiscoverRouters)

	mux.HandleFunc("GET /api/routers", handlers.ListRouters)
	mux.HandleFunc("POST /api/routers", handlers.AddRouter)
	mux.HandleFunc("PATCH /api/routers/{id}", handlers.UpdateRouter)
	mux.HandleFunc("DELETE /api/routers/{id}", handlers.DeleteRouter)
	mux.HandleFunc("GET /api/routers/{id}/test", handlers.TestRouter)

	mux.HandleFunc("GET /api/routers/{id}/system/resource", handlers.GetSystemResource)
	mux.HandleFunc("GET /api/routers/{id}/system/poll", handlers.GetPollSnapshot)
	mux.HandleFunc("GET /api/routers/{id}/system/wan-ip", handlers.GetWanIP)
	mux.HandleFunc("GET /api/routers/{id}/system/traffic", handlers.GetInterfaceTraffic)
	mux.HandleFunc("GET /api/routers/{id}/system/clock", handlers.GetSystemClock)
	mux.HandleFunc("GET /api/routers/{id}/system/logs", handlers.GetSystemLogs)

	mux.HandleFunc("GET /api/routers/{id}/hotspot/preflight", handlers.HotspotPreflight)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/setup", handlers.SetupHotspot)
	mux.HandleFunc("DELETE /api/routers/{id}/hotspot/setup", handlers.TeardownHotspot)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/settings", handlers.GetHotspotSettings)
	mux.HandleFunc("PUT /api/routers/{id}/hotspot/settings", handlers.PutHotspotSettings)
	mux.HandleFunc("PUT /api/routers/{id}/hotspot/login-page", handlers.UploadLoginPage)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/profile-metas", handlers.GetProfileMetas)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/cleanup", handlers.GetCleanupScheduler)
	mux.HandleFunc("PUT /api/routers/{id}/hotspot/cleanup", handlers.PutCleanupScheduler)
	mux.HandleFunc("GET /api/routers/{id}/sales", handlers.GetSales)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/users", handlers.ListHotspotUsers)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/users", handlers.CreateHotspotUser)
	mux.HandleFunc("PATCH /api/routers/{id}/hotspot/users/{userID}", handlers.UpdateHotspotUser)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/users/{userID}/toggle", handlers.ToggleHotspotUser)
	mux.HandleFunc("DELETE /api/routers/{id}/hotspot/users/{userID}", handlers.DeleteHotspotUser)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/active", handlers.ListHotspotActive)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/profiles", handlers.ListHotspotProfiles)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/profiles", handlers.CreateHotspotProfile)
	mux.HandleFunc("PATCH /api/routers/{id}/hotspot/profiles/{profileID}", handlers.UpdateHotspotProfile)
	mux.HandleFunc("DELETE /api/routers/{id}/hotspot/profiles/{profileID}", handlers.DeleteHotspotProfile)

	mux.HandleFunc("POST /api/routers/{id}/speedtest", handlers.RunSpeedTest)
}

func registerSPA(mux *http.ServeMux, webFS fs.FS) {
	if webFS == nil {
		return
	}
	sub, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		return
	}
	fileServer := http.FileServer(http.FS(sub))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= 5 && r.URL.Path[:5] == "/api/" {
			http.NotFound(w, r)
			return
		}
		if _, err := sub.Open(r.URL.Path[1:]); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
