<template>
  <div class="flex gap-6 items-start">
    <div class="w-72 shrink-0 space-y-5">

      <button
        type="button"
        class="w-full flex items-center justify-center gap-1.5 px-3 py-2 rounded-lg border text-sm font-medium transition-colors"
        :class="showingLive
          ? 'border-text-primary bg-muted text-text-primary'
          : 'border-border text-text-secondary hover:border-text-muted'"
        :disabled="loadingLive"
        @click="toggleLive"
      >
        <span v-if="loadingLive" class="size-3 border-2 border-text-muted/30 border-t-text-muted rounded-full animate-spin" />
        {{ showingLive ? 'Back to editor preview' : 'View currently used' }}
      </button>
      <p v-if="liveError" class="text-xs text-red -mt-3">{{ liveError }}</p>

      <!-- Template picker -->
      <div class="flex flex-col gap-2">
        <span class="text-sm font-medium text-text-secondary">Template</span>
        <div class="flex flex-col gap-2">
          <button
            v-for="t in templates"
            :key="t.key"
            type="button"
            class="flex items-start gap-3 p-3 rounded-lg border text-left transition-colors"
            :class="loginPage.template === t.key
              ? 'border-text-primary bg-muted'
              : 'border-border hover:border-text-muted'"
            @click="selectTemplate(t.key)"
          >
            <span class="text-xl leading-none mt-0.5 shrink-0">{{ t.icon }}</span>
            <div>
              <div class="text-sm font-semibold text-text-primary">{{ t.label }}</div>
              <div class="text-xs text-text-secondary mt-0.5">{{ t.description }}</div>
            </div>
          </button>
        </div>
      </div>

      <div class="border-t border-border pt-4 space-y-4">
        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Title</span>
          <input v-model="loginPage.title" class="input" :placeholder="hotspotName || 'Sign in to continue'" />
          <p class="text-xs text-text-muted">Leave blank to use your hotspot name.</p>
        </label>
        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Subtitle</span>
          <input v-model="loginPage.subtitle" class="input" placeholder="$(hostname)" />
          <p class="text-xs text-text-muted">Leave blank to show the hotspot DNS name.</p>
        </label>
      </div>

      <p v-if="error" class="text-xs text-red">{{ error }}</p>
      <p v-if="saved" class="text-xs text-green">Uploaded to router.</p>

      <div class="flex flex-col gap-2 pt-1">
        <button type="button" class="btn btn-primary" :disabled="uploading" @click="upload">
          <span v-if="uploading" class="size-3.5 border-2 border-black/20 border-t-black rounded-full animate-spin" />
          Upload to router
        </button>
        <button type="button" class="btn btn-ghost" @click="reset">Reset to default</button>
      </div>
    </div>

    <div class="flex-1 min-w-0">
      <p class="text-xs font-medium text-text-muted mb-2">
        {{ showingLive ? 'Currently used on router' : 'Preview' }}
      </p>
      <div class="border border-border rounded-xl overflow-hidden bg-surface" style="height: 520px">
        <iframe
          :srcdoc="showingLive ? liveHTML : preview"
          class="w-full h-full border-0"
          sandbox=""
          title="Login page preview"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoutersStore } from '@/stores/routers'
import { useToastStore } from '@/stores/toast'
import {
  getHotspotSettings, putHotspotSettings, uploadLoginPage as apiUploadLoginPage,
  getLoginPageHTML,
  type LoginPageSettings,
} from '@/api'
import {
  MINIMAL_TEMPLATE, WAVE_TEMPLATE, CARD_TEMPLATE,
  MINIMAL_LOGOUT_TEMPLATE, MINIMAL_STATUS_TEMPLATE, MINIMAL_REDIRECT_TEMPLATE,
  WAVE_LOGOUT_TEMPLATE, WAVE_STATUS_TEMPLATE, WAVE_REDIRECT_TEMPLATE,
  CARD_LOGOUT_TEMPLATE, CARD_STATUS_TEMPLATE, CARD_REDIRECT_TEMPLATE,
} from './loginPageTemplates'

const props = defineProps<{ hotspotName: string }>()

const store = useRoutersStore()
const toast = useToastStore()

const loginPage = ref<LoginPageSettings>({
  title: '',
  subtitle: '',
  template: 'minimal',
})
const uploading = ref(false)
const error = ref('')
const saved = ref(false)

const showingLive = ref(false)
const loadingLive = ref(false)
const liveError = ref('')
const liveHTML = ref('')

