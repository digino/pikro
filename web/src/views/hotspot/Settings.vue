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

      <!-- ── Tab: General ── -->
      <div v-if="tab === 'general'" class="max-w-lg">
        <form class="space-y-5" @submit.prevent="save">
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary">Business name</span>
            <input v-model="form.voucherBusinessName" class="input" placeholder="e.g. CyberCafé Lumière" />
            <p class="text-xs text-text-muted">Shown on vouchers and as the login page title.</p>
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary">Hotspot name</span>
            <input v-model="form.hotspotName" class="input" placeholder="e.g. myspot" />
            <p class="text-xs text-text-muted">The name of the hotspot server configured on your router.</p>
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary">DNS name</span>
            <input v-model="form.dnsName" class="input" placeholder="e.g. myspot.spot" />
            <p class="text-xs text-text-muted">The DNS name clients are redirected to for the login page.</p>
          </label>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary">Currency</span>
            <select v-model="form.currency" class="input">
              <option value="">None</option>
              <option v-for="c in CURRENCIES" :key="c.value" :value="c.value">{{ c.label }}</option>
              <option value="custom">Custom…</option>
            </select>
            <input
              v-if="form.currency === 'custom'"
              v-model="customCurrency"
              class="input mt-1"
              placeholder="Enter currency code, e.g. MGA"
              maxlength="6"
            />
            <p class="text-xs text-text-muted">Used when displaying prices on profiles and vouchers.</p>
          </label>

          <p v-if="settingsError" class="text-xs text-red">{{ settingsError }}</p>
          <p v-if="settingsSaved" class="text-xs text-green">Settings saved.</p>

          <div class="flex justify-end pt-1">
            <button type="submit" class="btn btn-primary" :disabled="saving">
              <span v-if="saving" class="size-4 border-2 border-black/20 border-t-black rounded-full animate-spin" />
              Save settings
            </button>
          </div>
        </form>
      </div>

      <!-- ── Tab: Login page ── -->
      <SettingsLoginPage
        v-else-if="tab === 'login'"
        ref="loginPageRef"
        :business-name="form.voucherBusinessName"
      />

      <!-- ── Tab: Vouchers ── -->
      <SettingsVoucher
        v-else-if="tab === 'voucher'"
        ref="voucherRef"
        :business-name="form.voucherBusinessName"
        :effective-currency="effectiveCurrency"
      />

      <!-- ── Tab: Cleanup ── -->
      <div v-else-if="tab === 'cleanup'" class="max-w-lg space-y-5">
        <div>
          <h3 class="text-sm font-semibold text-text-primary">Auto-cleanup</h3>
          <p class="text-sm text-text-muted mt-0.5">
            Expired hotspot users are automatically removed on a schedule.
            The scheduler is installed on the router during hotspot setup.
          </p>
        </div>

        <div class="flex items-center gap-2 p-3 border border-border rounded-lg text-xs text-text-secondary bg-surface">
          <CheckCircleIcon class="size-4 shrink-0 text-green" />
          Scheduler <span class="font-mono text-text-primary mx-1">pikro-cleanup</span>
          {{ cleanup.installed ? 'active on router' : 'not yet installed — run hotspot setup' }}
        </div>

        <div class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Run every</span>
          <select v-model="cleanup.interval" class="input" @change="saveCleanup">
            <option value="10m">Every 10 minutes (testing)</option>
            <option value="1h">Every hour</option>
            <option value="1d">Daily</option>
            <option value="7d">Weekly (recommended)</option>
          </select>
        </div>

        <div v-if="deviceModeBlocked" class="p-4 border rounded-xl space-y-2 bg-amber/8 border-amber/20">
          <div class="flex items-start gap-2">
            <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5 text-amber" />
            <div class="space-y-1">
              <p class="text-xs font-semibold text-amber">Scheduler blocked by RouterOS device-mode</p>
              <p class="text-xs text-text-secondary">
                Your router runs RouterOS 7.17+ with the scheduler disabled. Run the command below in Winbox Terminal or SSH,
                then press the physical button on your router (or cold-reboot) to confirm.
              </p>
            </div>
          </div>
          <div class="border border-border rounded-lg px-3 py-2 font-mono text-xs text-text-primary select-all bg-base">
            /system/device-mode/update scheduler=yes
          </div>
          <p class="text-xs text-text-muted">After confirming, re-run hotspot setup to install the scheduler.</p>
        </div>
        <p v-else-if="cleanupError" class="text-xs text-red">{{ cleanupError }}</p>
      </div>

      <!-- ── Tab: Migration ── -->
      <SettingsMigration v-else-if="tab === 'migration'" />
    </template>
  </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { CheckCircleIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import {
  getHotspotSettings, putHotspotSettings,
  getCleanupScheduler, putCleanupScheduler,
} from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import { CURRENCIES, KNOWN_CURRENCY_VALUES } from '@/utils/currencies'
import SettingsLoginPage from './settings/SettingsLoginPage.vue'
import SettingsVoucher from './settings/SettingsVoucher.vue'
import SettingsMigration from './settings/SettingsMigration.vue'

const store = useRoutersStore()

const tab = ref<'general' | 'login' | 'voucher' | 'cleanup' | 'migration'>('general')
const tabs = [
  { key: 'general' as const,  label: 'General' },
  { key: 'login' as const,    label: 'Login page' },
  { key: 'voucher' as const,  label: 'Vouchers' },
  { key: 'cleanup' as const,  label: 'Cleanup' },
  { key: 'migration' as const, label: 'Migration' },
]

const loading = ref(false)
const saving = ref(false)
const settingsSaved = ref(false)
const settingsError = ref('')
const customCurrency = ref('')

const form = ref({ hotspotName: '', dnsName: '', currency: '', voucherBusinessName: '' })

const effectiveCurrency = computed(() =>
  form.value.currency === 'custom' ? customCurrency.value : form.value.currency
)

const cleanupSaving = ref(false)
const cleanupError = ref('')
const cleanup = ref({ installed: false, interval: '7d' })

const deviceModeBlocked = computed(() =>
  cleanupError.value.toLowerCase().includes('device-mode')
)

const loginPageRef = ref<InstanceType<typeof SettingsLoginPage> | null>(null)
const voucherRef = ref<InstanceType<typeof SettingsVoucher> | null>(null)

async function load() {
  if (!store.activeId) return
  loading.value = true
  settingsError.value = ''
  cleanupError.value = ''
  try {
    const [s, c] = await Promise.all([
      getHotspotSettings(store.activeId),
      getCleanupScheduler(store.activeId).catch(() => ({ installed: false, interval: '' })),
    ])
    form.value.hotspotName = s.hotspotName ?? ''
    form.value.dnsName = s.dnsName ?? ''
    form.value.currency = s.currency ?? ''
    form.value.voucherBusinessName = s.voucher?.businessName ?? ''
    if (form.value.currency && !KNOWN_CURRENCY_VALUES.includes(form.value.currency)) {
      customCurrency.value = form.value.currency
      form.value.currency = 'custom'
    }
    loginPageRef.value?.init(s.loginPage)
    voucherRef.value?.init(s.voucher)
    cleanup.value = { installed: c.installed, interval: c.interval || '7d' }
  } catch (e: any) {
    settingsError.value = e?.response?.data?.error ?? e?.message ?? 'Failed to load settings'
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!store.activeId) return
  saving.value = true; settingsSaved.value = false; settingsError.value = ''
  try {
    const existing = await getHotspotSettings(store.activeId)
    await putHotspotSettings(store.activeId, {
      ...existing,
      hotspotName: form.value.hotspotName,
      dnsName: form.value.dnsName,
      currency: effectiveCurrency.value,
      voucher: {
        showValidity: existing.voucher?.showValidity ?? true,
        showPrice: existing.voucher?.showPrice ?? true,
        layout: existing.voucher?.layout,
        businessName: form.value.voucherBusinessName,
      },
    })
    settingsSaved.value = true
    setTimeout(() => { settingsSaved.value = false }, 3000)
  } catch (e: any) {
    settingsError.value = e?.response?.data?.error ?? e?.message ?? 'Failed to save settings'
  } finally { saving.value = false }
}

async function saveCleanup() {
  if (!store.activeId) return
  cleanupSaving.value = true; cleanupError.value = ''
  try {
    const result = await putCleanupScheduler(store.activeId, true, cleanup.value.interval)
    cleanup.value.installed = result.installed
    if (result.interval) cleanup.value.interval = result.interval
  } catch (e: any) {
    cleanupError.value = e?.response?.data?.error ?? e?.message ?? 'Failed to update scheduler'
  } finally { cleanupSaving.value = false }
}

watch(() => store.activeId, load, { immediate: true })
</script>
