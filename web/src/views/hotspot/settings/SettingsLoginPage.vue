<template>
  <div class="flex gap-6 items-start">
    <div class="w-72 shrink-0 space-y-5">

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
            @click="loginPage.template = t.key"
          >
            <span class="text-xl leading-none mt-0.5 shrink-0">{{ t.icon }}</span>
            <div>
              <div class="text-xs font-semibold text-text-primary">{{ t.label }}</div>
              <div class="text-xs text-text-muted mt-0.5">{{ t.description }}</div>
            </div>
          </button>
        </div>
      </div>

      <div class="border-t border-border pt-4 space-y-4">
        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Title</span>
          <input v-model="loginPage.title" class="input" :placeholder="businessName || 'Sign in to continue'" />
          <p class="text-xs text-text-muted">Leave blank to use your business name.</p>
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
          <p class="text-xs text-text-muted">Button and highlight color.</p>
        </div>
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
      <p class="text-xs font-medium text-text-muted mb-2">Preview</p>
      <div class="border border-border rounded-xl overflow-hidden bg-surface" style="height: 520px">
        <iframe
          :srcdoc="preview"
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
import {
  getHotspotSettings, putHotspotSettings, uploadLoginPage as apiUploadLoginPage,
  type LoginPageSettings,
} from '@/api'
import { MINIMAL_TEMPLATE, BOLD_TEMPLATE } from './loginPageTemplates'

const props = defineProps<{ businessName: string }>()

const store = useRoutersStore()

const loginPage = ref<LoginPageSettings>({
  title: '',
  subtitle: '',
  accentColor: '#111827',
  template: 'minimal',
})
const uploading = ref(false)
const error = ref('')
const saved = ref(false)
const colorPickerRef = ref<HTMLInputElement | null>(null)

const templates: { key: LoginPageSettings['template']; icon: string; label: string; description: string }[] = [
  { key: 'minimal', icon: '◈', label: 'Minimal', description: 'Dark mesh background, glass card — sharp and focused' },
  { key: 'bold',    icon: '▨', label: 'Bold',    description: 'Accent header with business name, white form — mobile-first' },
]

function init(lp?: LoginPageSettings) {
  if (lp) {
    loginPage.value = {
      title:       lp.title       ?? '',
      subtitle:    lp.subtitle    ?? '',
      accentColor: lp.accentColor ?? '#111827',
      template:    lp.template    ?? 'minimal',
    }
  }
}

defineExpose({ init })

const preview = computed(() => {
  const accent = loginPage.value.accentColor || '#111827'
  const title = loginPage.value.title || props.businessName || 'My Hotspot'
  const subtitle = loginPage.value.subtitle || 'myspot.spot'
  const tpl = loginPage.value.template === 'bold' ? BOLD_TEMPLATE : MINIMAL_TEMPLATE
  return tpl
    .replace(/__ACCENT__/g, accent)
    .replace(/__TITLE__/g, title)
    .replace(/__SUBTITLE__/g, subtitle)
})

function reset() {
  loginPage.value = { title: '', subtitle: '', accentColor: '#111827', template: 'minimal' }
}

async function upload() {
  if (!store.activeId) return
  if (!confirm('This will replace the current hotspot login page on the router. Continue?')) return
  uploading.value = true; error.value = ''; saved.value = false
  try {
    const existing = await getHotspotSettings(store.activeId)
    await apiUploadLoginPage(store.activeId, loginPage.value)
    await putHotspotSettings(store.activeId, { ...existing, loginPage: loginPage.value })
    saved.value = true
    setTimeout(() => { saved.value = false }, 3000)
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? e?.message ?? 'Failed to upload'
  } finally { uploading.value = false }
}
</script>
