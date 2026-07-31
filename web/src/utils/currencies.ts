export interface CurrencyOption {
  value: string
  label: string
}

export const CURRENCIES: CurrencyOption[] = [
  { value: 'USD', label: 'USD — US Dollar' },
  { value: 'EUR', label: 'EUR — Euro' },
  { value: 'GBP', label: 'GBP — British Pound' },
  { value: 'XOF', label: 'XOF — CFA Franc (West Africa)' },
  { value: 'XAF', label: 'XAF — CFA Franc (Central Africa)' },
  { value: 'NGN', label: 'NGN — Nigerian Naira' },
  { value: 'KES', label: 'KES — Kenyan Shilling' },
  { value: 'GHS', label: 'GHS — Ghanaian Cedi' },
  { value: 'ZAR', label: 'ZAR — South African Rand' },
  { value: 'TZS', label: 'TZS — Tanzanian Shilling' },
  { value: 'UGX', label: 'UGX — Ugandan Shilling' },
  { value: 'ETB', label: 'ETB — Ethiopian Birr' },
  { value: 'RWF', label: 'RWF — Rwandan Franc' },
  { value: 'MZN', label: 'MZN — Mozambican Metical' },
]

export const KNOWN_CURRENCY_VALUES = ['', ...CURRENCIES.map(c => c.value)]

const compactFormatter = new Intl.NumberFormat('en', { notation: 'compact', maximumFractionDigits: 1 })

// Formats a revenue amount with k/M/B suffixes (e.g. 1234 -> "1.2K"), optionally
// suffixed with a currency code (e.g. "1.2K XOF").
export function formatCompactAmount(n: number, currency?: string): string {
  if (!n) return currency ? `0 ${currency}` : '0'
  const s = compactFormatter.format(n)
  return currency ? `${s} ${currency}` : s
}
