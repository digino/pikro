<template>
  <PageLayout title="Hotspot" subtitle="Sales Reports">
    <template #actions>
      <button class="btn btn-ghost" @click="exportCsv">
        <ArrowDownTrayIcon class="size-3.5" />
        Export CSV
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
            :class="
              period === p.key
                ? 'bg-accent text-base'
                : 'text-text-secondary hover:text-text-primary hover:bg-surface'
            "
            @click="period = p.key"
          >
            {{ p.label }}
          </button>
        </div>

        <AppSelect
          v-model="selectedYear"
          :options="years.map((y) => ({ value: y, label: String(y) }))"
          numeric
        />

        <AppSelect
          v-if="period === 'daily'"
          v-model="selectedMonth"
          :options="MONTH_NAMES.map((m, i) => ({ value: i, label: m }))"
          numeric
        />

        <AppSelect
          v-model="filterProfile"
          :options="
            allProfileNames.map((name) => ({ value: name, label: name }))
          "
          placeholder="All profiles"
        />

        <div
          v-if="loading"
          class="flex items-center gap-1.5 text-xs text-text-muted ml-auto"
        >
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
            <p class="text-sm font-semibold mb-1">Vouchers generated</p>
            <div class="flex items-baseline gap-2">
              <p class="font-mono text-2xl font-bold text-text-primary">
                {{ summary.generated }}
              </p>
              <TrendBadge :delta="generatedDelta" />
            </div>
            <p class="text-xs text-text-secondary mt-0.5">
              vs. {{ prevSummary.generated }} prior
            </p>
          </div>
          <div class="rounded-xl border border-border p-4 bg-surface">
            <p class="text-sm font-semibold mb-1">Revenue</p>
            <div class="flex items-baseline gap-2">
              <p class="font-mono text-2xl font-bold">
                {{ fmtAmount(summary.revenue) }}
              </p>
              <TrendBadge :delta="revenueDelta" />
            </div>
            <p class="text-xs text-text-secondary mt-0.5">
              vs. {{ fmtAmount(prevSummary.revenue) }} prior
            </p>
          </div>
          <div class="rounded-xl border border-border p-4 bg-surface">
            <p class="text-sm font-semibold mb-1">Activated on router</p>
            <p class="font-mono text-2xl font-bold">
              {{ activatedCount }}
            </p>
            <p class="text-xs text-text-secondary font-medium mt-0.5">
              bytes-in &gt; 0
            </p>
          </div>
          <div class="rounded-xl border border-border p-4 bg-surface">
            <p class="text-sm font-medium mb-1">Unused inventory</p>
            <p class="font-mono text-2xl font-bold">
              {{ unusedCount }}
            </p>
            <p class="text-xs text-text-secondary mt-0.5">
              on router, not used yet
            </p>
          </div>
        </div>

        <!-- Top profiles -->
        <div
          v-if="topProfiles.length > 1"
          class="grid grid-cols-2 gap-3 lg:grid-cols-4 xl:grid-cols-5"
        >
          <div
            v-for="tp in topProfiles"
            :key="tp.profile"
            class="rounded-xl border border-border p-4 bg-surface"
          >
            <p class="text-sm text-text-secondary font-medium mb-1 truncate">
              {{ tp.profile }}
            </p>
            <p class="font-mono text-lg font-bold text-text-primary">
              {{ tp.count }}
            </p>
            <p class="text-xs text-text-muted mt-0.5">
              {{ fmtAmount(tp.revenue) }}
            </p>
          </div>
        </div>

        <!-- Chart -->
        <div class="rounded-xl border border-border p-5 bg-surface">
          <div class="flex items-center justify-between mb-4">
            <span class="font-semibold text-text-primary">
              {{ period === "monthly" ? "Monthly" : "Daily" }} sales
            </span>
            <span class="text-xs text-text-muted">{{
              activeCurrency || ""
            }}</span>
          </div>
          <div class="h-96">
            <SalesBarChart :points="chartPoints" :currency="activeCurrency" />
          </div>
        </div>

        <!-- Breakdown table -->
        <div class="rounded-xl border border-border overflow-hidden bg-surface">
          <div
            class="flex items-center justify-between px-5 py-4 border-b border-border"
          >
            <span class="font-semibold text-text-primary">
              {{
                period === "monthly"
                  ? String(selectedYear)
                  : `${MONTH_NAMES[selectedMonth]} ${selectedYear}`
              }}
            </span>
            <span class="text-xs text-text-muted">{{
              activeCurrency || ""
            }}</span>
          </div>

          <EmptyState
            v-if="rows.length === 0"
            message="No vouchers generated in this period."
          />

          <table v-else class="w-full text-sm">
            <thead>
              <tr class="border-b border-border text-text-muted text-xs">
                <th class="px-5 py-3 text-left font-medium">
                  {{ period === "monthly" ? "Month" : "Date" }}
                </th>
                <th class="px-5 py-3 text-right font-medium">Generated</th>
                <th
                  v-for="name in visibleProfiles"
                  :key="name"
                  class="px-5 py-3 text-right font-medium"
                >
                  {{ name }}
                </th>
                <th class="px-5 py-3 text-right font-medium">Revenue</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="row in rows"
                :key="row.label"
                class="border-b border-border/50 last:border-0 hover:bg-muted/30 transition-colors"
              >
                <td class="px-5 py-3 font-mono text-text-primary">
                  {{ row.label }}
                </td>
                <td class="px-5 py-3 text-right font-mono text-text-secondary">
                  {{ row.count }}
                </td>
                <td
                  v-for="name in visibleProfiles"
                  :key="name"
                  class="px-5 py-3 text-right font-mono text-text-muted"
                >
                  {{ row.byProfile[name] ?? "—" }}
                </td>
                <td
                  class="px-5 py-3 text-right font-mono font-semibold text-text-primary"
                >
                  {{ fmtAmount(row.revenue) }}
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr
                class="border-t border-border bg-muted/20 text-xs font-semibold"
              >
                <td class="px-5 py-3 text-text-secondary">Total</td>
                <td class="px-5 py-3 text-right font-mono text-text-primary">
                  {{ summary.generated }}
                </td>
                <td
                  v-for="name in visibleProfiles"
                  :key="name"
                  class="px-5 py-3 text-right font-mono text-text-primary"
                >
                  {{ totalByProfile[name] ?? "—" }}
                </td>
                <td class="px-5 py-3 text-right font-mono text-primary">
                  {{ fmtAmount(summary.revenue) }}
                </td>
              </tr>
            </tfoot>
          </table>
        </div>

        <p class="text-xs text-text-secondary font-medium">
          * Activated count reflects users currently on the router with usage
          recorded. Cleaned-up users are not counted.
        </p>
      </template>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { ArrowDownTrayIcon } from "@heroicons/vue/24/outline";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import EmptyState from "@/components/EmptyState.vue";
