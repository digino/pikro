<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <main v-else class="flex flex-col gap-8 px-4 sm:px-6 lg:px-10 pt-7 pb-14">
    <!-- Router header -->
    <header class="flex flex-wrap items-center justify-between gap-4">
      <div class="flex flex-wrap items-center gap-3.5">
        <img :src="inkIcon('routers')" alt="" width="56" height="56" />
        <div>
          <div class="eyebrow">{{ t('dashboard.connectedRouter') }}</div>
          <div class="flex items-center gap-2.5 mt-1.5">
            <span class="type-display text-2xl">{{ store.active()?.name }}</span>
            <span v-if="error" class="chip bg-danger-wash text-danger-text border-danger-text">
              <span class="size-2 rotate-45 bg-danger" />{{ t('dashboard.unreachable') }}
            </span>
            <span v-else-if="isOnline" class="chip bg-online-wash">
              <span class="size-2 rounded-full bg-online" />{{ t('dashboard.online') }}
            </span>
          </div>
          <div class="font-mono text-xs text-text-muted mt-1">{{ routerMeta }}</div>
        </div>
        <RouterLink to="/routers" class="btn btn-sm">
          {{ t('dashboard.switchRouter') }}
          <ChevronDownIcon class="size-4" />
        </RouterLink>
      </div>
      <RouterLink to="/hotspot/vouchers" class="btn btn-primary">
        <img :src="inkIcon('vouchers')" alt="" width="28" height="28" class="-ml-1" />
        {{ t('dashboard.generateVouchers') }}
      </RouterLink>
    </header>

    <!-- Hero -->
    <section class="flex flex-wrap items-end justify-between gap-5">
      <i18n-t
        keypath="dashboard.heroTitle"
        tag="h1"
        class="type-display m-0 text-[clamp(28px,3vw,40px)] leading-tight"
      >
        <template #highlight>
          <span
            class="inline-block px-3 bg-primary text-primary-fg border-2 border-ink rounded-2xl shadow-btn -rotate-2"
            >{{ t('dashboard.heroHighlight') }}</span
          >
        </template>
      </i18n-t>
      <p class="m-0 max-w-85 text-[15px] leading-normal text-text-muted">
        {{ t('dashboard.onlineNow', activeSessions) }}
        <template v-if="cleanupInstalled">{{ t('dashboard.cleanupKeepsTidy') }}</template>
      </p>
    </section>

    <section class="flex flex-wrap gap-6">
      <!-- Sales this month — on amber, contents use primary-fg so they stay dark in both themes -->
      <div class="card flex-[2_1_520px] min-w-0 p-7 flex flex-col gap-5.5 bg-primary text-primary-fg">
        <div class="flex justify-between items-start gap-4">
          <div>
            <div class="eyebrow text-primary-fg">{{ t('dashboard.salesThisMonth') }}</div>
            <div class="type-display mt-3 text-[clamp(44px,6vw,72px)] leading-none">
              {{ fmtNumber(monthSales.revenue) }}
              <span class="text-2xl">{{ salesCurrency }}</span>
            </div>
          </div>
          <img :src="inkIcon('reports')" alt="" width="104" height="104" class="shrink-0 max-sm:hidden" />
        </div>

        <div v-if="salesLoading && !hasSales" class="h-37.5 flex items-center justify-center">
          <span class="spinner spinner--sm" />
        </div>
        <p v-else-if="!hasSales" class="m-0 font-semibold">
          {{ t('dashboard.noSalesRecorded') }}
        </p>
        <div v-else>
          <div
            role="img"
            :aria-label="t('dashboard.dailySales')"
            class="flex items-end gap-1.5 h-37.5 border-b-2 border-primary-fg"
          >
            <div
              v-for="d in salesBars"
              :key="d.day"
              class="flex-1 rounded-t-sm"
              :class="d.today ? 'bg-white border-2 border-b-0 border-primary-fg' : 'bg-primary-fg'"
              :style="{ height: `${d.pct}%` }"
              :title="d.title"
            />
          </div>
          <div class="flex justify-between mt-2 font-mono text-xs">
            <span>{{ firstDayLabel }}</span>
            <span>{{ t('dashboard.today') }}</span>
          </div>
        </div>

        <div class="flex flex-wrap gap-3">
          <RouterLink to="/hotspot/reports" class="btn btn-ink">
            {{ t('dashboard.viewReports') }}
          </RouterLink>
        </div>
      </div>

      <div class="flex-[1_1_300px] min-w-0 flex flex-col gap-5">
        <div class="card px-5 py-4.5 flex items-center gap-4">
          <img :src="inkIcon('users')" alt="" width="64" height="64" class="shrink-0" />
          <div class="flex-1 min-w-0">
            <div class="eyebrow">{{ t('dashboard.activeSessions') }}</div>
            <div class="type-display text-4xl leading-tight">{{ activeSessions }}</div>
            <div class="text-sm text-text-muted">{{ t('dashboard.usersOnRouter', totalUsers) }}</div>
          </div>
          <RouterLink
            to="/hotspot/users?tab=active"
            class="inline-flex items-center min-h-11 text-sm font-semibold text-ink underline underline-offset-2"
          >
            {{ t('dashboard.viewAll') }}
          </RouterLink>
        </div>
        <div class="card px-5 py-4.5 flex items-center gap-4">
          <img :src="inkIcon('vouchers')" alt="" width="64" height="64" class="shrink-0" />
          <div class="flex-1 min-w-0">
            <div class="eyebrow">{{ t('dashboard.vouchersGenerated') }}</div>
            <div class="type-display text-4xl leading-tight">{{ fmtNumber(monthSales.generated) }}</div>
            <div class="text-sm text-text-muted">{{ t('dashboard.thisMonth') }}</div>
          </div>
        </div>
        <SystemHealthCard
          :loading="loading"
          :error="error"
          :resource="resource"
          :has-active-id="!!store.activeId"
        />
      </div>
    </section>

    <section class="flex flex-wrap gap-6">
      <!-- Bandwidth monitor -->
      <div class="card flex-[2_1_520px] min-w-0 p-6 flex flex-col gap-4">
        <div class="flex flex-wrap justify-between items-center gap-3">
          <div class="flex items-center gap-3">
            <img :src="inkIcon('monitor')" alt="" width="48" height="48" />
            <h2 class="type-display m-0 text-2xl">{{ t('dashboard.bandwidth') }}</h2>
            <span v-if="isOnline" class="chip bg-ink text-paper">
              <span class="size-1.75 rounded-full bg-primary" />{{ t('dashboard.live') }}
            </span>
          </div>
          <div class="flex flex-wrap gap-2">
            <span class="chip bg-primary text-primary-fg">
              {{ t('dashboard.download') }} <span class="font-mono">{{ curDown }} Mbps</span>
            </span>
            <span class="chip border-dashed">
              {{ t('dashboard.upload') }} <span class="font-mono">{{ curUp }} Mbps</span>
            </span>
          </div>
        </div>
        <div class="h-55">
          <BandwidthChart :history="bwHistory" />
        </div>
      </div>

      <AutoCleanupCard
        class="flex-[1_1_300px] min-w-0"
        :installed="cleanupInstalled"
        :interval="cleanupInterval"
        :toggling="cleanupToggling"
        :saving="cleanupSaving"
        :error="cleanupError"
        @toggle="toggleCleanup"
        @update:interval="saveCleanupInterval"
      />
    </section>

    <AppDialog
      :open="showDisableCleanupConfirm"
      :title="t('dashboard.turnOffAutoCleanupTitle')"
      @update:open="showDisableCleanupConfirm = $event"
    >
      <div class="space-y-4">
        <p class="text-sm text-text-secondary">
          {{ t('dashboard.turnOffAutoCleanupBody') }}
        </p>
        <div class="flex justify-end gap-2 pt-1">
          <button
            type="button"
            class="btn btn-ghost"
            @click="showDisableCleanupConfirm = false"
          >
            {{ t('dashboard.cancel') }}
          </button>
          <button
            type="button"
            class="btn btn-danger"
            @click="confirmDisableCleanup"
          >
            {{ t('dashboard.turnOff') }}
          </button>
        </div>
      </div>
    </AppDialog>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { RouterLink } from "vue-router";
