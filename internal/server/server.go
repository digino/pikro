package server

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/digino/pikro/internal/handlers"
)

// Start serves Pikro to this computer only: it holds router passwords and
// the technician's Cloud session, so nothing else on the network may use it.
func Start(port int, webFS fs.FS) error {
	mux := http.NewServeMux()
	registerAPI(mux)
	registerSPA(mux, webFS)
	srv := &http.Server{Handler: localOnly(requestLogger(mux)), ReadHeaderTimeout: 10 * time.Second}

	// "localhost" may resolve to either loopback address, so listen on both.
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	if ln6, err := net.Listen("tcp", fmt.Sprintf("[::1]:%d", port)); err == nil {
		go srv.Serve(ln6)
	}
	return srv.Serve(ln)
}

// localOnly refuses requests not addressed to localhost (a website pointing
// its own name at 127.0.0.1) and requests sent from another site's page.
func localOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isLoopback(r.Host) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if origin := r.Header.Get("Origin"); origin != "" {
			u, err := url.Parse(origin)
			if err != nil || !isLoopback(u.Host) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func isLoopback(hostport string) bool {
	host := hostport
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		host = h
	}
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
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
	mux.HandleFunc("GET /api/update", handlers.GetUpdate)
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
	mux.HandleFunc("GET /api/routers/{id}/system/dhcp-leases", handlers.GetDHCPLeases)
	mux.HandleFunc("GET /api/routers/{id}/system/dhcp-servers", handlers.GetDHCPServers)

	mux.HandleFunc("GET /api/routers/{id}/hotspot/preflight", handlers.HotspotPreflight)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/dhcp-check", handlers.CheckExistingDHCP)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/setup", handlers.SetupHotspot)
	mux.HandleFunc("DELETE /api/routers/{id}/hotspot/setup", handlers.TeardownHotspot)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/settings", handlers.GetHotspotSettings)
	mux.HandleFunc("PUT /api/routers/{id}/hotspot/settings", handlers.PutHotspotSettings)
	mux.HandleFunc("PUT /api/routers/{id}/hotspot/login-page", handlers.UploadLoginPage)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/login-page", handlers.GetLoginPageHTML)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/profile-metas", handlers.GetProfileMetas)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/cleanup", handlers.GetCleanupScheduler)
	mux.HandleFunc("PUT /api/routers/{id}/hotspot/cleanup", handlers.PutCleanupScheduler)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/migrate-mikhmon", handlers.PostMigrateMikhmon)
	mux.HandleFunc("GET /api/routers/{id}/sales", handlers.GetSales)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/users", handlers.ListHotspotUsers)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/users", handlers.CreateHotspotUser)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/users/batch", handlers.CreateHotspotUsersBatch)
	mux.HandleFunc("PATCH /api/routers/{id}/hotspot/users/{userID}", handlers.UpdateHotspotUser)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/users/{userID}/toggle", handlers.ToggleHotspotUser)
	mux.HandleFunc("DELETE /api/routers/{id}/hotspot/users/{userID}", handlers.DeleteHotspotUser)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/active", handlers.ListHotspotActive)
	mux.HandleFunc("DELETE /api/routers/{id}/hotspot/active/{sessionID}", handlers.DisconnectHotspotActive)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/hosts", handlers.ListHotspotHosts)
	mux.HandleFunc("GET /api/routers/{id}/hotspot/profiles", handlers.ListHotspotProfiles)
	mux.HandleFunc("POST /api/routers/{id}/hotspot/profiles", handlers.CreateHotspotProfile)
	mux.HandleFunc("PATCH /api/routers/{id}/hotspot/profiles/{profileID}", handlers.UpdateHotspotProfile)
	mux.HandleFunc("DELETE /api/routers/{id}/hotspot/profiles/{profileID}", handlers.DeleteHotspotProfile)
	mux.HandleFunc("GET /api/routers/{id}/address-pools", handlers.ListAddressPools)

	mux.HandleFunc("GET /api/cloud/session", handlers.GetCloudSession)
	mux.HandleFunc("POST /api/cloud/code", handlers.PostCloudCode)
	mux.HandleFunc("POST /api/cloud/session", handlers.PostCloudSession)
	mux.HandleFunc("DELETE /api/cloud/session", handlers.DeleteCloudSession)
	mux.HandleFunc("GET /api/routers/{id}/cloud", handlers.GetRouterCloud)
	mux.HandleFunc("POST /api/routers/{id}/cloud", handlers.PostRouterCloud)
	mux.HandleFunc("DELETE /api/routers/{id}/cloud", handlers.DeleteRouterCloud)
	mux.HandleFunc("POST /api/routers/{id}/cloud/button", handlers.PostRouterCloudButton)

	mux.HandleFunc("POST /api/routers/{id}/speedtest", handlers.RunSpeedTest)
	mux.HandleFunc("GET /api/routers/{id}/tools/internet-check", handlers.GetInternetCheck)
	mux.HandleFunc("POST /api/routers/{id}/tools/public-dns", handlers.PostPublicDNS)
	mux.HandleFunc("GET /api/routers/{id}/tools/clock", handlers.GetRouterClock)
	mux.HandleFunc("POST /api/routers/{id}/tools/clock/auto", handlers.PostAutoTime)
	mux.HandleFunc("GET /api/routers/{id}/tools/updates", handlers.GetUpdates)
	mux.HandleFunc("POST /api/routers/{id}/tools/reboot", handlers.PostReboot)
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
