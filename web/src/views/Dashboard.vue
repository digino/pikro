<template>
  <PageLayout title="Dashboard" subtitle="Router Overview">
    <template v-if="store.activeId" #actions>
      <span class="text-xs font-mono text-text-secondary">
        {{ routerDate || "—" }} {{ routerTime || "" }}
      </span>
    </template>

    <NoRouterSelected v-if="!store.activeId" />

    <div v-else class="grid gap-4">
      <DashboardCard title="Hotspot">
        <HotspotSummaryCard
          :loading="hotspotLoading"
          :total-users="totalUsers"
          :active-sessions="activeSessions"
          :disabled-users="disabledUsers"
        />
        <AutoCleanupCard
          :installed="cleanupInstalled"
          :interval="cleanupInterval"
          :toggling="cleanupToggling"
          :saving="cleanupSaving"
          :error="cleanupError"
          @toggle="toggleCleanup"
          @update:interval="saveCleanupInterval"
        />
      </DashboardCard>

      <div class="grid grid-cols-2 gap-4">
        <!-- Sales report -->
        <DashboardCard title="Sales this month">
          <template #actions>
            <RouterLink to="/hotspot/reports" class="btn btn-ghost">
              View reports
            </RouterLink>
          </template>

          <div
            v-if="salesLoading && salesLedger.length === 0"
            class="flex items-center justify-center py-4"
          >
            <span class="spinner spinner--sm" />
          </div>
          <EmptyState
            v-else-if="salesLedger.length === 0"
            message="No sales recorded yet"
          />
          <template v-else>
            <div class="grid grid-cols-3 gap-3">
              <div class="rounded-lg border border-border p-3 bg-base">
                <div class="text-sm font-medium text-text-secondary">
                  Vouchers generated
                </div>
                <div
                  class="font-mono text-2xl font-bold mt-1 text-text-primary"
                >
                  {{ monthSales.generated }}
                </div>
              </div>
              <div class="rounded-lg border border-border p-3 bg-base">
                <div class="text-sm font-medium text-text-secondary">
                  Revenue
                </div>
                <div class="font-mono text-2xl font-bold mt-1">
                  {{ fmtAmount(monthSales.revenue) }}
                </div>
              </div>
              <div
                class="rounded-lg border border-border p-3 bg-base"
                title="This month's revenue vs. last month."
              >
                <div class="text-sm font-medium text-text-secondary">
                  Performance
                </div>
                <div
                  v-if="salesPerformance === null"
                  class="text-sm font-medium mt-1.5 text-text-muted"
                >
                  No prior month to compare
                </div>
                <div
                  v-else
                  class="font-mono text-2xl font-bold mt-1 flex items-center gap-1"
                  :class="salesPerformance >= 0 ? 'text-green' : 'text-red'"
                >
                  <ArrowTrendingUpIcon
                    v-if="salesPerformance >= 0"
                    class="size-4"
                  />
                  <ArrowTrendingDownIcon v-else class="size-4" />
                  {{ salesPerformance >= 0 ? "+" : ""
                  }}{{ salesPerformance.toFixed(0) }}%
                </div>
              </div>
            </div>
            <div class="h-40">
              <SalesBarChart
                :points="dailySalesPoints"
                :currency="salesCurrency"
              />
            </div>
          </template>
        </DashboardCard>

        <!-- Active sessions -->
        <DashboardCard title="Active sessions">
          <template #actions>
            <RouterLink to="/hotspot/users?tab=active" class="btn btn-ghost">
              View all
            </RouterLink>
          </template>

          <ActiveSessionsTable :sessions="activeList" :loading="hotspotLoading" />
        </DashboardCard>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 items-stretch">
        <!-- System card -->
        <DashboardCard title="System">
          <SystemHealthCard
            :loading="loading"
            :error="error"
            :resource="resource"
            :router-name="store.active()?.name ?? ''"
            :active-sessions="activeSessions"
            :has-active-id="!!store.activeId"
          />
        </DashboardCard>

        <!-- Bandwidth monitor -->
        <DashboardCard title="Bandwidth" class="lg:col-span-2">
          <template #actions>
            <div class="flex items-center gap-1.5">
              <span
                class="size-2 rounded-sm shrink-0"
                style="background: #22d3ee"
              />
              <span class="text-xs text-text-secondary">Download</span>
              <span class="font-mono text-xs font-semibold text-text-primary"
                >{{ curDown }} Mbps</span
              >
            </div>
            <div class="flex items-center gap-1.5">
              <span
                class="size-2 rounded-sm shrink-0"
                style="background: #f59e0b"
              />
              <span class="text-xs text-text-secondary">Upload</span>
              <span class="font-mono text-xs font-semibold text-text-primary"
                >{{ curUp }} Mbps</span
              >
            </div>
          </template>
          <div class="h-full">
            <BandwidthChart :history="bwHistory" />
          </div>
        </DashboardCard>
      </div>
    </div>

    <AppDialog
      :open="showDisableCleanupConfirm"
      title="Turn off auto-cleanup?"
      @update:open="showDisableCleanupConfirm = $event"
    >
      <div class="space-y-4">
        <p class="text-sm text-text-secondary">
          Expired and quota-exhausted vouchers will no longer be removed
          automatically — they'll accumulate until you delete them manually.
        </p>
        <div class="flex justify-end gap-2 pt-1">
          <button
            type="button"
            class="btn btn-ghost"
            @click="showDisableCleanupConfirm = false"
          >
            Cancel
          </button>
          <button
            type="button"
            class="btn btn-danger"
            @click="confirmDisableCleanup"
          >
            Turn off
          </button>
        </div>
      </div>
    </AppDialog>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { RouterLink } from "vue-router";
