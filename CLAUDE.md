# Pikro — CLAUDE.md

Project-specific conventions. Global conventions (approach, delegation, tests, comments,
formatting, time estimates, safety) are in `~/.claude/CLAUDE.md` and apply here too.

---

## Project overview

Pikro is a zero-install web app for managing MikroTik RouterOS devices.
A single Go binary serves a Vue 3 SPA and communicates with routers over the RouterOS
native API (port 8728/8729 — the same protocol WinBox uses).

Target users: small ISPs, café/hotel operators, community network admins in markets
where MikroTik is dominant.

Key feature: hotspot voucher management — generate, print, expire, and clean up
hotspot users without touching the router CLI.

---

## Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22+, `net/http` mux, `//go:embed` for SPA |
| Router API | `github.com/go-routeros/routeros/v3` — native API port 8728 |
| Frontend | Vue 3, TypeScript, Vite, Tailwind CSS v4, Reka UI, Pinia |
| Config | `~/.pikro.json` — plain JSON, no database |
| Discovery | MNDP (MikroTik Neighbour Discovery Protocol) over UDP :5678 |

Full architecture: see [docs/architecture.md](docs/architecture.md).

---

## Complexity signals specific to this project

When estimating scope, flag these as high-risk multipliers:
- Change crosses the Go/Vue boundary (API contract change affects both sides)
- RouterOS scripting involved — hard to test locally, must run on real hardware
- Cross-compilation or release pipeline touched

---

## Tests

No test suite exists yet. TDD preferred. There is a `test-writer` agent in
`.claude/agents/` set up for Go table-driven tests — use it for non-trivial test work.

---

## Path-scoped conventions

Detailed, situational conventions live in `.claude/rules/` and load automatically when
you touch matching files — keeping this file lean. See:

- `.claude/rules/go-backend.md` — Go handler/client/config conventions + test patterns
  (paths: `internal/**/*.go`, `main.go`)
- `.claude/rules/vue-frontend.md` — Vue 3, Pinia, Reka UI, Tailwind v4 gotchas
  (paths: `web/src/**/*.vue`, `web/src/**/*.ts`)
- `.claude/rules/routeros-scripts.md` — RSC gotchas + `text/template` rule
  (paths: `**/*.rsc`, `internal/router/**`, `hotspot/**`)

---

## Makefile reference

```
make dev        # Go backend (:8080) + Vite dev server (:5173) together
make backend    # Go backend only
make build      # Vue build -> Go binary with version from git tag
make release    # Cross-compile 5 binaries into ./dist/
make clean      # Remove binary, dist/, web/dist/
```

Release binaries: `pikro-mac-arm64`, `pikro-mac-intel`, `pikro.exe`,
`pikro-linux-amd64`, `pikro-linux-arm64`.
