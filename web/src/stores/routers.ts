import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as api from '@/api'

export const useRoutersStore = defineStore('routers', () => {
  const routers = ref<api.RouterProfile[]>([])
  const activeId = ref<string | null>(null)
  const loading = ref(false)
  const showAddDialog = ref(false)

  const active = () => routers.value.find(r => r.id === activeId.value) ?? null

  async function load() {
    loading.value = true
    try {
      routers.value = await api.listRouters()
      if (!activeId.value && routers.value.length > 0) {
        activeId.value = routers.value[0].id
      }
    } finally {
      loading.value = false
    }
  }

  async function add(profile: api.NewRouterProfile) {
    const { id } = await api.addRouter(profile)
    await load()
    activeId.value = id
  }

  async function remove(id: string) {
    await api.deleteRouter(id)
    if (activeId.value === id) activeId.value = null
    await load()
  }

  function select(id: string) {
    activeId.value = id
  }

  return { routers, activeId, loading, showAddDialog, active, load, add, remove, select }
})
