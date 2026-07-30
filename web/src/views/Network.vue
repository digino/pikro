<template>
  <PageLayout title="Network" subtitle="Traffic overview">
    <NoRouterSelected v-if="!store.activeId" />

    <div v-else class="grid gap-4">
      <!-- bandwidth chart + interface traffic -->
      <div class="grid grid-cols-2 gap-4">
        <!-- Bandwidth chart -->
        <div class="rounded-xl border border-border p-5 bg-surface h-full">
          <div class="flex items-center gap-4 mb-4">
            <span class="font-semibold text-text-primary">Bandwidth</span>
            <div class="flex items-center gap-1.5">
              <span class="size-2 rounded-sm shrink-0" style="background:#22d3ee"/>
              <span class="text-xs text-text-secondary">Download</span>
              <span class="font-mono text-xs font-semibold text-text-primary">{{ curDown }} Mbps</span>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="size-2 rounded-sm shrink-0" style="background:#f59e0b"/>
              <span class="text-xs text-text-secondary">Upload</span>
              <span class="font-mono text-xs font-semibold text-text-primary">{{ curUp }} Mbps</span>
            </div>
          </div>
          <div class="h-60">
            <BandwidthChart :history="bwHistory" />
          </div>
        </div>

        <!-- Interface traffic -->
        <div class="rounded-xl border border-border p-5 bg-surface">
          <div class="flex items-center justify-between mb-4">
            <span class="font-semibold text-text-primary">Interfaces</span>
            <span v-if="trafficUpdatedAt" class="text-xs text-text-muted">{{ trafficUpdatedAt }}</span>
          </div>
          <div v-if="trafficLoading && traffic.length === 0" class="flex justify-center py-8">
            <span class="spinner spinner--sm" />
          </div>
          <div v-else-if="traffic.length === 0" class="text-xs text-text-muted py-8 text-center">
            No interface data
          </div>
          <table v-else class="w-full text-xs">
            <thead>
              <tr class="text-text-muted border-b border-border">
                <th class="text-left pb-2 font-medium">Interface</th>
                <th class="text-left pb-2 font-medium">IP</th>
                <th class="text-right pb-2 font-medium">RX</th>
                <th class="text-right pb-2 font-medium">TX</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="iface in traffic"
                :key="iface.name"
                class="border-b border-border/40 last:border-0"
              >
                <td class="py-1.5 font-mono text-text-primary">{{ iface.name }}</td>
                <td class="py-1.5 font-mono text-text-muted">{{ ifaceIP(iface.name) }}</td>
                <td class="py-1.5 text-right font-mono text-text-primary">{{ formatBps(iface['rx-bits-per-second']) }}</td>
                <td class="py-1.5 text-right font-mono text-text-secondary">{{ formatBps(iface['tx-bits-per-second']) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import { useRoutersStore } from '@/stores/routers'
import { getPollSnapshot } from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import BandwidthChart from '@/components/BandwidthChart.vue'

const store = useRoutersStore()

const traffic = ref<Record<string, string>[]>([])
const addresses = ref<Record<string, string>[]>([])
const trafficLoading = ref(false)
const trafficUpdatedAt = ref('')

const N = 46
interface BwPoint { down: number; up: number }
const bwHistory = ref<BwPoint[]>(Array.from({ length: N }, () => ({ down: 0, up: 0 })))

const curDown = computed(() => bwHistory.value[N - 1].down.toFixed(2))
const curUp = computed(() => bwHistory.value[N - 1].up.toFixed(2))

async function poll() {
  if (!store.activeId) return
  trafficLoading.value = true
  try {
    const snap = await getPollSnapshot(store.activeId)
    if (snap.traffic.length > 0) {
      traffic.value = snap.traffic.filter((i: Record<string, string>) => i.name)
      trafficUpdatedAt.value = new Date().toLocaleTimeString()
      const totalRx = snap.traffic.reduce((s: number, i: Record<string, string>) =>
        s + (parseInt(i['rx-bits-per-second'] ?? '0') || 0), 0)
      const totalTx = snap.traffic.reduce((s: number, i: Record<string, string>) =>
        s + (parseInt(i['tx-bits-per-second'] ?? '0') || 0), 0)
      bwHistory.value = [...bwHistory.value.slice(1), { down: totalRx / 1_000_000, up: totalTx / 1_000_000 }]
    }
    if (snap.addresses) addresses.value = snap.addresses
  } catch { /* non-critical */ } finally {
    trafficLoading.value = false
  }
}

let pollTimer: ReturnType<typeof setInterval>

function startTimers() {
  pollTimer = setInterval(() => { if (!document.hidden) poll() }, 3000)
}

function stopTimers() {
  clearInterval(pollTimer)
}

watch(() => store.activeId, async (id) => {
  stopTimers()
  bwHistory.value = Array.from({ length: N }, () => ({ down: 0, up: 0 }))
  traffic.value = []
  addresses.value = []
  if (!id) return
  await poll()
  startTimers()
}, { immediate: true })

onMounted(() => {
  document.addEventListener('visibilitychange', () => {
    if (!document.hidden && store.activeId) poll()
  })
})

onUnmounted(stopTimers)

function ifaceIP(name: string): string {
  const match = addresses.value.find(a => a.interface === name)
  return match ? match.address : '—'
}

function formatBps(val: string | undefined): string {
  const n = parseInt(val ?? '0') || 0
  if (n >= 1_000_000_000) return (n / 1_000_000_000).toFixed(2) + ' Gbps'
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(2) + ' Mbps'
  if (n >= 1_000) return (n / 1_000).toFixed(1) + ' Kbps'
  return n + ' bps'
}
</script>
