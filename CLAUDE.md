# Pikro — CLAUDE.md

Project-specific conventions. Global conventions (approach, delegation, tests, comments,
formatting, time estimates, safety) are in `~/.claude/CLAUDE.md` and apply here too.

---

## Project overview

Pikro is a zero-install web app for managing MikroTik RouterOS devices.
A single Go binary serves a Vue 3 SPA and communicates with routers over the RouterOS
native API (port 8728/8729 — the same protocol WinBox uses).

Target users: small ISPs, café/hotel operators, community network admins in markets
where MikroTik is dominant (especially West and Central Africa).

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

No test suite exists yet. When adding tests:
- Prefer integration-style over unit tests that mock everything
- Use Go's table-driven pattern: `[]struct{ name, input, want }{}`
- RouterOS scripts cannot be unit tested locally — validate against
  `mikrotik-routeros-rsc/` references and test on real hardware

---

## RouterOS scripting

Always consult the skill and references in `mikrotik-routeros-rsc/` before writing
or modifying any script. Comment non-obvious sections — RSC is obscure to future readers.

Key gotchas:
- `[:tonsec [:timestamp]]` fails silently when assigned to `:local` in scheduler context
- String comparison operators differ between v6 and v7 — use the v6-compatible
  Mikhmon dateint/timeint pattern when targeting both
- `=on-login=` scripts must be idempotent — they fire on every login, not just the first
- Empty string sent to RouterOS is treated as a literal value — omit fields entirely when blank

---

## Go conventions

- Handlers in `internal/handlers/` — thin: parse request, call client, write response
- RouterOS API calls in `internal/router/client.go` — no HTTP knowledge there
- Config via `internal/config/config.go` — never read `~/.pikro.json` directly
- Use `jsonOK` / `jsonError` / `jsonCreated` from `internal/handlers/helpers.go`
- Build tags: `scan_unix.go` for Linux/macOS/BSD, `scan_windows.go` for the rest
- Do not use `SO_REUSEPORT` — breaks Linux cross-compile
- Do not use `html/template` for hotspot pages — use `text/template` to avoid
  HTML-escaping RouterOS tokens like `$(username)`
- Version injected via ldflags: `main.Version` -> `handlers.AppVersion`

---

## Vue / TypeScript conventions

- API calls in `web/src/api/index.ts` — typed interfaces, axios at `/api`
- State in Pinia stores under `web/src/stores/`
- UI primitives from Reka UI + Tailwind CSS v4 utilities
- Icons from `@heroicons/vue/24/outline` — always outline variant
- Tailwind in `<style scoped>` requires `@reference "tailwindcss"` before `@apply`
- Raw RouterOS responses typed as `Record<string, string>` — no other `any`
- **Avoid arbitrary Tailwind values** like `text-[13px]` or `rounded-[7px]` — always
  prefer a named scale value (`text-sm`, `rounded-lg`). Only use `[value]` syntax when
  no standard token is close enough (e.g. `w-[220px]` for a fixed sidebar width, or
  sub-pixel font sizes like `text-[10.5px]` that have no scale equivalent).
- **`@apply` cannot use custom `@theme` tokens** (e.g. `bg-surface`, `text-text-primary`)
  — those only work as classes in the template. In `<style scoped>`, use plain CSS
  `var(--color-*)` for custom token references; `@apply` is safe for structural
  Tailwind utilities only (`flex`, `items-center`, `gap-2`, etc.).

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
