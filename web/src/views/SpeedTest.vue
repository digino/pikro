<template>
  <PageLayout title="Speed Test" subtitle="Internet download speed measured from the router">
    <div v-if="!store.activeId" class="flex flex-col items-center justify-center py-16 gap-3">
      <BoltIcon class="size-8 text-text-muted" />
      <p class="text-sm text-text-secondary">Select a router first.</p>
    </div>

    <template v-else>
 <div class="border border-border rounded-xl p-5 max-w-md bg-surface" >
        <p class="text-sm text-text-muted mb-4 leading-relaxed">
          The test runs <strong class="text-text-secondary">from the router</strong> — not your browser.
          The router must have internet access to reach the target.
        </p>
        <div class="flex flex-col gap-1 mb-4">
          <span class="text-sm font-medium text-text-secondary">File size</span>
          <div class="flex gap-2">
            <button
              v-for="s in FILE_SIZES"
              :key="s.value"
              type="button"
              class="flex-1 py-1.5 text-sm font-medium rounded-lg border transition-colors"
              :class="fileSize === s.value
                ? 'border-text-primary bg-muted text-text-primary'
                : 'border-border text-text-muted hover:border-muted hover:text-text-secondary'"
              @click="fileSize = s.value"
            >
              {{ s.label }}
            </button>
          </div>
          <p class="text-xs text-text-muted">Larger file = more accurate on fast connections</p>
        </div>
        <button class="btn btn-primary" :disabled="running" @click="run">
          <span v-if="running" class="size-4 border-2 border-black/20 border-t-black rounded-full animate-spin" />
          <BoltIcon v-else class="size-4" />
          {{ running ? 'Testing…' : 'Run Speed Test' }}
        </button>
        <p v-if="error" class="text-sm mt-3 text-red">{{ error }}</p>
      </div>

      <div v-if="result" class="flex flex-col gap-3 mt-4 max-w-md">
 <div class="border border-border rounded-xl p-5 text-center bg-surface" >
          <p class="text-xs text-text-muted mb-1">Download speed</p>
          <p class="text-4xl font-bold text-text-primary">{{ formatSpeed(result['rx-speed']) }}</p>
          <p class="text-xs text-text-muted mt-1">10 MB file · {{ result['duration'] || '—' }}</p>
        </div>
        <p class="text-xs text-text-muted text-center">
          Measured from the router via
          <span class="font-mono text-text-secondary">speedtest.tele2.net</span>
        </p>
      </div>
    </template>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { BoltIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import { runSpeedTest } from '@/api'
import PageLayout from '@/components/PageLayout.vue'

const FILE_SIZES = [
  { label: '10 MB', value: '10' },
  { label: '50 MB', value: '50' },
  { label: '100 MB', value: '100' },
]

const store = useRoutersStore()
const fileSize = ref('10')
const running = ref(false)
const error = ref('')
const result = ref<Record<string, string> | null>(null)

async function run() {
  if (!store.activeId) return
  running.value = true
  error.value = ''
  result.value = null
  try {
    result.value = await runSpeedTest(store.activeId, fileSize.value)
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? 'Speed test failed'
  } finally {
    running.value = false
  }
}

function formatSpeed(val: string) {
  const n = parseInt(val)
  if (isNaN(n)) return '—'
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + ' Mbps'
  if (n >= 1_000) return (n / 1_000).toFixed(0) + ' Kbps'
  return n + ' bps'
}
</script>
