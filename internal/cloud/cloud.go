// Package cloud is the client for Pikro Cloud, used by the opt-in Cloud
// setup a certified technician runs on an owner's router.
package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// URL is where Pikro Cloud runs. PIKRO_CLOUD_URL overrides it, e.g.
// http://localhost:8081 for a local broker.
func URL() string {
	if u := os.Getenv("PIKRO_CLOUD_URL"); u != "" {
		return strings.TrimRight(u, "/")
	}
	return "https://cloud.getpikro.com"
}

type Client struct {
	BaseURL string
	HTTP    http.Client
}

func New() *Client {
	return &Client{BaseURL: URL(), HTTP: http.Client{Timeout: 30 * time.Second}}
}

// Error is an answer from Cloud other than success; Message is Cloud's own.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string { return fmt.Sprintf("Pikro Cloud: %s (%d)", e.Message, e.Status) }

type Technician struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Country   string `json:"country"`
}

type Person struct {
	Email      string      `json:"email"`
	Admin      bool        `json:"admin"`
	Technician *Technician `json:"technician"` // nil unless an active certified technician
}

type Session struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Person    Person    `json:"person"`
}

type Plan struct {
	Label    string `json:"label"`
	Price    int64  `json:"price"`
	Currency string `json:"currency"`
	Profile  string `json:"profile"`
}

type Registration struct {
	HardwareID      string `json:"hardware_id"`
	RouterOSVersion string `json:"routeros_version"`
	OwnerPhone      string `json:"owner_phone"`
	Name            string `json:"name"`         // the hotspot's name, shown on the buy page
	RouterToken     string `json:"router_token"` // from the router's sync script; "" on first setup
	Plans           []Plan `json:"plans"`
}

type Registered struct {
	RouterID string `json:"router_id"`
	Script   string `json:"script"`
}

// call sends body (if any) as JSON with a technician's session token (if
// any), and decodes a 2xx answer into out (if any).
func (c *Client) call(ctx context.Context, method, path, token string, body, out any) error {
	header := map[string]string{}
	if token != "" {
		header["Authorization"] = "Bearer " + token
	}
	return c.send(ctx, method, path, header, body, out)
}

func (c *Client) send(ctx context.Context, method, path string, header map[string]string, body, out any) error {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		var e struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&e)
		if e.Error == "" {
			e.Error = http.StatusText(resp.StatusCode)
		}
		return &Error{Status: resp.StatusCode, Message: e.Error}
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// SendCode asks Cloud to email a sign-in code.
func (c *Client) SendCode(ctx context.Context, email string) error {
	return c.call(ctx, http.MethodPost, "/auth/code", "", map[string]string{"email": email}, nil)
}

func (c *Client) CreateSession(ctx context.Context, email, code string) (*Session, error) {
	var s Session
	err := c.call(ctx, http.MethodPost, "/auth/session", "", map[string]string{"email": email, "code": code}, &s)
	return &s, err
}

// SignOut ends the session on Cloud.
func (c *Client) SignOut(ctx context.Context, token string) error {
	return c.call(ctx, http.MethodDelete, "/auth/session", token, nil, nil)
}

// Disconnect tells Cloud to stop selling for a router, authenticated by the
// router's own token ("<router-id>.<secret>", from its sync script).
func (c *Client) Disconnect(ctx context.Context, routerToken string) error {
	return c.send(ctx, http.MethodPost, "/routers/disconnect", map[string]string{"X-Router-Token": routerToken}, nil, nil)
}

// Me tells who a session belongs to now (a technician may have been revoked).
func (c *Client) Me(ctx context.Context, token string) (*Person, error) {
	var p Person
	err := c.call(ctx, http.MethodGet, "/auth/me", token, nil, &p)
	return &p, err
}

// RegisterRouter registers (or re-registers) a router and returns the sync
// script to install, already holding its URL and token.
func (c *Client) RegisterRouter(ctx context.Context, token string, r Registration) (*Registered, error) {
	var res Registered
	err := c.call(ctx, http.MethodPost, "/routers", token, r, &res)
	return &res, err
}
