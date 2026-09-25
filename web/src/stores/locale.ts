import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Locale } from '@/locales'

export const useLocaleStore = defineStore('locale', () => {
  const { locale: i18nLocale } = useI18n()

  const stored = localStorage.getItem('pikro-locale') as Locale | null
  const locale = ref<Locale>(stored ?? 'en')

  function apply(l: Locale) {
    i18nLocale.value = l
    localStorage.setItem('pikro-locale', l)
  }

  function set(l: Locale) {
    locale.value = l
  }

  watch(locale, apply, { immediate: true })

  return { locale, set }
})
