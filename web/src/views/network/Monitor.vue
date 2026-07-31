<template>
  <PageLayout title="Network" subtitle="Monitor">
    <div class="grid gap-4">
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
          <EmptyState v-else-if="traffic.length === 0" message="No interface data" />
          <table v-else class="w-full text-xs">
            <thead>
              <tr class="text-text-muted border-b border-border">
                <th class="text-left pb-2 font-medium">Interface</th>
                <th class="text-left pb-2 font-medium">IP</th>
                <th class="text-left pb-2 font-medium">Status</th>
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
                <td class="py-1.5">
                  <StatusBadge
                    v-if="ifaceStatus(iface.name)"
                    :color="ifaceStatus(iface.name) === 'up' ? 'green' : ifaceStatus(iface.name) === 'down' ? 'red' : 'muted'"
                    :label="ifaceStatus(iface.name) === 'up' ? 'Up' : ifaceStatus(iface.name) === 'down' ? 'Down' : 'Disabled'"
                  />
                  <span v-else class="text-text-muted">—</span>
                </td>
                <td class="py-1.5 text-right font-mono text-text-primary">{{ formatBps(iface['rx-bits-per-second']) }}</td>
                <td class="py-1.5 text-right font-mono text-text-secondary">{{ formatBps(iface['tx-bits-per-second']) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Recent logs -->
      <div
        class="rounded-xl border border-border p-5 grid gap-3 bg-surface"
        style="grid-template-rows: auto 1fr"
      >
        <div class="flex items-center justify-between">
          <span class="font-semibold text-text-primary">Recent logs</span>
          <RouterLink to="/hotspot/logs" class="btn btn-ghost btn-sm"
            >View all</RouterLink
          >
        </div>
        <div
          v-if="logsLoading && logs.length === 0"
          class="flex justify-center py-4"
        >
          <span class="spinner spinner--sm" />
        </div>
        <EmptyState v-else-if="logs.length === 0" size="sm" message="No log entries" />
        <div v-else class="grid" style="align-content: start">
          <div
            v-for="(entry, i) in logs"
            :key="i"
            class="flex items-start gap-2 py-1.5 border-b border-border/40 last:border-0"
          >
            <span
              class="shrink-0 mt-1 size-1.5 rounded-full"
              :style="`background:${logColor(entry.topics)}`"
            />
            <div class="min-w-0 flex-1">
              <p class="text-xs text-text-primary truncate">
                {{ entry.message }}
              </p>
              <p class="text-xs text-text-muted font-mono">
                {{ entry.time }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import EmptyState from '@/components/EmptyState.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { useRoutersStore } from '@/stores/routers'
import { getPollSnapshot, getSystemLogs } from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import BandwidthChart from '@/components/BandwidthChart.vue'

const store = useRoutersStore()

const traffic = ref<Record<string, string>[]>([])
const addresses = ref<Record<string, string>[]>([])
const interfaces = ref<Record<string, string>[]>([])
const trafficLoading = ref(false)
const trafficUpdatedAt = ref('')

const logs = ref<Record<string, string>[]>([])
const logsLoading = ref(false)

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
    if (snap.interfaces) interfaces.value = snap.interfaces
  } catch { /* non-critical */ } finally {
    trafficLoading.value = false
  }
}

async function loadLogs() {
  if (!store.activeId) return
  logsLoading.value = true
  try {
    const all = await getSystemLogs(store.activeId)
    logs.value = all
      .filter((e) => (e.topics ?? '').toLowerCase().includes('hotspot'))
      .slice(0, 8)
  } catch {
    // non-critical
  } finally {
    logsLoading.value = false
  }
}

let pollTimer: ReturnType<typeof setInterval>
let logTimer: ReturnType<typeof setInterval>

function startTimers() {
  pollTimer = setInterval(() => { if (!document.hidden) poll() }, 5000)
  logTimer = setInterval(() => { if (!document.hidden) loadLogs() }, 30_000)
}

function stopTimers() {
  clearInterval(pollTimer)
  clearInterval(logTimer)
}

watch(() => store.activeId, async (id) => {
  stopTimers()
  bwHistory.value = Array.from({ length: N }, () => ({ down: 0, up: 0 }))
  traffic.value = []
  addresses.value = []
  interfaces.value = []
  logs.value = []
  if (!id) return
  await Promise.all([poll(), loadLogs()])
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

function ifaceStatus(name: string): 'up' | 'down' | 'disabled' | null {
  const iface = interfaces.value.find(i => i.name === name)
  if (!iface) return null
  if (iface.disabled === 'true') return 'disabled'
  return iface.running === 'true' ? 'up' : 'down'
}

function logColor(topics: string | undefined): string {
  if (!topics) return 'var(--color-text-muted)'
  const t = topics.toLowerCase()
  if (t.includes('error') || t.includes('critical')) return 'var(--color-red)'
  if (t.includes('warning')) return 'var(--color-amber)'
  if (t.includes('info')) return 'var(--color-green)'
  return 'var(--color-text-muted)'
}

function formatBps(val: string | undefined): string {
  const n = parseInt(val ?? '0') || 0
  if (n >= 1_000_000_000) return (n / 1_000_000_000).toFixed(2) + ' Gbps'
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(2) + ' Mbps'
  if (n >= 1_000) return (n / 1_000).toFixed(1) + ' Kbps'
  return n + ' bps'
}
</script>
