// Package update tells whether a newer Pikro release is out. It only checks:
// Pikro never downloads or replaces itself (see the roadmap's auto-update
// decision).
package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	latestReleaseURL = "https://api.github.com/repos/digino/pikro/releases/latest"
	// GitHub allows 60 unauthenticated requests an hour per IP; once a day is plenty.
	cacheFor = 24 * time.Hour
	// A failed check is retried sooner, but not on every page load.
	retryAfter = time.Hour
)

type Info struct {
	Current   string `json:"current"`
	Latest    string `json:"latest"`
	Available bool   `json:"available"`
}

type Checker struct {
	current string
	client  *http.Client

	mu        sync.Mutex
	latest    string
	checkedAt time.Time
	failedAt  time.Time
}

func NewChecker(current string) *Checker {
	return &Checker{current: current, client: &http.Client{Timeout: 8 * time.Second}}
}

// Check returns the cached answer, asking GitHub when it is stale or when
// force is set (the operator pressed "Check again").
func (c *Checker) Check(ctx context.Context, force bool) (Info, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	stale := now.Sub(c.checkedAt) > cacheFor && now.Sub(c.failedAt) > retryAfter
	if force || stale {
		latest, err := c.fetchLatest(ctx)
		if err != nil {
			c.failedAt = now
			return Info{Current: c.current}, err
		}
		c.latest, c.checkedAt = latest, now
	}
	return Info{
		Current:   c.current,
		Latest:    c.latest,
		Available: c.latest != "" && isNewer(c.latest, c.current),
	}, nil
}

func (c *Checker) fetchLatest(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestReleaseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "pikro/"+c.current)
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub answered %s", resp.Status)
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	return release.TagName, nil
}

// isNewer compares release versions like "v1.2.3" or "1.2.3-rc.1". Any
// number of leading v's is ignored (one release was tagged "vv0.4.0"); a
// release is newer than its own pre-releases. Development builds ("dev")
// and unreadable versions never count as outdated.
func isNewer(latest, current string) bool {
	l, okL := parse(latest)
	c, okC := parse(current)
	if !okL || !okC {
		return false
	}
	for i := range 3 {
		if l.nums[i] != c.nums[i] {
			return l.nums[i] > c.nums[i]
		}
	}
	// Same numbers: only a release following its pre-release is newer.
	return l.pre == "" && c.pre != ""
}

type version struct {
	nums [3]int
	pre  string
}

func parse(s string) (version, bool) {
	s = strings.TrimLeft(s, "vV")
	core, pre, _ := strings.Cut(s, "-")
	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return version{}, false
	}
	var v version
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return version{}, false
		}
		v.nums[i] = n
	}
	v.pre = pre
	return v, true
}
