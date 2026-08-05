<template>
  <PageLayout title="Hotspot" subtitle="Settings">
  <div>
    <NoRouterSelected v-if="!store.activeId" />

    <div v-else-if="loading" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <template v-else>
      <!-- Tabs -->
      <div class="flex items-center gap-1 border-b border-border mb-6 -mt-1">
        <button
          v-for="t in tabs" :key="t.key"
          class="px-3 py-2 text-sm font-semibold border-b-2 transition-colors"
          :class="tab === t.key
            ? 'border-text-primary text-text-primary'
            : 'border-transparent text-text-muted hover:text-text-secondary'"
          @click="tab = t.key"
        >
          {{ t.label }}
        </button>
      </div>

      <!-- ── Tab: Login page ── -->
      <SettingsLoginPage
        v-if="tab === 'login'"
        ref="loginPageRef"
        :hotspot-name="form.hotspotName"
      />

      <!-- ── Tab: Migration ── -->
      <SettingsMigration v-else-if="tab === 'migration'" />
    </template>
  </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoutersStore } from '@/stores/routers'
import { getHotspotSettings } from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import SettingsLoginPage from './settings/SettingsLoginPage.vue'
import SettingsMigration from './settings/SettingsMigration.vue'

const store = useRoutersStore()

const tab = ref<'login' | 'migration'>('login')
const tabs = [
  { key: 'login' as const,    label: 'Login page' },
  { key: 'migration' as const, label: 'Migration' },
]

const loading = ref(false)
const settingsError = ref('')

// hotspotName/dnsName/currency are edited from the router's Add/Edit dialog
// (see Routers page) — loaded here read-only for the login-page-title tab to consume.
const form = ref({ hotspotName: '', dnsName: '', currency: '' })

const loginPageRef = ref<InstanceType<typeof SettingsLoginPage> | null>(null)

async function load() {
  if (!store.activeId) return
  loading.value = true
  settingsError.value = ''
  try {
    const s = await getHotspotSettings(store.activeId)
    form.value.hotspotName = s.hotspotName ?? ''
    form.value.dnsName = s.dnsName ?? ''
    form.value.currency = s.currency ?? ''
    loginPageRef.value?.init(s.loginPage)
  } catch (e: any) {
    settingsError.value = e?.response?.data?.error ?? e?.message ?? 'Failed to load settings'
  } finally {
    loading.value = false
  }
}

watch(() => store.activeId, load, { immediate: true })
</script>
