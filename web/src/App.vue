<template>
  <div class="flex h-screen overflow-hidden bg-base">
    <AppSidebar :app-version="appVersion" />
    <div class="flex-1 overflow-y-auto min-w-0">
      <RouterView />
    </div>
    <AddRouterDialog v-model:open="store.showAddDialog" />
    <ToastHost />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getAppVersion } from '@/api'
import { useRoutersStore } from '@/stores/routers'
import { useThemeStore } from '@/stores/theme'
import AppSidebar from '@/components/AppSidebar.vue'
import AddRouterDialog from '@/components/AddRouterDialog.vue'
import ToastHost from '@/components/ToastHost.vue'

const store = useRoutersStore()
const appVersion = ref('dev')

// Init theme on boot — applies stored preference to <html data-theme>.
useThemeStore()

onMounted(() => {
  store.load()
  getAppVersion().then(v => { appVersion.value = v }).catch(() => {})
})
</script>
