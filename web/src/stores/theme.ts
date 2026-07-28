import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

type Theme = 'dark' | 'light'

export const useThemeStore = defineStore('theme', () => {
  const stored = localStorage.getItem('pikro-theme') as Theme | null
  const theme = ref<Theme>(stored ?? 'dark')

  function apply(t: Theme) {
    document.documentElement.setAttribute('data-theme', t)
    localStorage.setItem('pikro-theme', t)
  }

  function toggle() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
  }

  watch(theme, apply, { immediate: true })

  return { theme, toggle }
})