const TEMPLATE_SOURCE: Record<NonNullable<LoginPageSettings['template']>, string> = {
  minimal: MINIMAL_TEMPLATE,
  wave: WAVE_TEMPLATE,
  card: CARD_TEMPLATE,
}

// logout/status/redirect don't have a template picker of their own — they
// always match whichever login template is currently selected, so the whole
// hotspot page set stays visually consistent (alogin/error stay fixed and
// are never customized, handled entirely on the Go side).
const LOGOUT_SOURCE: Record<NonNullable<LoginPageSettings['template']>, string> = {
  minimal: MINIMAL_LOGOUT_TEMPLATE,
  wave: WAVE_LOGOUT_TEMPLATE,
  card: CARD_LOGOUT_TEMPLATE,
}
const STATUS_SOURCE: Record<NonNullable<LoginPageSettings['template']>, string> = {
  minimal: MINIMAL_STATUS_TEMPLATE,
  wave: WAVE_STATUS_TEMPLATE,
  card: CARD_STATUS_TEMPLATE,
}
const REDIRECT_SOURCE: Record<NonNullable<LoginPageSettings['template']>, string> = {
  minimal: MINIMAL_REDIRECT_TEMPLATE,
  wave: WAVE_REDIRECT_TEMPLATE,
  card: CARD_REDIRECT_TEMPLATE,
}

const templates: { key: LoginPageSettings['template']; icon: string; label: string; description: string }[] = [
  { key: 'minimal', icon: '◈', label: 'Minimal', description: 'Flat gray background, light card — sharp and focused' },
  { key: 'wave',    icon: '〜', label: 'Wave',    description: 'Light rounded card with a decorative wave accent' },
  { key: 'card',    icon: '▭', label: 'Card',    description: 'Classic centered card — traditional login form look' },
]

function init(lp?: LoginPageSettings) {
  if (lp) {
    loginPage.value = {
      title:    lp.title    ?? '',
      subtitle: lp.subtitle ?? '',
      template: lp.template ?? 'minimal',
    }
  }
}

defineExpose({ init })

function renderFrom(source: Record<NonNullable<LoginPageSettings['template']>, string>) {
  const title = loginPage.value.title || props.hotspotName || 'My Hotspot'
  const subtitle = loginPage.value.subtitle || 'myspot.spot'
  const tpl = source[loginPage.value.template ?? 'minimal']
  return tpl
    .replace(/__TITLE__/g, title)
    .replace(/__SUBTITLE__/g, subtitle)
}

function renderTemplate() {
  return renderFrom(TEMPLATE_SOURCE)
}

const preview = computed(renderTemplate)

function selectTemplate(key: LoginPageSettings['template']) {
  loginPage.value.template = key
  showingLive.value = false
}

async function toggleLive() {
  if (showingLive.value) {
    showingLive.value = false
    return
  }
  if (!store.activeId) return
  loadingLive.value = true; liveError.value = ''
  try {
    liveHTML.value = await getLoginPageHTML(store.activeId)
    showingLive.value = true
  } catch (e: any) {
    liveError.value = e?.response?.data?.error ?? e?.message ?? 'Failed to fetch'
  } finally { loadingLive.value = false }
}

function reset() {
  loginPage.value = { title: '', subtitle: '', template: 'minimal' }
  showingLive.value = false
}

async function upload() {
  if (!store.activeId) return
  if (!confirm('This will replace the current hotspot login page on the router. Continue?')) return
  uploading.value = true; error.value = ''; saved.value = false
  try {
    const existing = await getHotspotSettings(store.activeId)
    await apiUploadLoginPage(store.activeId, {
      ...loginPage.value,
      html: renderTemplate(),
      logoutHtml: renderFrom(LOGOUT_SOURCE),
      statusHtml: renderFrom(STATUS_SOURCE),
      redirectHtml: renderFrom(REDIRECT_SOURCE),
    })
    await putHotspotSettings(store.activeId, { ...existing, loginPage: loginPage.value })
    showingLive.value = false
    saved.value = true
    toast.success('Login page uploaded', 'The new page is now live on the router.')
    setTimeout(() => { saved.value = false }, 3000)
  } catch (e: any) {
    const message = e?.response?.data?.error ?? e?.message ?? 'Failed to upload'
    error.value = message
    toast.error('Upload failed', message)
  } finally { uploading.value = false }
}
</script>
