package cloud

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeCloud answers like pikro-cloud's broker for one technician session.
func fakeCloud(t *testing.T) *Client {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/code", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /auth/session", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email, Code string }
		json.NewDecoder(r.Body).Decode(&req)
		if req.Code != "123456" {
			http.Error(w, `{"error":"wrong or expired code"}`, http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`{"token":"TOK","expires_at":"2026-12-04T10:00:00Z","person":{"email":"` + req.Email +
			`","admin":false,"technician":{"first_name":"Kossi","last_name":"Mensah"}}}`))
	})
	mux.HandleFunc("DELETE /auth/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer TOK" {
			http.Error(w, `{"error":"not signed in"}`, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /routers", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer TOK" {
			http.Error(w, `{"error":"not signed in"}`, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"router_id":"R1","secret":"S","script":":local token \"R1.S\""}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return &Client{BaseURL: srv.URL}
}

func TestSignIn(t *testing.T) {
	ctx := context.Background()
	c := fakeCloud(t)
	if err := c.SendCode(ctx, "kossi@pikro.test"); err != nil {
		t.Fatal(err)
	}
	s, err := c.CreateSession(ctx, "kossi@pikro.test", "123456")
	if err != nil {
		t.Fatal(err)
	}
	if s.Token != "TOK" || s.Person.Technician == nil || s.Person.Technician.FirstName != "Kossi" || s.ExpiresAt.IsZero() {
		t.Errorf("session %+v", s)
	}

	if err := c.SignOut(ctx, s.Token); err != nil {
		t.Errorf("sign out: %v", err)
	}

	_, err = c.CreateSession(ctx, "kossi@pikro.test", "000000")
	var e *Error
	if !errors.As(err, &e) || e.Status != http.StatusUnauthorized || e.Message != "wrong or expired code" {
		t.Errorf("wrong code: err = %#v, want a 401 with Cloud's message", err)
	}
}

func TestRegisterRouter(t *testing.T) {
	ctx := context.Background()
	c := fakeCloud(t)
	reg := Registration{HardwareID: "HW1", RouterOSVersion: "6.49.23", OwnerPhone: "+22890000000",
		Plans: []Plan{{Label: "1 jour", Price: 300, Currency: "XOF", Profile: "1-jour"}}}

	res, err := c.RegisterRouter(ctx, "TOK", reg)
	if err != nil || res.RouterID != "R1" || res.Script == "" {
		t.Errorf("register: %+v, err %v", res, err)
	}
	var e *Error
	if _, err := c.RegisterRouter(ctx, "expired", reg); !errors.As(err, &e) || e.Status != http.StatusUnauthorized {
		t.Errorf("expired session: err = %v, want 401", err)
	}
}
