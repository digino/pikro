---
paths:
  - 'web/src/**/*.vue'
  - 'web/src/**/*.ts'
---

# Vue / TypeScript conventions

- API calls in `web/src/api/index.ts` — typed interfaces, axios at `/api`.
- State in Pinia stores under `web/src/stores/`.
- UI primitives from Reka UI + Tailwind CSS v4 utilities.
- Icons from `@heroicons/vue/24/outline` — always the outline variant.
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