import {
  ArrowTrendingUpIcon,
  ArrowTrendingDownIcon,
} from "@heroicons/vue/24/outline";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import EmptyState from "@/components/EmptyState.vue";
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
  type SaleEntry,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import { formatCompactAmount } from "@/utils/currencies";
import PageLayout from "@/components/PageLayout.vue";
import AppDialog from "@/components/AppDialog.vue";
import SalesBarChart from "@/components/SalesBarChart.vue";
import BandwidthChart from "@/components/BandwidthChart.vue";
import DashboardCard from "@/components/DashboardCard.vue";
import HotspotSummaryCard from "@/components/HotspotSummaryCard.vue";
import AutoCleanupCard from "@/components/AutoCleanupCard.vue";
import ActiveSessionsTable from "@/components/ActiveSessionsTable.vue";
import SystemHealthCard from "@/components/SystemHealthCard.vue";

const store = useRoutersStore();
const toast = useToastStore();

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

const salesLedger = ref<SaleEntry[]>([]);
const prevYearSalesLedger = ref<SaleEntry[]>([]);
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
const disabledUsers = computed(
  () => allUsers.value.filter((u) => u.disabled === "true").length,
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
      enabled ? "Auto-cleanup enabled" : "Auto-cleanup disabled",
    );
  } catch (e: any) {
    cleanupError.value = friendlyError(e, "Failed to update scheduler");
    toast.error("Failed to update auto-cleanup", cleanupError.value);
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
      "Cleanup schedule updated",
      `Now running every ${result.interval || interval}.`,
    );
  } catch (e: any) {
    cleanupError.value = friendlyError(e, "Failed to update scheduler");
    toast.error("Failed to update schedule", cleanupError.value);
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
    error.value = friendlyError(e, "Could not reach router");
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
    cleanupInterval.value = cleanup?.interval ?? "";
  } catch (e: any) {
    hotspotError.value = friendlyError(e, "Could not load hotspot data");
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
    // Comparing this month to last month dips into the previous year's
    // December only in January — only fetch it then.
    if (now.getMonth() === 0) {
      prevYearSalesLedger.value = await getSalesLedger(
        store.activeId,
        now.getFullYear() - 1,
      ).catch(() => []);
    } else {
      prevYearSalesLedger.value = [];
    }
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
    salesLedger.value = [];
    prevYearSalesLedger.value = [];
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
const monthSalesEntries = computed(() => {
  const now = new Date();
  return salesLedger.value.filter((e) => {
    const d = new Date(e.at);
    return (
      d.getFullYear() === now.getFullYear() && d.getMonth() === now.getMonth()
    );
  });
});

const monthSales = computed(() => {
  let generated = 0;
  let revenue = 0;
  for (const e of monthSalesEntries.value) {
    generated += e.count;
    revenue += (parseFloat(e.price) || 0) * e.count;
  }
  return { generated, revenue };
});

const salesCurrency = computed(
  () => salesLedger.value.find((e) => e.currency)?.currency ?? "",
);

// Previous calendar month's revenue, for a month-over-month comparison.
const prevMonthRevenue = computed(() => {
  const now = new Date();
  const prevMonth = now.getMonth() === 0 ? 11 : now.getMonth() - 1;
  const prevYear =
    now.getMonth() === 0 ? now.getFullYear() - 1 : now.getFullYear();
  let revenue = 0;
  for (const e of [...salesLedger.value, ...prevYearSalesLedger.value]) {
    const d = new Date(e.at);
    if (d.getFullYear() !== prevYear || d.getMonth() !== prevMonth) continue;
    revenue += (parseFloat(e.price) || 0) * e.count;
  }
  return revenue;
});

// Sales performance: this month's revenue vs. last month's, as a % change —
// a trend signal, not an absolute figure easily mistaken for revenue itself.
// null means "no valid baseline" (last month had zero revenue) — the template
// falls back to showing the plain revenue figure in that case.
const salesPerformance = computed(() => {
  if (prevMonthRevenue.value === 0) return null;
  return (
    ((monthSales.value.revenue - prevMonthRevenue.value) /
      prevMonthRevenue.value) *
    100
  );
});

// By-day breakdown for the current month, for the bar chart.
const dailySalesPoints = computed(() => {
  const bucket = new Map<string, { count: number; revenue: number }>();
  for (const e of monthSalesEntries.value) {
    const d = new Date(e.at);
    const key = String(d.getDate());
    const entry = bucket.get(key) ?? { count: 0, revenue: 0 };
    entry.count += e.count;
    entry.revenue += (parseFloat(e.price) || 0) * e.count;
    bucket.set(key, entry);
  }
  return [...bucket.entries()]
    .sort(([a], [b]) => parseInt(a) - parseInt(b))
    .map(([label, v]) => ({ label, count: v.count, revenue: v.revenue }));
});

function fmtAmount(n: number): string {
  return formatCompactAmount(n, salesCurrency.value);
}
</script>
