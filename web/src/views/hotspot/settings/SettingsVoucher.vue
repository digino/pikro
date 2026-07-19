<template>
  <div class="flex gap-6 items-start">
    <div class="w-72 shrink-0 space-y-5">

      <!-- Layout picker -->
      <div class="flex flex-col gap-2">
        <span class="text-sm font-medium text-text-secondary">Layout</span>
        <div class="flex flex-col gap-2">
          <button
            v-for="l in voucherLayouts"
            :key="l.key"
            type="button"
            class="flex items-start gap-3 p-3 rounded-lg border text-left transition-colors"
            :class="voucher.layout === l.key
              ? 'border-text-primary bg-muted'
              : 'border-border hover:border-text-muted'"
            @click="voucher.layout = l.key"
          >
            <span class="text-lg leading-none mt-0.5">{{ l.icon }}</span>
            <div>
              <div class="text-xs font-semibold text-text-primary">{{ l.label }}</div>
              <div class="text-xs text-text-muted mt-0.5">{{ l.description }}</div>
            </div>
          </button>
        </div>
      </div>

      <!-- Show toggles -->
      <div class="flex flex-col gap-1">
        <span class="text-sm font-medium text-text-secondary">Show on voucher</span>
        <div class="flex flex-col gap-3 mt-2">
          <label class="flex items-center justify-between gap-3 cursor-pointer">
            <span class="text-sm text-text-secondary">Validity (e.g. 1 day, 4 hours)</span>
            <SwitchRoot v-model="voucher.showValidity" class="relative inline-flex h-5 w-9 shrink-0 rounded-full transition-colors data-[state=checked]:bg-text-primary data-[state=unchecked]:bg-border focus:outline-none">
              <SwitchThumb class="pointer-events-none block h-4 w-4 rounded-full bg-white shadow transition-transform translate-x-0.5 data-[state=checked]:translate-x-4.5" />
            </SwitchRoot>
          </label>
          <label class="flex items-center justify-between gap-3 cursor-pointer">
            <span class="text-sm text-text-secondary">Price (currency from General tab)</span>
            <SwitchRoot v-model="voucher.showPrice" class="relative inline-flex h-5 w-9 shrink-0 rounded-full transition-colors data-[state=checked]:bg-text-primary data-[state=unchecked]:bg-border focus:outline-none">
              <SwitchThumb class="pointer-events-none block h-4 w-4 rounded-full bg-white shadow transition-transform translate-x-0.5 data-[state=checked]:translate-x-4.5" />
            </SwitchRoot>
          </label>
        </div>
      </div>

      <p v-if="saved" class="text-xs text-green">Voucher settings saved.</p>
      <p v-if="error" class="text-xs text-red">{{ error }}</p>

      <button type="button" class="btn btn-primary" :disabled="saving" @click="save">
        <span v-if="saving" class="size-3.5 border-2 border-black/20 border-t-black rounded-full animate-spin" />
        Save voucher settings
      </button>
    </div>

    <!-- Preview -->
    <div class="flex-1 min-w-0">
      <p class="text-xs font-medium text-text-muted mb-2">Preview</p>

      <!-- Card layout: 2-column grid of mini cards -->
      <div v-if="voucher.layout === 'card' || !voucher.layout" class="border border-border rounded-xl p-5 grid grid-cols-2 gap-3 bg-surface">
        <div
          v-for="sample in voucherSamples"
          :key="sample.name"
          class="border border-border rounded-lg p-3 flex flex-col gap-2 text-xs bg-base"
        >
          <div
            v-if="voucher.showValidity || (voucher.showPrice && effectiveCurrency)"
            class="flex items-baseline justify-between border-b border-border pb-1.5"
          >
            <span class="text-text-muted">{{ voucher.showValidity ? '1 day' : '' }}</span>
            <span v-if="voucher.showPrice && effectiveCurrency" class="font-bold text-text-primary">500 {{ effectiveCurrency }}</span>
          </div>
          <div class="grid gap-2 items-center" style="grid-template-columns: 1fr 1px 1fr">
            <div class="flex flex-col items-center gap-0.5">
              <span class="text-text-muted uppercase tracking-wide" style="font-size: 9px">Username</span>
              <span class="font-mono font-bold text-text-primary">{{ sample.name }}</span>
            </div>
            <div class="self-stretch border-r border-border" />
            <div class="flex flex-col items-center gap-0.5">
              <span class="text-text-muted uppercase tracking-wide" style="font-size: 9px">Password</span>
              <span class="font-mono font-bold text-text-primary">{{ sample.password }}</span>
            </div>
          </div>
          <div class="flex items-center justify-between border-t border-border pt-1.5 mt-auto">
            <span v-if="businessName" class="text-text-muted" style="font-size: 9px">{{ businessName }}</span>
            <span class="text-text-muted ml-auto" style="font-size: 9px">#{{ voucherSamples.indexOf(sample) + 1 }}</span>
          </div>
        </div>
      </div>

      <!-- Ticket layout: auto-fill grid (2-up on A4, fills page) -->
      <div v-else-if="voucher.layout === 'ticket'" class="border border-border rounded-xl p-5 bg-surface grid grid-cols-2 gap-3">
        <div
          v-for="sample in voucherSamples"
          :key="sample.name"
          class="rounded-lg border border-border bg-base overflow-hidden text-xs"
        >
          <!-- Header band -->
          <div class="bg-muted border-b border-border px-4 py-2 flex items-center justify-between gap-2">
            <span class="font-bold text-text-primary text-sm">{{ businessName || 'WiFi Voucher' }}</span>
            <span v-if="voucher.showPrice && effectiveCurrency" class="font-bold text-text-primary">500 {{ effectiveCurrency }}</span>
            <span class="text-text-muted ml-auto" style="font-size: 9px">#{{ voucherSamples.indexOf(sample) + 1 }}</span>
          </div>
          <!-- Credentials -->
          <div class="px-4 py-3 flex flex-col gap-2">
            <div class="flex items-center justify-between">
              <span class="text-text-muted uppercase tracking-wide" style="font-size: 9px">Username</span>
              <span class="font-mono font-bold text-text-primary">{{ sample.name }}</span>
            </div>
            <div class="flex items-center justify-between border-t border-border pt-2">
              <span class="text-text-muted uppercase tracking-wide" style="font-size: 9px">Password</span>
              <span class="font-mono font-bold text-text-primary">{{ sample.password }}</span>
            </div>
            <div v-if="voucher.showValidity" class="flex items-center justify-between border-t border-border pt-2">
              <span class="text-text-muted uppercase tracking-wide" style="font-size: 9px">Valid for</span>
              <span class="text-text-secondary text-xs">1 day</span>
            </div>
          </div>
        </div>
      </div>

      <p class="text-xs text-text-muted mt-2">Sample data — actual values come from the selected profile when printing.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { SwitchRoot, SwitchThumb } from 'reka-ui'
