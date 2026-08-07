import type { SelectOption } from '@/components/AppSelect.vue'

// Fake/local-only TLDs resolved by the router's own DNS for the hotspot
// login domain — must match the backend's validExtensions allowlist in
// internal/router/hotspot_setup.go.
export const HOTSPOT_EXTENSION_OPTIONS: SelectOption[] = [
  { value: '.spot', label: '.spot' },
  { value: '.hotspot', label: '.hotspot' },
  { value: '.info', label: '.info' },
  { value: '.wifi', label: '.wifi' },
]
