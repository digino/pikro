<template>
  <div>
    <div v-if="!store.activeId" class="border border-dashed border-border rounded-xl py-12 text-center">
      <p class="text-sm text-text-muted">Select a router to configure its hotspot settings.</p>
    </div>

    <div v-else-if="loading" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <template v-else>
      <!-- Tabs -->
      <div class="flex items-center gap-1 border-b border-border mb-6 -mt-1">
        <button
          v-for="t in tabs" :key="t.key"
          class="px-3 py-2 text-sm font-medium border-b-2 transition-colors"
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
              <option value="USD">USD — US Dollar</option>
              <option value="EUR">EUR — Euro</option>
              <option value="GBP">GBP — British Pound</option>
              <option value="XOF">XOF — CFA Franc (West Africa)</option>
              <option value="XAF">XAF — CFA Franc (Central Africa)</option>
              <option value="NGN">NGN — Nigerian Naira</option>
              <option value="KES">KES — Kenyan Shilling</option>
              <option value="GHS">GHS — Ghanaian Cedi</option>
              <option value="ZAR">ZAR — South African Rand</option>
              <option value="TZS">TZS — Tanzanian Shilling</option>
              <option value="UGX">UGX — Ugandan Shilling</option>
              <option value="ETB">ETB — Ethiopian Birr</option>
              <option value="RWF">RWF — Rwandan Franc</option>
              <option value="MZN">MZN — Mozambican Metical</option>
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

      <!-- ── Tab: Cleanup ── -->
      <div v-else-if="tab === 'cleanup'" class="max-w-lg space-y-5">
        <div>
          <h3 class="text-sm font-semibold text-text-primary">Auto-cleanup</h3>
          <p class="text-sm text-text-muted mt-0.5">
            Expired hotspot users are automatically removed on a schedule.
            The scheduler is installed on the router during hotspot setup.
          </p>
        </div>

 <div class="flex items-center gap-2 p-3 border border-border rounded-lg text-xs text-text-secondary bg-surface" >
 <CheckCircleIcon class="size-4 shrink-0 text-green" />
          Scheduler <span class="font-mono text-text-primary mx-1">pikro-cleanup</span>
          {{ cleanup.installed ? 'active on router' : 'not yet installed — run hotspot setup' }}
        </div>

        <div class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Run every</span>
          <select v-model="cleanup.interval" class="input" @change="saveCleanup">
            <option value="10m">Every 10 minutes (testing)</option>
            <option value="1h">Every hour</option>
            <option value="1d">Daily (recommended)</option>
            <option value="7d">Weekly</option>
          </select>
        </div>

 <div v-if="deviceModeBlocked" class="p-4 border rounded-xl space-y-2 bg-amber/8 border-amber/20" >
          <div class="flex items-start gap-2">
 <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5 text-amber" />
            <div class="space-y-1">
 <p class="text-xs font-semibold text-amber" >Scheduler blocked by RouterOS device-mode</p>
              <p class="text-xs text-text-secondary">
                Your router runs RouterOS 7.17+ with the scheduler disabled. Run the command below in Winbox Terminal or SSH,
                then press the physical button on your router (or cold-reboot) to confirm.
              </p>
            </div>
          </div>
 <div class="border border-border rounded-lg px-3 py-2 font-mono text-xs text-text-primary select-all bg-base" >
            /system/device-mode/update scheduler=yes
          </div>
          <p class="text-xs text-text-muted">After confirming, re-run hotspot setup to install the scheduler.</p>
        </div>
 <p v-else-if="cleanupError" class="text-xs text-red" >{{ cleanupError }}</p>
      </div>

      <!-- ── Tab: Login page ── -->
      <div v-else-if="tab === 'login'" class="flex gap-6 items-start">
        <div class="w-72 shrink-0 space-y-4">
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary">Title</span>
            <input v-model="loginPage.title" class="input" placeholder="Sign in to continue" />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary">Subtitle</span>
            <input v-model="loginPage.subtitle" class="input" placeholder="$(hostname)" />
            <p class="text-xs text-text-muted">Leave blank to show the hotspot DNS name.</p>
          </label>
          <div class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary">Accent color</span>
            <div class="flex items-center gap-2">
              <input ref="colorPickerRef" type="color" v-model="loginPage.accentColor" class="sr-only" />
              <button
                type="button"
                class="h-9 w-9 shrink-0 rounded-lg border border-border cursor-pointer hover:scale-105 transition-transform"
                :style="{ background: loginPage.accentColor || '#111827' }"
                @click="colorPickerRef?.click()"
                title="Pick color"
              />
              <input v-model="loginPage.accentColor" class="input font-mono" placeholder="#111827" />
            </div>
            <p class="text-xs text-text-muted">Button and focus ring color.</p>
          </div>

 <p v-if="loginPageError" class="text-xs text-red" >{{ loginPageError }}</p>
 <p v-if="loginPageSaved" class="text-xs text-green" >Uploaded to router.</p>

          <div class="flex flex-col gap-2 pt-1">
            <button type="button" class="btn btn-primary btn-sm" :disabled="loginPageUploading" @click="uploadLoginPage">
              <span v-if="loginPageUploading" class="size-3.5 border-2 border-black/20 border-t-black rounded-full animate-spin" />
              Upload to router
            </button>
            <button type="button" class="btn btn-ghost btn-sm" @click="resetLoginPage">Reset to default</button>
          </div>
        </div>

        <div class="flex-1 min-w-0">
          <p class="text-xs font-medium text-text-muted mb-2">Preview</p>
 <div class="border border-border rounded-xl overflow-hidden bg-surface" style="height: 480px" >
            <iframe
              :srcdoc="loginPagePreview"
              class="w-full h-full border-0"
              sandbox=""
              title="Login page preview"
            />
          </div>
        </div>
      </div>

      <!-- ── Tab: Vouchers ── -->
      <div v-else-if="tab === 'voucher'" class="flex gap-6 items-start">
        <div class="w-72 shrink-0 space-y-4">
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary">Business name</span>
            <input v-model="voucher.businessName" class="input" placeholder="e.g. CyberCafé Lumière" />
            <p class="text-xs text-text-muted">Shown at the top of each voucher card.</p>
          </label>

          <div class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary">Show on voucher</span>
            <div class="flex flex-col gap-2 mt-1">
              <label class="flex items-center gap-2 text-xs text-text-secondary cursor-pointer">
                <input type="checkbox" v-model="voucher.showValidity" class="rounded" />
                Validity (e.g. 1 day, 4 hours)
              </label>
              <label class="flex items-center gap-2 text-xs text-text-secondary cursor-pointer">
                <input type="checkbox" v-model="voucher.showPrice" class="rounded" />
                Price (uses currency from General tab)
              </label>
            </div>
          </div>

 <p v-if="voucherSaved" class="text-xs text-green" >Voucher settings saved.</p>
 <p v-if="voucherError" class="text-xs text-red" >{{ voucherError }}</p>

          <button type="button" class="btn btn-primary btn-sm w-full justify-center" :disabled="voucherSaving" @click="saveVoucher">
            <span v-if="voucherSaving" class="size-3.5 border-2 border-black/20 border-t-black rounded-full animate-spin" />
            Save voucher settings
          </button>
        </div>

        <div class="flex-1 min-w-0">
          <p class="text-xs font-medium text-text-muted mb-2">Preview</p>
 <div class="border border-border rounded-xl p-5 grid grid-cols-2 gap-4 bg-surface" >
            <div
              v-for="sample in voucherSamples"
              :key="sample.name"
              class="border border-border rounded-lg p-3 flex flex-col gap-2 text-xs bg-base"
            >
              <div
                v-if="(voucher.showValidity) || (voucher.showPrice && effectiveCurrency)"
                class="flex items-baseline justify-between border-b border-border pb-2"
              >
                <span class="text-text-muted text-xs">{{ voucher.showValidity ? '1 day' : '' }}</span>
                <span v-if="voucher.showPrice && effectiveCurrency" class="font-bold text-text-primary">500 {{ effectiveCurrency }}</span>
              </div>
              <div class="grid gap-2 items-center" style="grid-template-columns: 1fr 1px 1fr">
                <div class="flex flex-col items-center gap-0.5">
                  <span class="text-xs text-text-muted uppercase tracking-wide">Username</span>
                  <span class="font-mono font-bold text-text-primary">{{ sample.name }}</span>
                </div>
                <div class="self-stretch border-r border-border" />
                <div class="flex flex-col items-center gap-0.5">
                  <span class="text-xs text-text-muted uppercase tracking-wide">Password</span>
                  <span class="font-mono font-bold text-text-primary">{{ sample.password }}</span>
                </div>
              </div>
              <div v-if="voucher.businessName" class="text-xs text-text-muted text-center border-t border-border pt-1.5 mt-auto">
                {{ voucher.businessName }}
              </div>
            </div>
          </div>
          <p class="text-xs text-text-muted mt-2">Sample data — actual values come from the selected profile when printing.</p>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { CheckCircleIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import {
  getHotspotSettings, putHotspotSettings, uploadLoginPage as apiUploadLoginPage,
  type HotspotSettings, type LoginPageSettings, type VoucherSettings,
  getCleanupScheduler, putCleanupScheduler,
} from '@/api'

const store = useRoutersStore()

const tab = ref<'general' | 'login' | 'voucher' | 'cleanup'>('general')
const tabs = [
  { key: 'general' as const,  label: 'General' },
  { key: 'login' as const,    label: 'Login page' },
  { key: 'voucher' as const,  label: 'Vouchers' },
  { key: 'cleanup' as const,  label: 'Cleanup' },
]

const loading = ref(false)
const saving = ref(false)
const settingsSaved = ref(false)
const settingsError = ref('')
const customCurrency = ref('')

const form = ref<HotspotSettings>({ hotspotName: '', dnsName: '', currency: '' })

const effectiveCurrency = computed(() =>
  form.value.currency === 'custom' ? customCurrency.value : form.value.currency
)

const cleanupSaving = ref(false)
const cleanupError = ref('')
const cleanup = ref({ installed: false, interval: '1d' })

const deviceModeBlocked = computed(() =>
  cleanupError.value.toLowerCase().includes('device-mode')
)

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
    form.value = { hotspotName: s.hotspotName ?? '', dnsName: s.dnsName ?? '', currency: s.currency ?? '' }
    const known = ['', 'USD', 'EUR', 'GBP', 'XOF', 'XAF', 'NGN', 'KES', 'GHS', 'ZAR', 'TZS', 'UGX', 'ETB', 'RWF', 'MZN']
    if (form.value.currency && !known.includes(form.value.currency)) {
      customCurrency.value = form.value.currency
      form.value.currency = 'custom'
    }
    if (s.loginPage) {
      loginPage.value = {
        title: s.loginPage.title ?? '',
        subtitle: s.loginPage.subtitle ?? '',
        accentColor: s.loginPage.accentColor ?? '#111827',
      }
    }
    if (s.voucher) {
      voucher.value = {
        businessName: s.voucher.businessName ?? '',
        showValidity: s.voucher.showValidity ?? true,
        showPrice: s.voucher.showPrice ?? true,
      }
    }
    cleanup.value = { installed: c.installed, interval: c.interval || '1d' }
  } catch (e: any) {
    settingsError.value = e?.response?.data?.error ?? e?.message ?? 'Failed to load settings'
  } finally {
    loading.value = false
  }
}