import AppSelect from "@/components/AppSelect.vue";
import SalesBarChart from "@/components/SalesBarChart.vue";
import TrendBadge from "@/components/TrendBadge.vue";
import PageLayout from "@/components/PageLayout.vue";
import { useRoutersStore } from "@/stores/routers";
import { getSalesLedger, listHotspotUsers, type SaleEntry } from "@/api";
import { formatCompactAmount } from "@/utils/currencies";
import { downloadCsv } from "@/utils/csv";

const store = useRoutersStore();

// ── Controls ──────────────────────────────────────────────
const MONTH_NAMES = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

const periods = [
  { key: "daily", label: "Daily" },
  { key: "monthly", label: "Monthly" },
] as const;
type Period = (typeof periods)[number]["key"];

const now = new Date();
const period = ref<Period>("daily");
const selectedYear = ref(now.getFullYear());
const selectedMonth = ref(now.getMonth());
const filterProfile = ref("");

const years = computed(() => {
  const y = now.getFullYear();
  return [y, y - 1, y - 2];
});

// ── Data ──────────────────────────────────────────────────
// Both the selected year and the previous year are fetched so month-over-month
// (Jan vs. prior Dec) and year-over-year comparisons work across boundaries.
const ledger = ref<SaleEntry[]>([]);
const prevYearLedger = ref<SaleEntry[]>([]);
const liveUsers = ref<Record<string, string>[]>([]);
const loading = ref(false);
const error = ref("");

