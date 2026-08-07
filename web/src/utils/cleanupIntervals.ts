import type { SelectOption } from '@/components/AppSelect.vue'

export const CLEANUP_INTERVAL_OPTIONS: SelectOption[] = [
  { value: '10m', label: 'Every 10 minutes (testing)' },
  { value: '1h', label: 'Every hour' },
  { value: '1d', label: 'Daily' },
  { value: '7d', label: 'Weekly (recommended)' },
  { value: '30d', label: 'Monthly' },
]
