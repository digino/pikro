<template>
  <PageLayout title="Hotspot" subtitle="Sales Reports">
    <template #actions>
      <button class="btn btn-ghost" @click="window.print()">
        <PrinterIcon class="size-3.5" />
        Print
      </button>
    </template>

    <NoRouterSelected v-if="!store.activeId" />

    <div v-else class="flex flex-col gap-5">

      <!-- Controls -->
      <div class="flex items-center gap-3 flex-wrap">
        <div class="flex rounded-lg border border-border overflow-hidden">
          <button
            v-for="p in periods"
            :key="p.key"
            type="button"
            class="px-3 py-1.5 text-sm font-medium transition-colors"
            :class="period === p.key
              ? 'bg-accent text-base'
              : 'text-text-secondary hover:text-text-primary hover:bg-surface'"
            @click="period = p.key"
          >
            {{ p.label }}
          </button>
        </div>

        <select v-model="selectedYear" class="input-select">
          <option v-for="y in years" :key="y" :value="y">{{ y }}</option>
        </select>

        <select v-if="period === 'daily'" v-model="selectedMonth" class="input-select">
          <option v-for="(m, i) in MONTH_NAMES" :key="i" :value="i">{{ m }}</option>
        </select>

        <select v-model="filterProfile" class="input-select">
          <option value="">All profiles</option>
          <option v-for="name in allProfileNames" :key="name" :value="name">{{ name }}</option>
        </select>

        <div v-if="loading" class="flex items-center gap-1.5 text-xs text-text-muted ml-auto">
          <span class="spinner spinner--sm" /> Loading…
        </div>
        <p v-else-if="error" class="text-xs text-red ml-auto">{{ error }}</p>
      </div>

      <!-- No data yet state -->
      <EmptyState
        v-if="!loading && ledger.length === 0"
        bordered
        title="No sales recorded yet"
        message="Sales are recorded automatically each time you generate vouchers through Pikro."
      />

      <template v-else>
        <!-- Summary cards -->
        <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
          <div class="rounded-xl border border-border p-4 bg-surface">
            <p class="text-xs text-text-muted mb-1">Vouchers generated</p>
            <p class="font-mono text-2xl font-bold text-text-primary">{{ summary.generated }}</p>
          </div>
          <div class="rounded-xl border border-border p-4 bg-surface">
            <p class="text-xs text-text-muted mb-1">Revenue</p>
            <p class="font-mono text-2xl font-bold text-primary">{{ fmtAmount(summary.revenue) }}</p>
          </div>
          <div class="rounded-xl border border-border p-4 bg-surface">
            <p class="text-xs text-text-muted mb-1">Activated on router</p>
            <p class="font-mono text-2xl font-bold text-green">{{ activatedCount }}</p>
            <p class="text-xs text-text-muted mt-0.5">bytes-in &gt; 0</p>
          </div>
          <div class="rounded-xl border border-border p-4 bg-surface">
            <p class="text-xs text-text-muted mb-1">Unused inventory</p>
            <p class="font-mono text-2xl font-bold text-amber">{{ unusedCount }}</p>
            <p class="text-xs text-text-muted mt-0.5">on router, not used yet</p>
          </div>
        </div>

        <!-- Top profiles -->
        <div v-if="topProfiles.length > 1" class="grid grid-cols-2 gap-3 sm:grid-cols-4">
          <div
            v-for="tp in topProfiles"
            :key="tp.profile"
            class="rounded-xl border border-border p-4 bg-surface"
          >
            <p class="text-xs text-text-muted mb-1 truncate">{{ tp.profile }}</p>
            <p class="font-mono text-lg font-bold text-text-primary">{{ tp.count }}</p>
            <p class="text-xs text-text-muted mt-0.5">{{ fmtAmount(tp.revenue) }}</p>
          </div>
        </div>

        <!-- Breakdown table -->
        <div class="rounded-xl border border-border overflow-hidden bg-surface" id="report-table">
          <div class="flex items-center justify-between px-5 py-4 border-b border-border">
            <span class="font-semibold text-text-primary">
              {{ period === 'monthly' ? String(selectedYear) : `${MONTH_NAMES[selectedMonth]} ${selectedYear}` }}
            </span>
            <span class="text-xs text-text-muted">{{ activeCurrency || '' }}</span>
          </div>

          <EmptyState v-if="rows.length === 0" message="No vouchers generated in this period." />

          <table v-else class="w-full text-sm">
            <thead>
              <tr class="border-b border-border text-text-muted text-xs">
                <th class="px-5 py-3 text-left font-medium">{{ period === 'monthly' ? 'Month' : 'Date' }}</th>
                <th class="px-5 py-3 text-right font-medium">Generated</th>
                <th
                  v-for="name in visibleProfiles"
                  :key="name"
                  class="px-5 py-3 text-right font-medium"
                >{{ name }}</th>
                <th class="px-5 py-3 text-right font-medium">Revenue</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="row in rows"
                :key="row.label"
                class="border-b border-border/50 last:border-0 hover:bg-muted/30 transition-colors"
              >
                <td class="px-5 py-3 font-mono text-text-primary">{{ row.label }}</td>
                <td class="px-5 py-3 text-right font-mono text-text-secondary">{{ row.count }}</td>
                <td
                  v-for="name in visibleProfiles"
                  :key="name"
                  class="px-5 py-3 text-right font-mono text-text-muted"
                >{{ row.byProfile[name] ?? '—' }}</td>
                <td class="px-5 py-3 text-right font-mono font-semibold text-text-primary">
                  {{ fmtAmount(row.revenue) }}
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="border-t border-border bg-muted/20 text-xs font-semibold">
                <td class="px-5 py-3 text-text-secondary">Total</td>
                <td class="px-5 py-3 text-right font-mono text-text-primary">{{ summary.generated }}</td>
                <td
                  v-for="name in visibleProfiles"
                  :key="name"
                  class="px-5 py-3 text-right font-mono text-text-primary"
                >{{ totalByProfile[name] ?? '—' }}</td>
                <td class="px-5 py-3 text-right font-mono text-primary">{{ fmtAmount(summary.revenue) }}</td>
              </tr>
            </tfoot>
          </table>
        </div>

        <p class="text-xs text-text-muted">
          * Activated count reflects users currently on the router with usage recorded.
          Cleaned-up users are not counted.
        </p>
      </template>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { PrinterIcon } from '@heroicons/vue/24/outline'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import EmptyState from '@/components/EmptyState.vue'