function currentSettings(): HotspotSettings {
  const currency = effectiveCurrency.value
  return { ...form.value, currency, loginPage: loginPage.value, voucher: voucher.value }
}

async function save() {
  if (!store.activeId) return
  saving.value = true; settingsSaved.value = false; settingsError.value = ''
  try {
    await putHotspotSettings(store.activeId, currentSettings())
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

const loginPage = ref<LoginPageSettings>({ title: '', subtitle: '', accentColor: '#111827' })
const loginPageUploading = ref(false)
const loginPageError = ref('')
const loginPageSaved = ref(false)
const colorPickerRef = ref<HTMLInputElement | null>(null)

const LOGIN_PAGE_TEMPLATE = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: ui-monospace, "Cascadia Code", Menlo, Consolas, monospace;
    background: #f9fafb;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1rem;
  }
  .card {
    background: #fff;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    padding: 2rem;
    width: 100%;
    max-width: 360px;
    box-shadow: 0 1px 3px rgba(0,0,0,.07);
  }
  .logo { display: flex; align-items: center; gap: .5rem; margin-bottom: 1.5rem; justify-content: center; }
  .logo-text { font-size: .875rem; font-weight: 600; color: #111827; letter-spacing: -.02em; }
  h1 { font-size: 1rem; font-weight: 600; color: #111827; margin-bottom: .25rem; text-align: center; }
  .subtitle { font-size: .75rem; color: #9ca3af; text-align: center; margin-bottom: 1.5rem; }
  label { display: flex; flex-direction: column; gap: .25rem; margin-bottom: .75rem; }
  .label-text { font-size: .75rem; font-weight: 500; color: #6b7280; }
  input[type="text"], input[type="password"] {
    width: 100%; padding: .5rem .75rem; font-size: .875rem; font-family: inherit;
    border: 1px solid #e5e7eb; border-radius: 8px; background: #fff; color: #111827;
    outline: none; transition: border-color .15s;
  }
  button[type="submit"] {
    width: 100%; padding: .625rem 1rem; background: __ACCENT__; color: #fff;
    border: none; border-radius: 8px; font-size: .875rem; font-weight: 500;
    font-family: inherit; cursor: pointer; margin-top: .25rem;
  }
  .footer { font-size: .7rem; color: #d1d5db; text-align: center; margin-top: 1.25rem; }
</style>
</head>
<body>
<div class="card">
  <div class="logo">
    <span class="logo-text">pikro</span>
  </div>
  <h1>__TITLE__</h1>
  <p class="subtitle">__SUBTITLE__</p>
  <form>
    <label><span class="label-text">Username</span><input type="text" /></label>
    <label><span class="label-text">Password</span><input type="password" /></label>
    <button type="submit">Sign In</button>
  </form>
  <p class="footer">Powered by Pikro</p>
</div>
</body>
</html>`

const loginPagePreview = computed(() => {
  const accent = loginPage.value.accentColor || '#111827'
  const title = loginPage.value.title || 'Sign in to continue'
  const subtitle = loginPage.value.subtitle || 'myspot.spot'
  return LOGIN_PAGE_TEMPLATE
    .replace('__ACCENT__', accent)
    .replace('__TITLE__', title)
    .replace('__SUBTITLE__', subtitle)
})

function resetLoginPage() {
  loginPage.value = { title: '', subtitle: '', accentColor: '#111827' }
}

async function uploadLoginPage() {
  if (!store.activeId) return
  if (!confirm('This will replace the current hotspot login page on the router. Continue?')) return
  loginPageUploading.value = true; loginPageError.value = ''; loginPageSaved.value = false
  try {
    await apiUploadLoginPage(store.activeId, loginPage.value)
    await putHotspotSettings(store.activeId, currentSettings())
    loginPageSaved.value = true
    setTimeout(() => { loginPageSaved.value = false }, 3000)
  } catch (e: any) {
    loginPageError.value = e?.response?.data?.error ?? e?.message ?? 'Failed to upload'
  } finally { loginPageUploading.value = false }
}

const voucher = ref<VoucherSettings>({ businessName: '', showValidity: true, showPrice: true })
const voucherSaving = ref(false)
const voucherSaved = ref(false)
const voucherError = ref('')

const voucherSamples = [
  { name: 'ab3f', password: 'ab3f' },
  { name: 'x9kz', password: 'x9kz' },
]

async function saveVoucher() {
  if (!store.activeId) return
  voucherSaving.value = true; voucherSaved.value = false; voucherError.value = ''
  try {
    await putHotspotSettings(store.activeId, currentSettings())
    voucherSaved.value = true
    setTimeout(() => { voucherSaved.value = false }, 3000)
  } catch (e: any) {
    voucherError.value = e?.response?.data?.error ?? e?.message ?? 'Failed to save'
  } finally { voucherSaving.value = false }
}

watch(() => store.activeId, load, { immediate: true })
</script>