import { useI18n } from "vue-i18n";
import { ChevronDownIcon } from "@heroicons/vue/24/outline";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import { useRoutersStore } from "@/stores/routers";
import { useToastStore } from "@/stores/toast";
import {
  getSystemResource,
  getSystemClock,
  getPollSnapshot,
  listHotspotUsers,
  listHotspotActive,
  getCleanupScheduler,
  putCleanupScheduler,
  getSalesLedger,
  type SalesData,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import { inkIcon } from "@/utils/inkIcons";
import AppDialog from "@/components/AppDialog.vue";
import BandwidthChart from "@/components/BandwidthChart.vue";
import AutoCleanupCard from "@/components/AutoCleanupCard.vue";
import SystemHealthCard from "@/components/SystemHealthCard.vue";

const store = useRoutersStore();
const toast = useToastStore();
const { t, locale } = useI18n();

// ── State ─────────────────────────────────────────────────────
const resource = ref<Record<string, string>>({});
const loading = ref(false);
const error = ref("");
const routerTime = ref("");
const routerDate = ref("");

const allUsers = ref<Record<string, string>[]>([]);
const activeList = ref<Record<string, string>[]>([]);
const cleanupInstalled = ref<boolean | null>(null);
const cleanupInterval = ref("7d");
const cleanupToggling = ref(false);
const cleanupSaving = ref(false);
const cleanupError = ref("");
const showDisableCleanupConfirm = ref(false);
const hotspotLoading = ref(false);
const hotspotError = ref("");

const salesLedger = ref<SalesData>({ batches: [], activations: [] });
const salesLoading = ref(false);

const N = 46;
interface BwPoint {
  down: number;
  up: number;
}
const bwHistory = ref<BwPoint[]>(
  Array.from({ length: N }, () => ({ down: 0, up: 0 })),
);
const curDown = computed(() => bwHistory.value[N - 1].down.toFixed(2));
const curUp = computed(() => bwHistory.value[N - 1].up.toFixed(2));

// ── Hotspot computed counts ───────────────────────────────────
const totalUsers = computed(() => allUsers.value.length);
const activeSessions = computed(() => activeList.value.length);

// ── Router header ─────────────────────────────────────────────
const isOnline = computed(() => !error.value && Object.keys(resource.value).length > 0);
const routerMeta = computed(() =>
  [
    resource.value["board-name"],
    store.active()?.host,
    resource.value["version"] && `RouterOS ${resource.value["version"].split(" ")[0]}`,
    [routerDate.value, routerTime.value].filter(Boolean).join(" "),
  ]
    .filter(Boolean)
    .join(" · "),
);

function toggleCleanup(enabled: boolean) {
  if (!store.activeId || cleanupToggling.value) return;
  if (!enabled) {
    showDisableCleanupConfirm.value = true;
    return;
  }
  applyCleanupToggle(true);
}

async function applyCleanupToggle(enabled: boolean) {
  if (!store.activeId) return;
  cleanupToggling.value = true;
  cleanupError.value = "";
  try {
    const result = await putCleanupScheduler(
      store.activeId,
      enabled,
      cleanupInterval.value,
    );
    cleanupInstalled.value = result.installed;
    if (result.interval) cleanupInterval.value = result.interval;
    toast.success(
      enabled ? t("dashboard.autoCleanupEnabled") : t("dashboard.autoCleanupDisabled"),
    );
  } catch (e: any) {
    cleanupError.value = friendlyError(e, t("dashboard.failedToUpdateScheduler"));
    toast.error(t("dashboard.failedToUpdateAutoCleanup"), cleanupError.value);
  } finally {
    cleanupToggling.value = false;
  }
}

function confirmDisableCleanup() {
  showDisableCleanupConfirm.value = false;
  applyCleanupToggle(false);
}

async function saveCleanupInterval(interval: string) {
  if (!store.activeId) return;
  cleanupInterval.value = interval;
  cleanupSaving.value = true;
  cleanupError.value = "";
  try {
    const result = await putCleanupScheduler(store.activeId, true, interval);
    cleanupInstalled.value = result.installed;
    if (result.interval) cleanupInterval.value = result.interval;
    toast.success(
      t("dashboard.cleanupScheduleUpdated"),
      t("dashboard.cleanupScheduleUpdatedBody", { interval: result.interval || interval }),
    );
  } catch (e: any) {
    cleanupError.value = friendlyError(e, t("dashboard.failedToUpdateScheduler"));
    toast.error(t("dashboard.failedToUpdateSchedule"), cleanupError.value);
  } finally {
    cleanupSaving.value = false;
  }
}

// ── Data loading ──────────────────────────────────────────────
async function loadStatic() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    const [res, clock] = await Promise.all([
      getSystemResource(store.activeId),
      getSystemClock(store.activeId).catch(
        () => ({}) as Record<string, string>,
      ),
    ]);
    resource.value = res;
    routerTime.value = clock["time"] ?? "";
    routerDate.value = clock["date"] ?? "";
  } catch (e: any) {
    error.value = friendlyError(e, t("dashboard.couldNotReachRouter"));
  } finally {
    loading.value = false;
  }
}

