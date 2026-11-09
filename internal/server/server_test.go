package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocalOnly(t *testing.T) {
	tests := []struct {
		name, host, origin string
		want               int
	}{
		{"the app", "localhost:8080", "", http.StatusOK},
		{"the app, same origin", "localhost:8080", "http://localhost:8080", http.StatusOK},
		{"vite dev server", "localhost:5173", "http://localhost:5173", http.StatusOK},
		{"ipv4 loopback", "127.0.0.1:8080", "", http.StatusOK},
		{"ipv6 loopback", "[::1]:8080", "", http.StatusOK},
		{"dns rebinding", "evil.example:8080", "", http.StatusForbidden},
		{"lan address", "192.168.1.20:8080", "", http.StatusForbidden},
		{"another website's page", "localhost:8080", "https://evil.example", http.StatusForbidden},
		{"a local file", "localhost:8080", "null", http.StatusForbidden},
	}
	ok := localOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/api/routers", nil)
			r.Host = tt.host
			if tt.origin != "" {
				r.Header.Set("Origin", tt.origin)
			}
			w := httptest.NewRecorder()
			ok.ServeHTTP(w, r)
			if w.Code != tt.want {
				t.Errorf("status %d, want %d", w.Code, tt.want)
			}
		})
	}
}
