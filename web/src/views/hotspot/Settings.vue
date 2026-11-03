<template>
  <PageLayout :title="t('settings.crumb')" :subtitle="t('settings.title')" icon="settings">
    <template #actions>
      <div v-if="store.activeId" role="tablist" :aria-label="t('settings.sections')" class="flex flex-wrap gap-2">
        <button
          v-for="key in TABS" :key="key"
          type="button"
          role="tab"
          :aria-selected="tab === key"
          class="inline-flex items-center gap-2 min-h-11 px-4 rounded-full border-2 border-ink text-sm font-semibold cursor-pointer transition-colors"
          :class="tab === key ? 'bg-ink text-paper shadow-[3px_3px_0_var(--color-primary)]' : 'bg-paper text-ink hover:bg-primary-wash'"
          @click="tab = key"
        >
          {{ t(`settings.tabs.${key}`) }}
          <span
            v-if="key === 'migration'"
            class="rounded-full border-[1.5px] border-current px-1.75 py-0.5 text-xs"
          >{{ t('settings.beta') }}</span>
        </button>
      </div>
    </template>

    <NoRouterSelected v-if="!store.activeId" />

    <div v-else-if="loading" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <div
      v-else-if="settingsError"
      role="alert"
      class="flex items-center gap-2.5 p-4 rounded-control border-2 border-danger-text bg-danger-wash text-danger-text text-sm"
    >
      <ExclamationTriangleIcon class="size-5 shrink-0" />
      <span>{{ settingsError }}</span>
      <button type="button" class="btn btn-sm ml-auto" @click="load">{{ t('settings.retry') }}</button>
    </div>

    <template v-else-if="settings">
      <SettingsLoginPage v-if="tab === 'login'" :settings="settings" @saved="settings = $event" />
      <SettingsMigration v-else />
    </template>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import { getHotspotSettings, type HotspotSettings } from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import SettingsLoginPage from './settings/SettingsLoginPage.vue'
import SettingsMigration from './settings/SettingsMigration.vue'

const { t } = useI18n()
const store = useRoutersStore()

const TABS = ['login', 'migration'] as const
const tab = ref<(typeof TABS)[number]>('login')

const loading = ref(false)
const settingsError = ref('')

// hotspotName/dnsName/currency are edited from the router's Add/Edit dialog
// (see Routers page) — loaded here read-only for the login-page-title tab to consume.
const settings = ref<HotspotSettings | null>(null)

async function load() {
  if (!store.activeId) return
  loading.value = true
  settingsError.value = ''
  try {
    settings.value = await getHotspotSettings(store.activeId)
  } catch (e: any) {
    settingsError.value = e?.response?.data?.error ?? e?.message ?? t('settings.failedToLoad')
  } finally {
    loading.value = false
  }
}

watch(() => store.activeId, load, { immediate: true })
</script>
