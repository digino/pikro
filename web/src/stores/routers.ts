import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as api from '@/api'

const ACTIVE_ID_KEY = 'pikro-active-router-id'

export const useRoutersStore = defineStore('routers', () => {
  const routers = ref<api.RouterProfile[]>([])
  const activeId = ref<string | null>(localStorage.getItem(ACTIVE_ID_KEY))
  const loading = ref(false)
  const showAddDialog = ref(false)

  const active = () => routers.value.find(r => r.id === activeId.value) ?? null

  function setActiveId(id: string | null) {
    activeId.value = id
    if (id) localStorage.setItem(ACTIVE_ID_KEY, id)
    else localStorage.removeItem(ACTIVE_ID_KEY)
  }

  async function load() {
    loading.value = true
    try {
      routers.value = await api.listRouters()
      // The saved router may have been deleted while the app was closed —
      // fall back to the first router only when there's truly no valid selection.
      if (activeId.value && !routers.value.some(r => r.id === activeId.value)) {
        setActiveId(null)
      }
      if (!activeId.value && routers.value.length > 0) {
        setActiveId(routers.value[0].id)
      }
    } finally {
      loading.value = false
    }
  }

  async function add(profile: api.NewRouterProfile) {
    const { id } = await api.addRouter(profile)
    await load()
    setActiveId(id)
  }

  async function remove(id: string) {
    await api.deleteRouter(id)
    if (activeId.value === id) setActiveId(null)
    await load()
  }

  function select(id: string) {
    setActiveId(id)
  }

  return { routers, activeId, loading, showAddDialog, active, load, add, remove, select }
})
