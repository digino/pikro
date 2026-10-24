---
paths:
  - 'web/src/**/*.vue'
  - 'web/src/**/*.ts'
---

# Vue / TypeScript conventions

- API calls in `web/src/api/index.ts` — typed interfaces, axios at `/api`.
- State in Pinia stores under `web/src/stores/`.
- UI primitives from Reka UI + Tailwind CSS v4 utilities.
- Icons: two sets with distinct jobs.
  - Ink icons (`web/src/assets/icons/ink/`, via `inkIcon()` in `utils/inkIcons.ts`) for
    navigation, page headers and empty states. They are multicolour, so render them as
    `<img alt="">`, never inline (each SVG carries its own filter id). New ones: the
    `icon-forge` skill, preset `ink`, accent `#FFBF00`.
  - `@heroicons/vue/24/outline` for small UI glyphs (chevron, search, close) — always
    the outline variant.
- Ink & Amber styling lives in `style.css`: use the `card`, `chip`, `btn`, `eyebrow` and
  `type-display` utilities and the palette tokens (`ink`, `paper`, `ground`, `hairline`,
  `graphite`, `primary`, `primary-wash`, `online`, `danger`, ...). `btn-primary` (amber)
  is for the ONE main action per screen; use `btn-ink` or plain `btn` for the rest.
  Every token has a `[data-theme="dark"]` value, so never hard-code hex colours.
- Raw RouterOS responses typed as `Record<string, string>` — no other `any`.
- Use `<script setup lang="ts">`; leverage defineProps, defineEmits, defineModel and
  other macros for type-safe component APIs.
- Prefer Composition API over Options API; avoid reactive props destructuring to
  maintain reactivity.
- Core reactivity: ref, shallowRef, computed, watch, watchEffect, lifecycle hooks —
  prefer `shallowRef` when deep reactivity isn't needed.
- Built-in components (Transition, Teleport, Suspense, KeepAlive) and directives like
  `v-memo` are available for advanced patterns.

## Tailwind (v4) gotchas
- In `<style scoped>`, `@apply` requires `@reference "tailwindcss"` first.
- `@apply` CANNOT use custom `@theme` tokens (e.g. `bg-surface`, `text-text-primary`)
  — those only work as classes in the template. In `<style scoped>`, use plain CSS
  `var(--color-*)` for custom token references; `@apply` is safe for structural
  utilities only (`flex`, `items-center`, `gap-2`).
- Avoid arbitrary values like `text-[13px]` or `rounded-[7px]` — prefer a named scale
  value (`text-sm`, `rounded-lg`). Only use `[value]` when no standard token is close
  enough (e.g. `w-[220px]` for a fixed sidebar width, or sub-pixel sizes like
  `text-[10.5px]`).