async function loadHotspot() {
  if (!store.activeId) return;
  hotspotLoading.value = true;
  hotspotError.value = "";
  try {
    const [users, active, cleanup] = await Promise.all([
      listHotspotUsers(store.activeId),
      listHotspotActive(store.activeId),
      getCleanupScheduler(store.activeId).catch(() => null),
    ]);
    allUsers.value = users;
    activeList.value = active;
    cleanupInstalled.value = cleanup?.installed ?? null;
    if (cleanup?.interval) cleanupInterval.value = cleanup.interval;
  } catch (e: any) {
    hotspotError.value = friendlyError(e, t("dashboard.couldNotLoadHotspotData"));
  } finally {
    hotspotLoading.value = false;
  }
}

async function loadSales() {
  if (!store.activeId) return;
  salesLoading.value = true;
  try {
    const now = new Date();
    salesLedger.value = await getSalesLedger(store.activeId, now.getFullYear());
  } catch {
    // non-critical
  } finally {
    salesLoading.value = false;
  }
}

async function poll() {
  if (!store.activeId) return;
  try {
    const snap = await getPollSnapshot(store.activeId);
    resource.value = snap.resource;
    error.value = "";
    if (snap.clock["time"]) {
      routerTime.value = snap.clock["time"];
    }
    if (snap.clock["date"]) {
      routerDate.value = snap.clock["date"];
    }
    if (snap.traffic.length > 0) {
      const totalRx = snap.traffic.reduce(
        (s, i) => s + (parseInt(i["rx-bits-per-second"] ?? "0") || 0),
        0,
      );
      const totalTx = snap.traffic.reduce(
        (s, i) => s + (parseInt(i["tx-bits-per-second"] ?? "0") || 0),
        0,
      );
      bwHistory.value = [
        ...bwHistory.value.slice(1),
        { down: totalRx / 1_000_000, up: totalTx / 1_000_000 },
      ];
    }
  } catch {
    // non-critical
  }
}

