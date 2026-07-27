---
paths:
  - '**/*.rsc'
  - 'internal/router/**'
  - 'hotspot/**'
---

# RouterOS scripting (RSC)

Always consult the skill and references in `mikrotik-routeros-rsc/` before writing
or modifying any script. Comment non-obvious sections — RSC is obscure to future readers.

RSC scripts cannot be unit tested locally — validate against `mikrotik-routeros-rsc/`
references and test on real hardware.

Key gotchas:
- `[:tonsec [:timestamp]]` fails silently when assigned to `:local` in scheduler context
- String comparison operators differ between v6 and v7 — use the v6-compatible
  Mikhmon dateint/timeint pattern when targeting both
- `=on-login=` scripts must be idempotent — they fire on every login, not just the first
- Empty string sent to RouterOS is treated as a literal value — omit fields entirely when blank

When generating hotspot pages, do NOT use `html/template` — use `text/template` to avoid
HTML-escaping RouterOS tokens like `$(username)`.
