<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <PageLayout v-else :title="t('nav.network')" :subtitle="t('nav.monitor')" icon="monitor">
    <template #actions>
      <button
        type="button"
        class="btn"
        :class="live && 'btn-ink shadow-[3px_3px_0_var(--color-primary)]'"
        :aria-pressed="live"
        :title="t('monitor.liveHint')"
        @click="live = !live"
      >
        <span
          class="size-2.25 rounded-full"
          :class="live ? 'bg-primary shadow-[0_0_0_4px_color-mix(in_srgb,var(--color-primary)_30%,transparent)]' : 'bg-text-muted'"
        />
        {{ live ? t('monitor.live') : t('monitor.paused') }}
      </button>
    </template>

    <section class="grid grid-cols-[repeat(auto-fit,minmax(220px,1fr))] gap-5">
      <div class="card bg-primary text-primary-fg px-5 py-4.5 flex flex-col gap-1.5">
        <span class="eyebrow text-primary-fg">{{ t('monitor.downloadNow') }}</span>
        <span class="type-display text-[44px] leading-none">
          {{ now.down }} <span class="text-lg">Mbps</span>
        </span>
        <span class="text-[13px]">{{ t('monitor.on', { iface: selected }) }}</span>
      </div>
      <div class="card px-5 py-4.5 flex flex-col gap-1.5">
        <span class="eyebrow">{{ t('monitor.uploadNow') }}</span>
        <span class="type-display text-[44px] leading-none">
          {{ now.up }} <span class="text-lg">Mbps</span>
        </span>
        <span class="text-[13px] text-text-muted">{{ t('monitor.on', { iface: selected }) }}</span>
      </div>
      <div class="card px-5 py-4.5 flex flex-col gap-1.5">
        <span class="eyebrow">{{ t('monitor.peakSince', { time: since }) }}</span>
        <span class="type-display text-[44px] leading-none">
          {{ peak ? peak.mbps.toFixed(1) : '0.0' }} <span class="text-lg">Mbps</span>
        </span>
        <span class="text-[13px] text-text-muted">
          {{ peak ? t('monitor.peakAt', { time: peak.at }) : t('monitor.noPeak') }}
        </span>
      </div>
      <div class="card bg-ink text-paper px-5 py-4.5 flex flex-col gap-2 shadow-[5px_5px_0_var(--color-primary)]">
        <span class="eyebrow text-paper/75">{{ t('monitor.wan') }}</span>
        <span class="flex items-center gap-2 text-lg font-bold">
          <span
            class="size-2.5 rounded-full"
            :class="wan.online ? 'bg-online shadow-[0_0_0_4px_color-mix(in_srgb,var(--color-online)_25%,transparent)]' : 'bg-danger'"
          />
          {{ wan.online ? t('monitor.online') : t('monitor.offline') }}
        </span>
        <span class="font-mono text-[15px]">{{ wan.ip }}</span>
        <span class="text-[13px] text-paper/75">{{ wan.detail }}</span>
      </div>
    </section>

    <section class="card p-6 flex flex-col gap-4">
      <div class="flex flex-wrap justify-between items-center gap-3">
        <h2 class="type-display m-0 text-2xl">{{ t('monitor.bandwidth', { iface: selected }) }}</h2>
        <div class="flex gap-3.5 text-[13px]">
          <span class="inline-flex items-center gap-1.5">
            <span class="w-4.5 h-2.5 border-2 border-ink bg-primary" />{{ t('monitor.download') }}
          </span>
          <span class="inline-flex items-center gap-1.5">
            <span class="w-4.5 border-t-[2.5px] border-dashed border-ink" />{{ t('monitor.upload') }}
          </span>
        </div>
      </div>
      <div
        class="h-65"
        role="img"
        :aria-label="t('monitor.chartLabel', { iface: selected, down: now.down, up: now.up })"
      >
        <BandwidthChart :history="history[selected] ?? EMPTY_HISTORY" />
      </div>
    </section>

    <section class="flex flex-col gap-3.5">
      <div class="flex flex-wrap items-baseline gap-3">
        <h2 class="type-display m-0 text-[28px]">{{ t('monitor.interfaces') }}</h2>
        <span class="text-sm text-text-muted">{{ t('monitor.pickHint') }}</span>
      </div>

      <div v-if="loading && !cards.length" class="flex justify-center py-8">
        <span class="spinner spinner--sm" />
      </div>
      <EmptyState v-else-if="!cards.length" :message="t('monitor.noInterfaces')" />
      <div v-else class="grid grid-cols-[repeat(auto-fill,minmax(250px,1fr))] gap-4.5">
        <button
          v-for="c in cards"
          :key="c.name"
          type="button"
          :aria-pressed="c.name === selected"
          class="flex flex-col gap-3 p-4 text-left rounded-card border-2 border-ink bg-paper text-ink cursor-pointer transition-shadow"
          :class="[
            c.status === 'disabled' && 'border-dashed',
            c.name === selected
              ? 'shadow-[6px_6px_0_var(--color-primary),6px_6px_0_2px_var(--color-ink)]'
              : c.status === 'disabled' ? '' : 'shadow-[5px_5px_0_var(--color-shadow)]',
          ]"
          @click="picked = c.name"
        >
          <span class="flex items-center gap-3">
            <span
              aria-hidden="true"
              class="flex-none size-11.5 rounded-control border-2 border-ink flex items-center justify-center"
              :class="c.status === 'up' ? 'bg-ink text-primary' : 'bg-ground text-text-muted'"
            >
              <svg v-if="c.wifi" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M2 8.8a15.5 15.5 0 0 1 20 0M5 12.6a11 11 0 0 1 14 0M8.5 16.2a6 6 0 0 1 7 0" /><circle cx="12" cy="20" r="1.4" fill="currentColor" /></svg>
              <svg v-else width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linejoin="round"><path d="M4 5h16v11h-4v3H8v-3H4z" /><path d="M8 9v2M11 9v2M13 9v2M16 9v2" stroke-linecap="round" /></svg>
            </span>
            <span class="flex flex-col gap-0.5 min-w-0 flex-1">
              <span class="font-mono font-bold text-[15px] truncate">{{ c.name }}</span>
              <span class="text-xs text-text-muted">{{ c.role }}</span>
            </span>
            <span class="chip" :class="STATUS_CHIP[c.status]">
              <span class="size-2 rounded-full" :class="STATUS_DOT[c.status]" />{{ t(`monitor.status.${c.status}`) }}
            </span>
          </span>
          <span class="font-mono text-[13px]">{{ c.ip }}</span>
          <span class="grid grid-cols-2 gap-2 border-t-[1.5px] border-hairline pt-2.5">
            <span class="flex flex-col gap-0.5">
              <span class="eyebrow">{{ t('monitor.rx') }}</span>
              <span class="font-mono text-sm font-semibold">{{ c.rx }}</span>
            </span>
            <span class="flex flex-col gap-0.5">
              <span class="eyebrow">{{ t('monitor.tx') }}</span>
              <span class="font-mono text-sm font-semibold">{{ c.tx }}</span>
            </span>
          </span>
        </button>
      </div>
    </section>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoutersStore } from '@/stores/routers'