import { useRoutersStore } from '@/stores/routers'
import { getHotspotSettings, putHotspotSettings, type VoucherSettings } from '@/api'

const props = defineProps<{ businessName: string; effectiveCurrency: string }>()

const store = useRoutersStore()

const voucher = ref<VoucherSettings>({ businessName: '', showValidity: true, showPrice: true, layout: 'card' })
const saving = ref(false)
const saved = ref(false)
const error = ref('')

const voucherLayouts: { key: VoucherSettings['layout']; icon: string; label: string; description: string }[] = [
  { key: 'card',   icon: '▦', label: 'Card',   description: '6 per page — compact grid, maximises A4 space' },
  { key: 'ticket', icon: '▣', label: 'Ticket', description: '2 per page — large text, easy to read and hand out' },
]

const voucherSamples = [
  { name: 'ab3f', password: 'ab3f' },
  { name: 'x9kz', password: 'x9kz' },
  { name: 'p7mw', password: 'p7mw' },
  { name: 'q2nt', password: 'q2nt' },
]

function init(v?: VoucherSettings) {
  if (v) {
    voucher.value = {
      businessName: v.businessName ?? '',
      showValidity: v.showValidity ?? true,
      showPrice: v.showPrice ?? true,
      layout: v.layout ?? 'card',
    }
  }
}

defineExpose({ init })

async function save() {
  if (!store.activeId) return
  saving.value = true; saved.value = false; error.value = ''
  try {
    const existing = await getHotspotSettings(store.activeId)
    await putHotspotSettings(store.activeId, { ...existing, voucher: voucher.value })
    saved.value = true
    setTimeout(() => { saved.value = false }, 3000)
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? e?.message ?? 'Failed to save'
  } finally { saving.value = false }
}
</script>
