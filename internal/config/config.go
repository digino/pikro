package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ProfileMeta stores Pikro-only metadata for a hotspot user profile.
// These fields are not sent to RouterOS — they are used for ticket display,
// pricing, and computing expiry comments written onto individual users.
type ProfileMeta struct {
	Validity string `json:"validity"` // duration shorthand, e.g. "1d", "4h", "1w"
	Price    string `json:"price"`
}

type LoginPageSettings struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Template string `json:"template,omitempty"` // "minimal" | "wave" | "card"
}

type HotspotSettings struct {
	HotspotName  string                 `json:"hotspotName"`
	DNSName      string                 `json:"dnsName"`
	Currency     string                 `json:"currency"`
	ProfileMetas map[string]ProfileMeta `json:"profileMetas,omitempty"`
	LoginPage    *LoginPageSettings     `json:"loginPage,omitempty"`
}

type RouterProfile struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Host            string          `json:"host"`
	Port            int             `json:"port"`
	Username        string          `json:"username"`
	Password        string          `json:"password"`
	UseTLS          bool            `json:"useTls"`
	HotspotSettings *HotspotSettings `json:"hotspotSettings,omitempty"`
}

type Config struct {
	Routers []RouterProfile `json:"routers"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pikro.json")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath())
	if os.IsNotExist(err) {
		return &Config{Routers: []RouterProfile{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0600)
}