import { getPollSnapshot, type PollSnapshot } from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import EmptyState from '@/components/EmptyState.vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import BandwidthChart from '@/components/BandwidthChart.vue'
import { formatBytes } from '@/utils/bytes'

type Status = 'up' | 'down' | 'disabled'
interface BwPoint { down: number; up: number }

const STATUS_CHIP: Record<Status, string> = {
  up: 'bg-primary text-primary-fg',
  down: 'bg-danger-wash text-danger-text',
  disabled: 'bg-hairline text-text-secondary border-dashed',
}
const STATUS_DOT: Record<Status, string> = {
  up: 'bg-ink',
  down: 'bg-danger rotate-45 rounded-[1px]',
  disabled: 'border-[1.5px] border-dashed border-text-secondary',
}

// 5 minutes of points at one poll every 5 s (BandwidthChart's step).
const POINTS = 60
const POLL_MS = 5000
const EMPTY_HISTORY: BwPoint[] = Array.from({ length: POINTS }, () => ({ down: 0, up: 0 }))

const { t, locale } = useI18n()
const store = useRoutersStore()

const snap = ref<PollSnapshot | null>(null)
const loading = ref(false)
const live = ref(true)
const picked = ref<string | null>(null)
const history = ref<Record<string, BwPoint[]>>({})
const peaks = ref<Record<string, { mbps: number; at: string }>>({})
const since = ref('')

const clockTime = (d = new Date()) =>
  d.toLocaleTimeString(locale.value, { hour: '2-digit', minute: '2-digit' })

const interfaces = computed(() =>
  (snap.value?.interfaces ?? []).filter(i => i.name && i.type !== 'loopback'),
)
const wanName = computed(() => snap.value?.wanInterface ?? '')

const selected = computed(() =>
  picked.value
  ?? (wanName.value || interfaces.value.find(i => i.running === 'true')?.name)
  ?? '—',
)

const now = computed(() => {
  const points = history.value[selected.value] ?? EMPTY_HISTORY
  const last = points[points.length - 1]
  return { down: last.down.toFixed(1), up: last.up.toFixed(1) }
})
const peak = computed(() => peaks.value[selected.value])

function statusOf(i: Record<string, string>): Status {
  if (i.disabled === 'true') return 'disabled'
  return i.running === 'true' ? 'up' : 'down'
}

function addressOf(name: string) {
  return snap.value?.addresses.find(a => a.interface === name)?.address
}

function roleOf(i: Record<string, string>) {
  const s = snap.value
  const hotspots = s?.hotspotInterfaces ?? []
  const bridge = s?.bridgePorts[i.name]
  if (i.name === wanName.value) return t('monitor.roles.wan')
  if (hotspots.includes(i.name)) return t('monitor.roles.guestNetwork')
  if (bridge && hotspots.includes(bridge)) {
    return t(isWifi(i) ? 'monitor.roles.guestWifi' : 'monitor.roles.guestPort', { bridge })
  }
  if (bridge) return t('monitor.roles.bridgePort', { bridge })
  if (i.comment) return i.comment
  return addressOf(i.name) ? t('monitor.roles.local') : t('monitor.roles.spare')
}