let pollTimer: ReturnType<typeof setInterval>;
let hotspotTimer: ReturnType<typeof setInterval>;
let salesTimer: ReturnType<typeof setInterval>;

function startTimers() {
  pollTimer = setInterval(() => {
    if (!document.hidden) poll();
  }, 5000);
  hotspotTimer = setInterval(() => {
    if (!document.hidden) loadHotspot();
  }, 15_000);
  salesTimer = setInterval(() => {
    if (!document.hidden) loadSales();
  }, 60_000);
}

function stopTimers() {
  clearInterval(pollTimer);
  clearInterval(hotspotTimer);
  clearInterval(salesTimer);
}

watch(
  () => store.activeId,
  async (id) => {
    stopTimers();
    allUsers.value = [];
    activeList.value = [];
    salesLedger.value = { batches: [], activations: [] };
    bwHistory.value = Array.from({ length: N }, () => ({ down: 0, up: 0 }));
    if (!id) return;
    await Promise.all([loadStatic(), loadHotspot(), loadSales()]);
    await poll();
    startTimers();
  },
  { immediate: true },
);

function onVisibilityChange() {
  if (!document.hidden && store.activeId) poll();
}

onMounted(() => {
  document.addEventListener("visibilitychange", onVisibilityChange);
});

onUnmounted(() => {
  stopTimers();
  document.removeEventListener("visibilitychange", onVisibilityChange);
});

// ── Sales summary (current month) ──────────────────────────────
function inThisMonth(iso: string): boolean {
  const d = new Date(iso);
  const now = new Date();
  return d.getFullYear() === now.getFullYear() && d.getMonth() === now.getMonth();
}

// Generated = batches printed this month; revenue = vouchers actually used this month.
const monthSales = computed(() => {
  let generated = 0;
  let revenue = 0;
  for (const b of salesLedger.value.batches) if (inThisMonth(b.at)) generated += b.count;
  for (const a of salesLedger.value.activations) if (inThisMonth(a.at)) revenue += parseFloat(a.price) || 0;
  return { generated, revenue };
});

const hasSales = computed(
  () => salesLedger.value.batches.length > 0 || salesLedger.value.activations.length > 0,
);

const salesCurrency = computed(
  () =>
    salesLedger.value.activations.find((e) => e.currency)?.currency ??
    salesLedger.value.batches.find((e) => e.currency)?.currency ??
    "",
);

// By-day breakdown for the current month, for the bar chart.
const salesBars = computed(() => {
  const now = new Date();
  const revenueByDay = Array.from({ length: now.getDate() }, () => 0);
  for (const a of salesLedger.value.activations) {
    if (!inThisMonth(a.at)) continue;
    const day = new Date(a.at).getDate();
    // The router clock can run ahead of the browser's; skip days past "today".
    if (day > revenueByDay.length) continue;
    revenueByDay[day - 1] += parseFloat(a.price) || 0;
  }
  const max = Math.max(...revenueByDay, 1);
  return revenueByDay.map((revenue, i) => ({
    day: i + 1,
    today: i === revenueByDay.length - 1,
    pct: (revenue / max) * 100,
    title: `${fmtDay(i + 1)}: ${fmtNumber(revenue)} ${salesCurrency.value}`,
  }));
});

const firstDayLabel = computed(() => fmtDay(1));

function fmtDay(day: number): string {
  const now = new Date();
  return new Date(now.getFullYear(), now.getMonth(), day).toLocaleDateString(
    locale.value,
    { day: "numeric", month: "short" },
  );
}

function fmtNumber(n: number): string {
  return new Intl.NumberFormat(locale.value).format(Math.round(n));
}
</script>