import PageLayout from '@/components/PageLayout.vue'
import { useRoutersStore } from '@/stores/routers'
import { getSalesLedger, listHotspotUsers, type SaleEntry } from '@/api'

const store = useRoutersStore()

// ── Controls ──────────────────────────────────────────────
const MONTH_NAMES = [
  'January','February','March','April','May','June',
  'July','August','September','October','November','December',
]

const periods = [
  { key: 'daily',   label: 'Daily' },
  { key: 'monthly', label: 'Monthly' },
] as const
type Period = typeof periods[number]['key']

const now = new Date()
const period        = ref<Period>('daily')
const selectedYear  = ref(now.getFullYear())
const selectedMonth = ref(now.getMonth())
const filterProfile = ref('')

const years = computed(() => {
  const y = now.getFullYear()
  return [y, y - 1, y - 2]
})

// ── Data ──────────────────────────────────────────────────
const ledger      = ref<SaleEntry[]>([])
const liveUsers   = ref<Record<string, string>[]>([])
const loading     = ref(false)
const error       = ref('')

async function load() {
  if (!store.activeId) return
  loading.value = true
  error.value = ''
  try {
    const [entries, users] = await Promise.all([
      getSalesLedger(store.activeId, selectedYear.value),
      listHotspotUsers(store.activeId).catch(() => [] as Record<string, string>[]),
    ])
    ledger.value    = entries
    liveUsers.value = users as Record<string, string>[]
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? e?.message ?? 'Failed to load'
  } finally {
    loading.value = false
  }
}

