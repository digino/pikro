import { defineStore } from 'pinia'
import { ref } from 'vue'

export type ToastVariant = 'success' | 'error' | 'info'

export interface Toast {
  id: number
  title: string
  description?: string
  variant: ToastVariant
}

let nextId = 1

export const useToastStore = defineStore('toast', () => {
  const toasts = ref<Toast[]>([])

  function push(title: string, opts: { description?: string; variant?: ToastVariant } = {}) {
    const id = nextId++
    toasts.value.push({ id, title, description: opts.description, variant: opts.variant ?? 'info' })
    return id
  }

  function dismiss(id: number) {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  function success(title: string, description?: string) {
    return push(title, { description, variant: 'success' })
  }

  function error(title: string, description?: string) {
    return push(title, { description, variant: 'error' })
  }

  return { toasts, push, dismiss, success, error }
})