function isWifi(i: Record<string, string>) {
  return /^(wlan|wifi)/.test(i.type ?? '')
}

// Internet first, then the guest network and its ports, then the rest in
// the router's own order.
function rankOf(name: string) {
  const hotspots = snap.value?.hotspotInterfaces ?? []
  if (name === wanName.value) return 0
  if (hotspots.includes(name)) return 1
  const bridge = snap.value?.bridgePorts[name]
  return bridge && hotspots.includes(bridge) ? 2 : 3
}

const cards = computed(() =>
  [...interfaces.value]
    .sort((a, b) => rankOf(a.name) - rankOf(b.name))
    .map(i => ({
      name: i.name,
      role: roleOf(i),
      status: statusOf(i),
      wifi: isWifi(i),
      ip: addressOf(i.name) ?? t('monitor.noAddress'),
      rx: formatBytes(Number(i['rx-byte'])),
      tx: formatBytes(Number(i['tx-byte'])),
    })),
)

const wan = computed(() => {
  const i = interfaces.value.find(x => x.name === wanName.value)
  if (!i) return { online: false, ip: '—', detail: t('monitor.noWan') }
  const uptime = linkUptime(i['last-link-up-time'], snap.value?.clock)
  return {
    online: statusOf(i) === 'up',
    ip: addressOf(i.name)?.split('/')[0] ?? t('monitor.noAddress'),
    detail: uptime
      ? t('monitor.viaUp', { iface: i.name, uptime })
      : t('monitor.via', { iface: i.name }),
  }
})

// RouterOS 7 writes dates as 2026-10-29, RouterOS 6 as oct/29/2026. Both the
// link time and the clock are router-local, so their difference is the uptime.
const MONTHS = ['jan', 'feb', 'mar', 'apr', 'may', 'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec']
function routerMs(date = '', time = '') {
  const v7 = /^(\d{4})-(\d{2})-(\d{2})$/.exec(date)
  const v6 = /^([a-z]{3})\/(\d{2})\/(\d{4})$/.exec(date)
  const hms = /^(\d{2}):(\d{2}):(\d{2})$/.exec(time)
  if (!hms || !(v7 || v6)) return null
  const [y, m, d] = v7
    ? [+v7[1], +v7[2] - 1, +v7[3]]
    : [+v6![3], MONTHS.indexOf(v6![1]), +v6![2]]
  return Date.UTC(y, m, d, +hms[1], +hms[2], +hms[3])
}

function linkUptime(linkUp?: string, clock?: Record<string, string>) {
  const [date, time] = (linkUp ?? '').split(' ')
  const from = routerMs(date, time)
  const to = routerMs(clock?.date, clock?.time)
  if (from === null || to === null || to < from) return ''
  const mins = Math.floor((to - from) / 60_000)
  const d = Math.floor(mins / 1440), h = Math.floor(mins / 60) % 24, m = mins % 60
  if (d) return `${d}d ${h}h`
  if (h) return `${h}h ${m}m`
  return `${m}m`
}

function record(s: PollSnapshot) {
  const next: Record<string, BwPoint[]> = {}
  for (const tr of s.traffic) {
    if (!tr.name) continue
    const rx = (parseInt(tr['rx-bits-per-second'] ?? '0') || 0) / 1_000_000
    const tx = (parseInt(tr['tx-bits-per-second'] ?? '0') || 0) / 1_000_000
    // The router counts from its own side: the WAN receives what guests
    // download, every other port sends it to them.
    const point = tr.name === s.wanInterface ? { down: rx, up: tx } : { down: tx, up: rx }
    next[tr.name] = [...(history.value[tr.name] ?? EMPTY_HISTORY).slice(1), point]
    if (point.down > (peaks.value[tr.name]?.mbps ?? 0)) {
      peaks.value[tr.name] = { mbps: point.down, at: clockTime() }
    }
  }
  history.value = next
}

async function poll() {
  const id = store.activeId
  if (!id) return
  loading.value = true
  try {
    const s = await getPollSnapshot(id)
    if (id !== store.activeId) return
    snap.value = s
    record(s)
  } catch { /* non-critical: the next poll retries */ } finally {
    loading.value = false
  }
}

let pollTimer: ReturnType<typeof setInterval> | undefined

function stopTimer() {
  clearInterval(pollTimer)
  pollTimer = undefined
}

function startTimer() {
  stopTimer()
  if (!live.value || !store.activeId) return
  poll()
  pollTimer = setInterval(() => { if (!document.hidden) poll() }, POLL_MS)
}

watch(() => store.activeId, () => {
  snap.value = null
  history.value = {}
  peaks.value = {}
  picked.value = null
  since.value = clockTime()
  startTimer()
}, { immediate: true })

watch(live, startTimer)

function onVisibilityChange() {
  if (!document.hidden && live.value && store.activeId) poll()
}

onMounted(() => document.addEventListener('visibilitychange', onVisibilityChange))
onUnmounted(() => {
  stopTimer()
  document.removeEventListener('visibilitychange', onVisibilityChange)
})
</script>
