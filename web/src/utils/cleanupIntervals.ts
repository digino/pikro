import type { SelectOption } from '@/components/AppSelect.vue'

export const CLEANUP_INTERVAL_OPTIONS: SelectOption[] = [
  { value: '10m', label: 'Every 10 minutes' },
  { value: '30m', label: 'Every 30 minutes' },
  { value: '1h', label: 'Every hour (recommended)' },
]