async function load() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    const [entries, prevEntries, users] = await Promise.all([
      getSalesLedger(store.activeId, selectedYear.value),
      getSalesLedger(store.activeId, selectedYear.value - 1).catch(() => []),
      listHotspotUsers(store.activeId).catch(
        () => [] as Record<string, string>[],
      ),
    ]);
    ledger.value = entries;
    prevYearLedger.value = prevEntries;
    liveUsers.value = users as Record<string, string>[];
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? e?.message ?? "Failed to load";
  } finally {
    loading.value = false;
  }
}

watch(() => store.activeId, load, { immediate: true });
watch(selectedYear, load);

// ── Currency ──────────────────────────────────────────────
// Derived from ledger entries — avoids a separate settings fetch.
const activeCurrency = computed(() => {
  for (const e of ledger.value) {
    if (e.currency) return e.currency;
  }
  return "";
});

// ── Live inventory stats ──────────────────────────────────
// "Activated" = users on the router that have recorded usage (bytes-in > 0).
// "Unused"    = users on the router with bytes-in == 0 or missing.
const activatedCount = computed(
  () =>
    liveUsers.value.filter((u) => parseInt(u["bytes-in"] ?? "0") > 0).length,
);
const unusedCount = computed(
  () =>
    liveUsers.value.filter(
      (u) => !(parseInt(u["bytes-in"] ?? "0") > 0) && u.disabled !== "true",
    ).length,
);

// ── Profile names from ledger ─────────────────────────────
const allProfileNames = computed(() => {
  const names = new Set<string>();
  for (const e of ledger.value) if (e.profile) names.add(e.profile);
  return [...names].sort();
});

// ── Filtered rows ─────────────────────────────────────────
interface Row {
  label: string;
  count: number;
  revenue: number;
  byProfile: Record<string, number>;
}

function filterByProfile(entries: SaleEntry[]): SaleEntry[] {
  return filterProfile.value
    ? entries.filter((e) => e.profile === filterProfile.value)
    : entries;
}

const filteredLedger = computed(() => filterByProfile(ledger.value));

