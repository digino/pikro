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
