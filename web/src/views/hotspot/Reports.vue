<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <main v-else class="flex flex-col gap-6 px-4 sm:px-6 lg:px-10 pt-7 pb-14">
    <header class="flex flex-wrap items-end justify-between gap-4">
      <div class="flex items-center gap-4">
        <img :src="inkIcon('reports')" alt="" width="64" height="64" class="shrink-0" />
        <div>
          <nav :aria-label="t('reports.breadcrumb')" class="eyebrow">
            {{ t('reports.breadcrumb') }} / {{ t('nav.reports') }}
          </nav>
          <h1 class="type-display m-0 mt-1.5 text-4xl leading-none">{{ t('reports.title') }}</h1>
        </div>
      </div>

      <div class="flex flex-wrap items-center gap-2.5">
        <div
          role="group"
          :aria-label="t('reports.period')"
          class="flex border-2 border-ink rounded-control overflow-hidden bg-paper"
        >
          <button
            v-for="(p, i) in (['daily', 'monthly'] as const)"
            :key="p"
            type="button"
            :aria-pressed="period === p"
            class="min-h-10 px-4 text-sm font-semibold cursor-pointer"
            :class="[i > 0 && 'border-l-2 border-ink', period === p ? 'bg-primary text-primary-fg' : 'text-ink hover:bg-primary-wash']"
            @click="period = p"
          >
            {{ t(`reports.${p}`) }}
          </button>
        </div>
        <label class="flex items-center gap-2 min-h-11 px-3 rounded-control border-2 border-ink bg-paper">
          <span class="eyebrow">{{ t('reports.year') }}</span>
          <select v-model.number="year" class="bg-transparent font-semibold text-sm outline-none cursor-pointer">
            <option v-for="y in YEARS" :key="y" :value="y">{{ y }}</option>
          </select>
        </label>
        <label
          v-if="period === 'daily'"
          class="flex items-center gap-2 min-h-11 px-3 rounded-control border-2 border-ink bg-paper"
        >
          <span class="eyebrow">{{ t('reports.month') }}</span>
          <select v-model.number="month" class="bg-transparent font-semibold text-sm outline-none cursor-pointer">
            <option v-for="(m, i) in monthNames" :key="i" :value="i">{{ m }}</option>
          </select>
        </label>
        <label class="flex items-center gap-2 min-h-11 px-3 rounded-control border-2 border-ink bg-paper">
          <span class="eyebrow">{{ t('reports.plan') }}</span>
          <select v-model="plan" class="bg-transparent font-semibold text-sm outline-none cursor-pointer">
            <option value="">{{ t('reports.allPlans') }}</option>
            <option v-for="p in allPlans" :key="p" :value="p">{{ p }}</option>
          </select>
        </label>
        <button type="button" class="btn" :disabled="!hasAny" @click="exportCsv">
          <ArrowDownTrayIcon class="size-4" />
          {{ t('reports.exportCsv') }}
        </button>
      </div>
    </header>

    <div v-if="loading && !hasAny" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <div
      v-else-if="error"
      role="alert"
      class="flex items-center gap-2.5 p-4 rounded-control border-2 border-danger-text bg-danger-wash text-danger-text text-sm"
    >
      <ExclamationTriangleIcon class="size-5 shrink-0" />
      {{ error }}
      <button type="button" class="btn btn-sm ml-auto" @click="load">{{ t('reports.retry') }}</button>
    </div>

    <div v-else-if="!hasAny" class="card border-dashed shadow-none p-10 text-center">
      <img :src="inkIcon('reports')" alt="" width="64" height="64" class="mx-auto" />
      <p class="m-0 mt-3 font-semibold">{{ t('reports.emptyTitle') }}</p>
      <p class="m-0 mt-1 mx-auto max-w-110 text-sm text-text-muted">{{ t('reports.emptyBody') }}</p>
    </div>

    <template v-else>
      <!-- Summary -->
      <section :aria-label="t('reports.summary')" class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-5">
        <div class="card bg-primary text-primary-fg px-5 py-4.5 flex flex-col gap-2">
          <span class="eyebrow text-primary-fg">{{ t('reports.revenue') }}</span>
          <span class="type-display text-4xl leading-none">
            {{ fmtNumber(current.revenue) }} <span class="text-base">{{ currency }}</span>
          </span>
          <span v-if="revenueDelta !== null" class="chip self-start bg-primary-fg text-primary border-primary-fg">
            {{ t('reports.vsPrev', { delta: fmtDelta(revenueDelta), period: prevLabel }) }}
          </span>
        </div>
        <div class="card px-5 py-4.5 flex flex-col gap-2">
          <span class="eyebrow">{{ t('reports.generated') }}</span>
          <span class="type-display text-4xl leading-none">{{ fmtNumber(current.generated) }}</span>
          <span class="text-[13px] text-text-muted">{{ t('reports.inBatches', current.batchCount) }}</span>
        </div>
        <div class="card px-5 py-4.5 flex flex-col gap-2">
          <span class="eyebrow">{{ t('reports.activated') }}</span>
          <span class="type-display text-4xl leading-none">{{ fmtNumber(current.activated) }}</span>
          <div
            role="meter"
            aria-valuemin="0"
            aria-valuemax="100"
            :aria-valuenow="activatedPct"
            :aria-label="t('reports.activated')"
            class="h-2.5 rounded-full border-2 border-ink overflow-hidden"
          >
            <div class="h-full bg-ink" :style="{ width: `${activatedPct}%` }" />
          </div>
          <span class="text-[13px] text-text-muted">{{ t('reports.ofGenerated', { pct: activatedPct }) }}</span>
        </div>
        <div class="card border-dashed shadow-none px-5 py-4.5 flex flex-col gap-2">
          <span class="eyebrow">{{ t('reports.unused') }}</span>
          <span class="type-display text-4xl leading-none">{{ fmtNumber(current.unused) }}</span>
          <span class="text-[13px] text-text-muted">{{ t('reports.unusedHint') }}</span>
        </div>
      </section>
      <p class="m-0 -mt-2 text-[13px] text-text-muted">{{ t('reports.note') }}</p>

      <div class="grid grid-cols-1 xl:grid-cols-[2fr_1fr] items-stretch gap-6">
        <!-- Stacked bars by plan -->
        <section class="card p-6 flex flex-col gap-4 min-w-0">
          <div class="flex flex-wrap items-center justify-between gap-3">
            <h2 class="type-display m-0 text-2xl">{{ chartTitle }}</h2>
            <div class="flex flex-wrap gap-3.5 text-[13px]">
              <span v-for="s in series" :key="s.key" class="inline-flex items-center gap-1.5">
                <span class="size-3.5 rounded-[3px] border-2 border-ink" :class="s.fill" />{{ s.label }}
              </span>
            </div>
          </div>
          <div
            role="img"
            :aria-label="t('reports.chartLabel')"
            class="flex items-end gap-1.5 h-55 border-b-2 border-ink"
          >
            <div
              v-for="b in buckets"
              :key="b.index"
              class="flex-1 min-w-0 h-full flex flex-col justify-end"
              :title="b.title"
            >
              <div
                v-for="(seg, si) in b.segments"
                :key="seg.key"
                class="border-2 border-b-0 border-ink"
                :class="[seg.fill, si === 0 && 'rounded-t-sm']"
                :style="{ height: `${seg.pct}%` }"
              />
            </div>
          </div>
          <div class="flex justify-between font-mono text-[11px] text-text-muted">
            <span v-for="a in axis" :key="a">{{ a }}</span>
          </div>
        </section>

        <!-- By plan -->
        <section class="card p-6 flex flex-col gap-3.5 min-w-0">
          <h2 class="type-display m-0 text-2xl">{{ t('reports.byPlan') }}</h2>
          <p v-if="byPlan.length === 0" class="m-0 text-sm text-text-muted">{{ t('reports.nothingInPeriod') }}</p>
          <div
            v-for="p in byPlan"
            :key="p.name"
            class="flex flex-col gap-1.5 pb-3 border-b border-hairline last:border-0"
          >
            <div class="flex items-baseline justify-between gap-2.5">
              <span class="chip" :class="p.chip">{{ p.name }}</span>
              <span class="type-display text-lg">{{ fmtNumber(p.revenue) }} {{ currency }}</span>
            </div>
            <div class="flex items-center gap-2.5">
              <div class="flex-1 h-2.5 rounded-full border-2 border-ink overflow-hidden">
                <div class="h-full" :class="p.fill" :style="{ width: `${p.share}%` }" />
              </div>
              <span class="min-w-24 text-right font-mono text-xs">{{ t('reports.soldShare', { n: p.sold, pct: p.share }) }}</span>
            </div>
          </div>
        </section>
      </div>

      <!-- Recent days / months -->
      <section class="card overflow-hidden">
        <h2 class="type-display m-0 px-6 py-4.5 text-2xl">
          {{ period === 'daily' ? t('reports.last7Days') : t('reports.byMonth') }}
        </h2>
        <div class="overflow-x-auto border-t-2 border-ink">
          <table class="w-full min-w-140 text-sm">
            <thead>
              <tr class="border-b-2 border-ink">
                <th class="eyebrow text-left px-4 py-3.5">{{ period === 'daily' ? t('reports.colDay') : t('reports.colMonth') }}</th>
                <th class="eyebrow text-right px-4 py-3.5">{{ t('reports.colGenerated') }}</th>
                <th class="eyebrow text-right px-4 py-3.5">{{ t('reports.colActivated') }}</th>
                <th class="eyebrow text-right px-4 py-3.5">{{ t('reports.colRevenue') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="r in tableRows" :key="r.index" class="border-b border-hairline">
                <td class="px-4 py-3 font-semibold">{{ r.label }}</td>
                <td class="px-4 py-3 text-right font-mono">{{ fmtNumber(r.generated) }}</td>
                <td class="px-4 py-3 text-right font-mono">{{ fmtNumber(r.activated) }}</td>
                <td class="px-4 py-3 text-right font-mono font-semibold">{{ fmtNumber(r.revenue) }} {{ currency }}</td>
              </tr>
              <tr class="bg-primary text-primary-fg">
                <td class="px-4 py-3 font-extrabold border-t-2 border-ink">{{ t('reports.total') }}</td>
                <td class="px-4 py-3 text-right font-mono font-bold border-t-2 border-ink">{{ fmtNumber(tableTotal.generated) }}</td>
                <td class="px-4 py-3 text-right font-mono font-bold border-t-2 border-ink">{{ fmtNumber(tableTotal.activated) }}</td>
                <td class="px-4 py-3 text-right font-mono font-extrabold border-t-2 border-ink">
                  {{ fmtNumber(tableTotal.revenue) }} {{ currency }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </template>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useI18n } from "vue-i18n";
import { ArrowDownTrayIcon, ExclamationTriangleIcon } from "@heroicons/vue/24/outline";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import { useRoutersStore } from "@/stores/routers";
import {
  getSalesLedger,
  getHotspotSettings,
  listHotspotUsers,
  type SalesData,
  type SalesBatch,
  type SalesActivation,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import { downloadCsv } from "@/utils/csv";
import { inkIcon } from "@/utils/inkIcons";

const store = useRoutersStore();
const { t, locale } = useI18n();

const NOW = new Date();
const YEARS = [NOW.getFullYear(), NOW.getFullYear() - 1, NOW.getFullYear() - 2];
// The top two plans get solid fills; the rest share a hatched "other" fill.
const FILLS = ["bg-primary", "bg-ink", "bg-[repeating-linear-gradient(45deg,var(--color-ink)_0_2px,var(--color-paper)_2px_5px)]"];
const CHIPS = ["bg-primary text-primary-fg", "bg-ink text-paper", ""];

const period = ref<"daily" | "monthly">("daily");
const year = ref(NOW.getFullYear());
const month = ref(NOW.getMonth());
const plan = ref("");

const empty = (): SalesData => ({ batches: [], activations: [] });
const ledger = ref<SalesData>(empty());
const prevYearLedger = ref<SalesData>(empty());
const currency = ref("");
const loading = ref(false);
const error = ref("");

async function load() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    // Listing users records any newly used vouchers, so the report is current.
    await listHotspotUsers(store.activeId).catch(() => null);
    const [cur, prev, settings] = await Promise.all([
      getSalesLedger(store.activeId, year.value),
      getSalesLedger(store.activeId, year.value - 1).catch(empty),
      getHotspotSettings(store.activeId).catch(() => null),
    ]);
    ledger.value = cur;
    prevYearLedger.value = prev;
    currency.value =
      settings?.currency ||
      cur.activations.find((a) => a.currency)?.currency ||
      cur.batches.find((b) => b.currency)?.currency ||
      "";
  } catch (e: any) {
    error.value = friendlyError(e, t("reports.failedToLoad"));
  } finally {
    loading.value = false;
  }
}
watch([() => store.activeId, year], load, { immediate: true });

const hasAny = computed(() => ledger.value.batches.length > 0 || ledger.value.activations.length > 0);

const allPlans = computed(() =>
  [...new Set([...ledger.value.batches.map((b) => b.profile), ...ledger.value.activations.map((a) => a.profile)])]
    .filter(Boolean)
    .sort(),
);

const monthNames = computed(() =>
  Array.from({ length: 12 }, (_, i) => new Date(2000, i, 1).toLocaleDateString(locale.value, { month: "long" })),
);

// ── Periods ──────────────────────────────────────────────────────────────
interface Span { year: number; month: number | null } // month null = whole year

const currentSpan = computed<Span>(() => ({ year: year.value, month: period.value === "daily" ? month.value : null }));
const prevSpan = computed<Span>(() =>
  period.value === "monthly"
    ? { year: year.value - 1, month: null }
    : month.value === 0
      ? { year: year.value - 1, month: 11 }
      : { year: year.value, month: month.value - 1 },
);

function inSpan(iso: string, s: Span): boolean {
  const d = new Date(iso);
  return d.getFullYear() === s.year && (s.month === null || d.getMonth() === s.month);
}

const allBatches = computed(() => [...prevYearLedger.value.batches, ...ledger.value.batches]);
const allActivations = computed(() => [...prevYearLedger.value.activations, ...ledger.value.activations]);

function matchesPlan(x: { profile: string }): boolean {
  return !plan.value || x.profile === plan.value;
}

function totals(span: Span) {
  const batches = allBatches.value.filter((b) => matchesPlan(b) && inSpan(b.at, span));
  const activations = allActivations.value.filter((a) => matchesPlan(a) && inSpan(a.at, span));
  const ids = new Set(batches.map((b) => b.id));
  const generated = batches.reduce((s, b) => s + b.count, 0);
  const usedFromThese = allActivations.value.filter((a) => ids.has(a.batch)).length;
  return {
    batches,
    activations,
    generated,
    batchCount: ids.size,
    activated: activations.length,
    revenue: activations.reduce((s, a) => s + (parseFloat(a.price) || 0), 0),
    unused: Math.max(0, generated - usedFromThese),
  };
}

const current = computed(() => totals(currentSpan.value));
const previous = computed(() => totals(prevSpan.value));

const activatedPct = computed(() =>
  current.value.generated ? Math.min(100, Math.round((current.value.activated / current.value.generated) * 100)) : 0,
);

// No delta without a previous-period baseline (avoids a misleading "+100%").
const revenueDelta = computed(() =>
  previous.value.revenue ? ((current.value.revenue - previous.value.revenue) / previous.value.revenue) * 100 : null,
);

const prevLabel = computed(() =>
  prevSpan.value.month === null ? String(prevSpan.value.year) : monthNames.value[prevSpan.value.month],
);

// ── Chart & table buckets ────────────────────────────────────────────────
// Daily: days of the selected month (up to today); monthly: months of the year (up to now).
const bucketCount = computed(() => {
  const isCurrent = year.value === NOW.getFullYear() && (period.value === "monthly" || month.value === NOW.getMonth());
  if (period.value === "monthly") return isCurrent ? NOW.getMonth() + 1 : 12;
  return isCurrent ? NOW.getDate() : new Date(year.value, month.value + 1, 0).getDate();
});

function bucketIndex(iso: string): number {
  const d = new Date(iso);
  return period.value === "monthly" ? d.getMonth() : d.getDate() - 1;
}

function bucketLabel(i: number, long = false): string {
  if (period.value === "monthly") return monthNames.value[i];
  const d = new Date(year.value, month.value, i + 1);
  return d.toLocaleDateString(locale.value, long ? { weekday: "short", day: "numeric" } : { day: "numeric", month: "short" });
}

// Plans ranked by vouchers sold in the period; the top two get their own fill.
const rankedPlans = computed(() => {
  const counts = new Map<string, number>();
  for (const a of current.value.activations) counts.set(a.profile, (counts.get(a.profile) ?? 0) + 1);
  return [...counts.entries()].sort((a, b) => b[1] - a[1]).map(([name]) => name);
});

function rankOf(profile: string): number {
  return Math.min(2, rankedPlans.value.indexOf(profile));
}

const series = computed(() => {
  const ranked = rankedPlans.value;
  const out = ranked.slice(0, 2).map((name, i) => ({ key: name, label: name, fill: FILLS[i] }));
  if (ranked.length > 2) out.push({ key: "__other", label: t("reports.otherPlans"), fill: FILLS[2] });
  return out;
});

const buckets = computed(() => {
  const n = bucketCount.value;
  const perBucket = Array.from({ length: n }, () => [0, 0, 0]);
  const revenue = Array.from({ length: n }, () => 0);
  for (const a of current.value.activations) {
    const i = bucketIndex(a.at);
    if (i >= n) continue; // router clock ahead of the browser's
    perBucket[i][rankOf(a.profile)]++;
    revenue[i] += parseFloat(a.price) || 0;
  }
  const max = Math.max(1, ...perBucket.map((c) => c[0] + c[1] + c[2]));
  return perBucket.map((c, i) => ({
    index: i,
    title: `${bucketLabel(i)}: ${c[0] + c[1] + c[2]} · ${fmtNumber(revenue[i])} ${currency.value}`,
    // Top segment first (rendered on top), largest plan at the bottom.
    segments: [2, 1, 0]
      .filter((r) => c[r] > 0)
      .map((r) => ({ key: r, fill: FILLS[r], pct: (c[r] / max) * 100 })),
  }));
});

const axis = computed(() => {
  const n = bucketCount.value;
  const picks = [...new Set([0, Math.floor((n - 1) / 3), Math.floor((2 * (n - 1)) / 3), n - 1])];
  return picks.map((i) => bucketLabel(i));
});

const chartTitle = computed(() =>
  period.value === "daily"
    ? t("reports.chartDaily", { month: monthNames.value[month.value] })
    : t("reports.chartMonthly", { year: year.value }),
);

const byPlan = computed(() => {
  const total = current.value.revenue || 1;
  const rows = new Map<string, { sold: number; revenue: number }>();
  for (const a of current.value.activations) {
    const r = rows.get(a.profile) ?? { sold: 0, revenue: 0 };
    r.sold++;
    r.revenue += parseFloat(a.price) || 0;
    rows.set(a.profile, r);
  }
  return [...rows.entries()]
    .sort((a, b) => b[1].revenue - a[1].revenue)
    .map(([name, r]) => ({
      name,
      ...r,
      share: Math.round((r.revenue / total) * 100),
      fill: FILLS[rankOf(name)],
      chip: CHIPS[rankOf(name)],
    }));
});

function rowsFor(batches: SalesBatch[], activations: SalesActivation[]) {
  const n = bucketCount.value;
  const rows = Array.from({ length: n }, (_, i) => ({ index: i, label: bucketLabel(i, true), generated: 0, activated: 0, revenue: 0 }));
  for (const b of batches) {
    const i = bucketIndex(b.at);
    if (i < n) rows[i].generated += b.count;
  }
  for (const a of activations) {
    const i = bucketIndex(a.at);
    if (i < n) {
      rows[i].activated++;
      rows[i].revenue += parseFloat(a.price) || 0;
    }
  }
  return rows;
}

const allRows = computed(() => rowsFor(current.value.batches, current.value.activations));

// Daily view shows the last 7 days of the period; monthly shows every month so far.
const tableRows = computed(() => {
  const rows = [...allRows.value].reverse();
  return period.value === "daily" ? rows.slice(0, 7) : rows;
});

const tableTotal = computed(() =>
  tableRows.value.reduce(
    (s, r) => ({ generated: s.generated + r.generated, activated: s.activated + r.activated, revenue: s.revenue + r.revenue }),
    { generated: 0, activated: 0, revenue: 0 },
  ),
);

// ── Formatting & export ─────────────────────────────────────────────────
function fmtNumber(n: number): string {
  return new Intl.NumberFormat(locale.value).format(Math.round(n));
}

function fmtDelta(d: number): string {
  return `${d >= 0 ? "+" : "−"}${Math.abs(Math.round(d))}%`;
}

function exportCsv() {
  const scope = period.value === "daily" ? `${year.value}-${String(month.value + 1).padStart(2, "0")}` : String(year.value);
  downloadCsv(
    `pikro-sales-${scope}${plan.value ? "-" + plan.value : ""}.csv`,
    [
      period.value === "daily" ? t("reports.colDay") : t("reports.colMonth"),
      t("reports.colGenerated"),
      t("reports.colActivated"),
      `${t("reports.colRevenue")}${currency.value ? ` (${currency.value})` : ""}`,
    ],
    allRows.value.map((r) => [bucketLabel(r.index), r.generated, r.activated, Math.round(r.revenue)]),
  );
}
</script>
