<template>
  <div class="flex gap-6 items-start">
    <div class="w-72 shrink-0 space-y-5">
      <div>
        <h3 class="text-sm font-semibold text-text-primary">Voucher templates</h3>
        <p class="text-sm text-text-muted mt-0.5">
          Vouchers always show validity and price when the profile has them.
          The template itself is picked each time you print, from the print dialog.
        </p>
      </div>

      <!-- Template picker (preview only) -->
      <div class="flex flex-col gap-2">
        <span class="text-sm font-medium text-text-secondary">Preview a template</span>
        <div class="flex flex-col gap-2">
          <button
            v-for="t in VOUCHER_TEMPLATES"
            :key="t.key"
            type="button"
            class="flex items-start gap-3 p-3 rounded-lg border text-left transition-colors"
            :class="preview === t.key
              ? 'border-text-primary bg-muted'
              : 'border-border hover:border-text-muted'"
            @click="preview = t.key"
          >
            <div>
              <div class="text-xs font-semibold text-text-primary">{{ t.label }}</div>
              <div class="text-xs text-text-muted mt-0.5">{{ t.description }}</div>
            </div>
          </button>
        </div>
      </div>
    </div>

    <!-- Preview -->
    <div class="flex-1 min-w-0">
      <p class="text-xs font-medium text-text-muted mb-2">Preview</p>

      <!-- Classic: dense, heavy borders, zero gap -->
      <div v-if="preview === 'classic'" class="border border-border rounded-xl p-5 bg-surface overflow-hidden">
        <div class="grid grid-cols-4 font-mono text-xs">
          <div
            v-for="sample in voucherSamples"
            :key="sample.name"
            class="border-2 border-text-primary p-2 flex flex-col gap-1 -mt-px -ml-px"
          >
            <div class="flex justify-between"><span class="text-[8px] uppercase text-text-muted">User</span><span class="font-bold">{{ sample.name }}</span></div>
            <div class="flex justify-between"><span class="text-[8px] uppercase text-text-muted">Pass</span><span class="font-bold">{{ sample.password }}</span></div>
            <div class="text-[8px] border-t border-text-primary pt-1 text-center">1 day · 500 {{ effectiveCurrency }}</div>
          </div>
        </div>
      </div>

      <!-- Modern: compact cards, heavy border, small gap -->
      <div v-else-if="preview === 'modern'" class="border border-border rounded-xl p-5 bg-surface grid grid-cols-4 gap-2">
        <div
          v-for="sample in voucherSamples"
          :key="sample.name"
          class="border-2 border-text-primary rounded p-2.5 flex flex-col gap-1.5 text-xs bg-base"
        >
          <div v-if="hotspotName" class="text-[8px] font-bold text-center uppercase border-b border-text-primary pb-1">{{ hotspotName }}</div>
          <div class="grid gap-1.5 items-center" style="grid-template-columns: 1fr 1px 1fr">
            <div class="flex flex-col items-center gap-0.5">
              <span class="text-text-muted uppercase tracking-wide text-[8px]">User</span>
              <span class="font-mono font-bold text-text-primary">{{ sample.name }}</span>
            </div>
            <div class="self-stretch border-r border-text-primary" />
            <div class="flex flex-col items-center gap-0.5">
              <span class="text-text-muted uppercase tracking-wide text-[8px]">Pass</span>
              <span class="font-mono font-bold text-text-primary">{{ sample.password }}</span>
            </div>
          </div>
          <div class="text-[8px] text-center border-t border-text-primary pt-1">1 day · 500 {{ effectiveCurrency }}</div>
        </div>
      </div>

      <!-- Business: roomier, QR on the side -->
      <div v-else class="border border-border rounded-xl p-5 bg-surface grid grid-cols-2 gap-3">
        <div
          v-for="sample in voucherSamples"
          :key="sample.name"
          class="border-2 border-text-primary rounded p-3 flex items-center gap-3 text-xs bg-base"
        >
          <div class="flex-1 min-w-0 flex flex-col gap-1">
            <div v-if="hotspotName" class="text-[9px] font-bold uppercase border-b border-text-primary pb-1 mb-0.5">{{ hotspotName }}</div>
            <div class="flex justify-between"><span class="text-[8px] uppercase text-text-muted">Username</span><span class="font-mono font-bold">{{ sample.name }}</span></div>
            <div class="flex justify-between"><span class="text-[8px] uppercase text-text-muted">Password</span><span class="font-mono font-bold">{{ sample.password }}</span></div>
            <div class="text-[8px] border-t border-text-primary pt-1">1 day · 500 {{ effectiveCurrency }}</div>
          </div>
          <div class="size-14 shrink-0 rounded bg-muted flex items-center justify-center text-text-muted" title="QR code (scan to open the login page)">
            <QrCodeIcon class="size-8" />
          </div>
        </div>
      </div>

      <p class="text-xs text-text-muted mt-2">Sample data — actual values come from the selected profile when printing.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { QrCodeIcon } from '@heroicons/vue/24/outline'
import { VOUCHER_TEMPLATES, type VoucherTemplate } from '@/utils/voucherTemplates'

defineProps<{ hotspotName: string; effectiveCurrency: string }>()

const preview = ref<VoucherTemplate['key']>('classic')

const voucherSamples = [
  { name: 'ab3f', password: 'ab3f' },
  { name: 'x9kz', password: 'x9kz' },
  { name: 'p7mw', password: 'p7mw' },
  { name: 'q2nt', password: 'q2nt' },
]
</script>