watch(() => store.activeId, load, { immediate: true })
watch(selectedYear, load)

// ── Currency ──────────────────────────────────────────────
// Derived from ledger entries — avoids a separate settings fetch.
const activeCurrency = computed(() => {
  for (const e of ledger.value) {
    if (e.currency) return e.currency
  }
  return ''
})

// ── Live inventory stats ──────────────────────────────────
// "Activated" = users on the router that have recorded usage (bytes-in > 0).
// "Unused"    = users on the router with bytes-in == 0 or missing.
const activatedCount = computed(() =>
  liveUsers.value.filter(u => parseInt(u['bytes-in'] ?? '0') > 0).length
)
const unusedCount = computed(() =>
  liveUsers.value.filter(u => !(parseInt(u['bytes-in'] ?? '0') > 0) && u.disabled !== 'true').length
)

// ── Profile names from ledger ─────────────────────────────
const allProfileNames = computed(() => {
  const names = new Set<string>()
  for (const e of ledger.value) if (e.profile) names.add(e.profile)
  return [...names].sort()
})

// ── Filtered rows ─────────────────────────────────────────
interface Row {
  label: string
  count: number
  revenue: number
  byProfile: Record<string, number>
}

const filteredLedger = computed(() =>
  filterProfile.value
    ? ledger.value.filter(e => e.profile === filterProfile.value)
    : ledger.value
)

const rows = computed((): Row[] => {
  const bucket = new Map<string, Row>()

  for (const entry of filteredLedger.value) {
    const d = new Date(entry.at)
    if (d.getFullYear() !== selectedYear.value) continue
    if (period.value === 'daily' && d.getMonth() !== selectedMonth.value) continue

    const key = period.value === 'monthly'
      ? MONTH_NAMES[d.getMonth()]
      : `${selectedYear.value}-${String(d.getMonth() + 1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`

    if (!bucket.has(key)) bucket.set(key, { label: key, count: 0, revenue: 0, byProfile: {} })
    const row = bucket.get(key)!
    row.count += entry.count
    row.revenue += (parseFloat(entry.price) || 0) * entry.count
    row.byProfile[entry.profile] = (row.byProfile[entry.profile] ?? 0) + entry.count
  }

  const sorted = [...bucket.values()]
  if (period.value === 'monthly') {
    sorted.sort((a, b) => MONTH_NAMES.indexOf(a.label) - MONTH_NAMES.indexOf(b.label))
  } else {
    sorted.sort((a, b) => a.label.localeCompare(b.label))
  }
  return sorted
})

// ── Summary & top profiles ────────────────────────────────
const summary = computed(() => {
  let generated = 0, revenue = 0
  for (const row of rows.value) {
    generated += row.count
    revenue   += row.revenue
  }
  return { generated, revenue }
})

const totalByProfile = computed(() => {
  const totals: Record<string, number> = {}
  for (const row of rows.value) {
    for (const [name, cnt] of Object.entries(row.byProfile)) {
      totals[name] = (totals[name] ?? 0) + cnt
    }
  }
  return totals
})

// Top profiles ranked by count — shown as quick-glance cards above the table.
const topProfiles = computed(() => {
  return Object.entries(totalByProfile.value)
    .map(([profile, count]) => ({
      profile,
      count,
      revenue: filteredLedger.value
        .filter(e => e.profile === profile)
        .reduce((s, e) => s + (parseFloat(e.price) || 0) * e.count, 0),
    }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 4)
})

// Profile columns shown in the breakdown table.
const visibleProfiles = computed(() =>
  filterProfile.value ? [filterProfile.value] : allProfileNames.value
)

// ── Formatter ─────────────────────────────────────────────
function fmtAmount(n: number): string {
  if (!n) return '0'
  const s = n % 1 === 0 ? String(n) : n.toFixed(0)
  return activeCurrency.value ? `${s} ${activeCurrency.value}` : s
}

const window = globalThis.window
</script>
