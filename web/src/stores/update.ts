import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getUpdateStatus, type UpdateStatus } from '@/api'

const DISMISSED_KEY = 'pikro-update-dismissed'

/** "vv0.4.0", "0.4.0" and "v0.4.0" all read "v0.4.0"; "dev" stays as is. */
export function displayVersion(v: string) {
  const bare = v.replace(/^v+/i, '')
  return /^\d/.test(bare) ? `v${bare}` : v
}

export const useUpdateStore = defineStore('update', () => {
  const status = ref<UpdateStatus | null>(null)
  const checking = ref(false)
  // The release the operator closed the banner for: it stays hidden until
  // a newer one comes out.
  const dismissed = ref(readDismissed())

  const showBanner = computed(() =>
    !!status.value?.available && status.value.latest !== dismissed.value,
  )

  async function check(refresh = false) {
    checking.value = true
    try {
      status.value = await getUpdateStatus(refresh)
    } catch {
      // Pikro's own server didn't answer: nothing to say about updates.
    } finally {
      checking.value = false
    }
  }

  function dismiss() {
    if (!status.value) return
    dismissed.value = status.value.latest
    try {
      localStorage.setItem(DISMISSED_KEY, status.value.latest)
    } catch { /* private mode: dismissed for this session only */ }
  }

  return { status, checking, showBanner, check, dismiss }
})

function readDismissed() {
  try {
    return localStorage.getItem(DISMISSED_KEY) ?? ''
  } catch {
    return ''
  }
}
