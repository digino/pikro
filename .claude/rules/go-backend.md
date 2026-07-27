---
paths:
  - 'internal/**/*.go'
  - 'main.go'
---

# Go backend conventions

- Handlers in `internal/handlers/` — thin: parse request, call client, write response.
- RouterOS API calls in `internal/router/client.go` — no HTTP knowledge there.
- Config via `internal/config/config.go` — never read `~/.pikro.json` directly.
- Use `jsonOK` / `jsonError` / `jsonCreated` from `internal/handlers/helpers.go`.
- Build tags: `scan_unix.go` for Linux/macOS/BSD, `scan_windows.go` for the rest.
- Do not use `SO_REUSEPORT` — breaks Linux cross-compile.
- Do not use `html/template` for hotspot pages — use `text/template` to avoid
  HTML-escaping RouterOS tokens like `$(username)`.
- Version injected via ldflags: `main.Version` -> `handlers.AppVersion`.
- Never send router passwords to the browser — return profiles with the password
  field stripped (see the `safeProfile` pattern in `internal/handlers/routers.go`).

## Tests
- Prefer integration-style over unit tests that mock everything.
- Use Go's table-driven pattern: `[]struct{ name, input, want }{}`.
- Prefer parameterized (table-driven) tests over many near-duplicate cases.