function bucketEntries(entries: SaleEntry[]): Row[] {
  const bucket = new Map<string, Row>();

  for (const entry of entries) {
    const d = new Date(entry.at);
    if (d.getFullYear() !== selectedYear.value) continue;
    if (period.value === "daily" && d.getMonth() !== selectedMonth.value)
      continue;

    const key =
      period.value === "monthly"
        ? MONTH_NAMES[d.getMonth()]
        : `${selectedYear.value}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;

    if (!bucket.has(key))
      bucket.set(key, { label: key, count: 0, revenue: 0, byProfile: {} });
    const row = bucket.get(key)!;
    row.count += entry.count;
    row.revenue += (parseFloat(entry.price) || 0) * entry.count;
    row.byProfile[entry.profile] =
      (row.byProfile[entry.profile] ?? 0) + entry.count;
  }

  const sorted = [...bucket.values()];
  if (period.value === "monthly") {
    sorted.sort(
      (a, b) => MONTH_NAMES.indexOf(a.label) - MONTH_NAMES.indexOf(b.label),
    );
  } else {
    sorted.sort((a, b) => a.label.localeCompare(b.label));
  }
  return sorted;
}

const rows = computed(() => bucketEntries(filteredLedger.value));

// ── Previous-period comparison ────────────────────────────
// Daily view (a specific month) compares to the previous month.
// Monthly view (a specific year) compares to the previous year.
const prevRows = computed(() => {
  if (period.value === "monthly") {
    return bucketEntriesForYear(
      filterByProfile(prevYearLedger.value),
      selectedYear.value - 1,
    );
  }
  // Previous month, possibly in the previous year (Jan -> prior Dec).
  const prevMonth = selectedMonth.value === 0 ? 11 : selectedMonth.value - 1;
  const prevYear =
    selectedMonth.value === 0 ? selectedYear.value - 1 : selectedYear.value;
  const source =
    prevYear === selectedYear.value ? ledger.value : prevYearLedger.value;
  return bucketEntriesForMonth(filterByProfile(source), prevYear, prevMonth);
});

// bucketEntries scoped to an explicit year/month rather than the reactive
// selectedYear/selectedMonth — needed for the previous-period comparison,
// which must always look at prevYear/prevMonth regardless of what's selected.
function bucketEntriesForYear(entries: SaleEntry[], year: number): Row[] {
  const bucket = new Map<string, Row>();
  for (const entry of entries) {
    const d = new Date(entry.at);
    if (d.getFullYear() !== year) continue;
    const key = MONTH_NAMES[d.getMonth()];
    if (!bucket.has(key))
      bucket.set(key, { label: key, count: 0, revenue: 0, byProfile: {} });
    const row = bucket.get(key)!;
    row.count += entry.count;
    row.revenue += (parseFloat(entry.price) || 0) * entry.count;
  }
  return [...bucket.values()];
}

function bucketEntriesForMonth(
  entries: SaleEntry[],
  year: number,
  month: number,
): Row[] {
  const bucket = new Map<string, Row>();
  for (const entry of entries) {
    const d = new Date(entry.at);
    if (d.getFullYear() !== year || d.getMonth() !== month) continue;
    const key = `${year}-${String(month + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
    if (!bucket.has(key))
      bucket.set(key, { label: key, count: 0, revenue: 0, byProfile: {} });
    const row = bucket.get(key)!;
    row.count += entry.count;
    row.revenue += (parseFloat(entry.price) || 0) * entry.count;
  }
  return [...bucket.values()];
}

function summarize(source: Row[]) {
  let generated = 0,
    revenue = 0;
  for (const row of source) {
    generated += row.count;
    revenue += row.revenue;
  }
  return { generated, revenue };
}

const summary = computed(() => summarize(rows.value));
const prevSummary = computed(() => summarize(prevRows.value));

// Percentage change vs. the previous period. null when there's no prior data
// to compare against (avoids a misleading "+100%" from a zero baseline).
function percentDelta(current: number, prev: number): number | null {
  if (prev === 0) return current > 0 ? null : 0;
  return ((current - prev) / prev) * 100;
}

const generatedDelta = computed(() =>
  percentDelta(summary.value.generated, prevSummary.value.generated),
);
const revenueDelta = computed(() =>
  percentDelta(summary.value.revenue, prevSummary.value.revenue),
);

const totalByProfile = computed(() => {
  const totals: Record<string, number> = {};
  for (const row of rows.value) {
    for (const [name, cnt] of Object.entries(row.byProfile)) {
      totals[name] = (totals[name] ?? 0) + cnt;
    }
  }
  return totals;
});

// Top profiles ranked by count — shown as quick-glance cards above the table.
const topProfiles = computed(() => {
  return Object.entries(totalByProfile.value)
    .map(([profile, count]) => ({
      profile,
      count,
      revenue: filteredLedger.value
        .filter((e) => e.profile === profile)
        .reduce((s, e) => s + (parseFloat(e.price) || 0) * e.count, 0),
    }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 4);
});

// Profile columns shown in the breakdown table.
const visibleProfiles = computed(() =>
  filterProfile.value ? [filterProfile.value] : allProfileNames.value,
);

// ── Chart points ───────────────────────────────────────────
const chartPoints = computed(() =>
  rows.value.map((r) => ({
    label: r.label,
    count: r.count,
    revenue: r.revenue,
  })),
);

// ── Formatter ─────────────────────────────────────────────
function fmtAmount(n: number): string {
  return formatCompactAmount(n, activeCurrency.value);
}

// ── CSV export ────────────────────────────────────────────
function exportCsv() {
  const headers = [
    period.value === "monthly" ? "Month" : "Date",
    "Generated",
    ...visibleProfiles.value,
    "Revenue",
  ];
  const csvRows = rows.value.map((row) => [
    row.label,
    row.count,
    ...visibleProfiles.value.map((name) => row.byProfile[name] ?? 0),
    row.revenue,
  ]);
  const scope =
    period.value === "monthly"
      ? String(selectedYear.value)
      : `${MONTH_NAMES[selectedMonth.value]}-${selectedYear.value}`;
  downloadCsv(`pikro-sales-${scope}.csv`, headers, csvRows);
}
</script>
