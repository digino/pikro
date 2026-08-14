import { createI18n } from 'vue-i18n'
import en from './en.json'
import fr from './fr.json'

export type Locale = 'en' | 'fr'

const stored = localStorage.getItem('pikro-locale') as Locale | null

const i18n = createI18n({
  legacy: false,
  locale: stored ?? 'en',
  fallbackLocale: 'en',
  messages: { en, fr },
})

export default i18n
