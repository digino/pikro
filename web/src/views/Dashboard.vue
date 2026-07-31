<template>
  <PageLayout title="Dashboard" subtitle="Router Overview">
    <NoRouterSelected v-if="!store.activeId" />

    <div v-else class="grid gap-4">
      <!-- ── Row 1: hotspot ── -->
      <DashboardCard title="Hotspot">
        <template #actions>
          <RouterLink to="/hotspot/vouchers" class="btn btn-ghost">
            Generate vouchers
          </RouterLink>
          <RouterLink to="/hotspot/users" class="btn btn-primary">
            Manage
          </RouterLink>
        </template>

        <div
          v-if="hotspotLoading && totalUsers === 0"
          class="flex items-center justify-center py-4"
        >
          <span class="spinner spinner--sm" />
        </div>
        <template v-else>
          <div class="grid grid-cols-4 gap-3 tracking-tight text-sm font-semibold text-secondary">
            <div class="rounded-lg border border-border p-3 bg-base">
              <div class="text-text-secondary">Users</div>
              <div class="font-mono text-2xl font-bold mt-1 text-text-primary">{{ totalUsers }}</div>
            </div>
            <div class="rounded-lg border border-border p-3 bg-base">
              <div class="text-text-secondary">Active</div>
              <div class="font-mono text-2xl font-bold mt-1">{{ activeSessions }}</div>
            </div>
            <div class="rounded-lg border border-border p-3 bg-base">
              <div class="text-text-secondary">Expired</div>
              <div class="font-mono text-2xl font-bold mt-1">{{ expiredUsers }}</div>
            </div>
            <div class="rounded-lg border border-border p-3 bg-base">
              <div class="text-text-secondary">Disabled</div>
              <div class="font-mono text-2xl font-bold mt-1 text-text-secondary">{{ disabledUsers }}</div>
            </div>
          </div>

          <div
            v-if="cleanupInstalled !== null"
            class="flex items-center gap-3 p-3 rounded-lg border border-border bg-base"
          >
            <component
              :is="cleanupInstalled ? CheckCircleIcon : ExclamationTriangleIcon"
              class="size-5 shrink-0"
              :class="cleanupInstalled ? 'text-green' : 'text-amber'"
            />
            <div class="flex-1 min-w-0">
              <div class="text-sm font-semibold text-text-primary">
                Auto-cleanup {{ cleanupInstalled ? "active" : "not configured" }}
              </div>
              <div class="text-xs text-text-secondary mt-0.5">
                {{
                  cleanupInstalled
                    ? `Expired vouchers are removed automatically every ${cleanupIntervalLabel}.`
                    : "Not configured — expired users accumulate until removed manually."
                }}
              </div>
            </div>
            <RouterLink
              to="/hotspot/settings"
              class="btn btn-ghost"
              :class="{ 'text-amber': !cleanupInstalled }"
            >
              {{ cleanupInstalled ? "Change" : "Configure" }}
            </RouterLink>
          </div>
        </template>
      </DashboardCard>

      <!-- ── Row 2: sales + active sessions ── -->
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
          <EmptyState v-else-if="salesLedger.length === 0" message="No sales recorded yet" />
          <template v-else>
            <div class="grid grid-cols-3 gap-3">
              <div class="rounded-lg border border-border p-3 bg-base">
                <div class="text-sm font-medium text-text-secondary">Vouchers generated</div>
                <div class="font-mono text-2xl font-bold mt-1 text-text-primary">{{ monthSales.generated }}</div>
              </div>
              <div class="rounded-lg border border-border p-3 bg-base">
                <div class="text-sm font-medium text-text-secondary">Revenue</div>
                <div class="font-mono text-2xl font-bold mt-1">{{ fmtAmount(monthSales.revenue) }}</div>
              </div>
              <div class="rounded-lg border border-border p-3 bg-base" title="Sum of revenue from the last 30 days — a rolling run-rate, not a true subscription MRR.">
                <div class="text-sm font-medium text-text-secondary">Revenue (30d)</div>
                <div class="font-mono text-2xl font-bold mt-1">{{ fmtAmount(revenue30d) }}</div>
              </div>
            </div>
            <div class="h-40">
              <SalesBarChart :points="dailySalesPoints" :currency="salesCurrency" />
            </div>
          </template>
        </DashboardCard>

        <!-- Active sessions -->
        <DashboardCard title="Active sessions">
          <template #actions>
            <RouterLink to="/hotspot/users" class="btn btn-ghost">
              View all
            </RouterLink>
          </template>

          <div
            v-if="hotspotLoading && activeList.length === 0"
            class="flex justify-center py-6"
          >
            <span class="spinner spinner--sm" />
          </div>
          <EmptyState v-else-if="activeList.length === 0" message="No active sessions" />
          <div v-else class="overflow-x-auto">
            <table class="w-full text-xs">
              <thead>
                <tr class="text-text-muted border-b border-border">
                  <th class="text-left pb-2 font-medium">User</th>
                  <th class="text-left pb-2 font-medium">IP</th>
                  <th class="text-right pb-2 font-medium">Uptime</th>
                  <th class="text-right pb-2 font-medium">Down</th>
                  <th class="text-right pb-2 font-medium">Up</th>
                  <th class="text-right pb-2 font-medium">Left</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="s in activeList.slice(0, 8)"
                  :key="s['.id']"
                  class="border-b border-border/40 last:border-0"
                >
                  <td class="py-2 font-mono font-medium text-text-primary">
                    {{ s.user }}
                  </td>
                  <td class="py-2 font-mono text-text-secondary">
                    {{ s.address }}
                  </td>
                  <td class="py-2 text-right font-mono text-text-secondary">
                    {{ s.uptime || "—" }}
                  </td>
                  <td class="py-2 text-right font-mono text-text-primary">
                    {{ formatBytes(parseInt(s["bytes-in"] ?? "0")) }}
                  </td>
                  <td class="py-2 text-right font-mono text-text-secondary">
                    {{ formatBytes(parseInt(s["bytes-out"] ?? "0")) }}
                  </td>
                  <td class="py-2 text-right font-mono text-text-secondary">
                    {{ s["session-time-left"] || "—" }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </DashboardCard>
      </div>

      <!-- ── Row 3: system + bandwidth ── -->
      <div class="grid grid-cols-1 lg:grid-cols-[1fr_2fr] gap-4 items-stretch">
        <!-- System card -->
        <DashboardCard title="System">
          <div v-if="loading" class="flex items-center gap-2 text-sm text-text-muted py-4 justify-center">
            <span class="spinner spinner--sm" /> Loading…
          </div>
          <div v-else-if="error" class="flex items-center gap-1.5 text-xs text-red py-4 justify-center">
            <ExclamationTriangleIcon class="size-3.5 shrink-0" />{{ error }}
          </div>
          <template v-else>
            <div class="flex flex-col items-center gap-2">
              <RouterArt
                :board-name="resource['board-name'] ?? ''"
                :size="88"
                :power-led="healthColor"
                :wifi-led="activeSessions > 0 ? 'var(--color-green)' : 'var(--color-border)'"
                wan-led="var(--color-amber)"
              />
              <div class="text-center space-y-0.5">
                <div class="text-sm font-semibold text-text-primary font-mono leading-tight">
                  {{ resource['board-name'] || store.active()?.name || '—' }}
                </div>
                <div class="text-xs text-text-muted font-mono leading-tight">
                  RouterOS {{ resource['version']?.split(' ')[0] || '—' }}
                </div>
              </div>
              <div class="relative mt-1">
                <svg width="72" height="72" viewBox="0 0 120 120">
                  <circle cx="60" cy="60" r="50" fill="none" stroke="var(--color-border)" stroke-width="8"/>
                  <circle
                    cx="60" cy="60" r="50" fill="none"
                    :stroke="healthColor"
                    stroke-width="8"
                    stroke-linecap="round"
                    :stroke-dasharray="ringCirc"
                    :stroke-dashoffset="ringOffset"
                    transform="rotate(-90 60 60)"
                    style="transition: stroke-dashoffset 0.6s ease, stroke 0.4s ease"
                  />
                </svg>
                <div class="absolute inset-0 flex flex-col items-center justify-center">
                  <span class="font-mono font-bold tracking-tight text-text-primary">{{ healthScore }}</span>
                  <span class="text-xs font-bold text-text-secondary">Health</span>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-1 text-xs">
              <div class="flex justify-between items-center py-1.5 border-b border-border/60">
                <span class="text-text-muted">Uptime</span>
                <span class="font-mono font-semibold text-text-primary">{{ resource["uptime"] ?? "—" }}</span>
              </div>
              <div class="flex justify-between items-center py-1.5 border-b border-border/60">
                <span class="text-text-muted">CPU</span>
                <span class="font-mono font-semibold text-text-primary">{{ resource["cpu-load"] ?? "—" }}%</span>
              </div>
              <div class="flex justify-between items-center py-1.5 border-b border-border/60">
                <span class="text-text-muted">Free RAM</span>
                <span class="font-mono font-semibold text-text-primary">{{ formatBytes(freeMemory) }}</span>
              </div>
              <div class="flex justify-between items-center py-1.5 border-b border-border/60">
                <span class="text-text-muted">Free disk</span>
                <span class="font-mono font-semibold text-text-primary">{{ formatBytes(parseInt(resource["free-hdd-space"] ?? "0") || 0) }}</span>
              </div>
              <div class="flex justify-between items-center py-1.5">
                <span class="text-text-muted">Date &amp; time</span>
                <span class="font-mono font-semibold text-text-primary">{{ routerDate || "—" }} {{ routerTime || "" }}</span>
              </div>
            </div>
          </template>
        </DashboardCard>

        <!-- Bandwidth monitor -->
        <DashboardCard title="Bandwidth">
          <template #actions>
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
          </template>
          <div class="h-full">
            <BandwidthChart :history="bwHistory" />
          </div>
        </DashboardCard>
      </div>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { RouterLink } from "vue-router";
import {
  ExclamationTriangleIcon,
  CheckCircleIcon,
} from "@heroicons/vue/24/outline";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import EmptyState from "@/components/EmptyState.vue";
import { useRoutersStore } from "@/stores/routers";
import {
  getSystemResource,
  getSystemClock,
  getPollSnapshot,
  listHotspotUsers,
  listHotspotActive,
  getCleanupScheduler,
  getSalesLedger,
  type SaleEntry,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import { formatCompactAmount } from "@/utils/currencies";
import PageLayout from "@/components/PageLayout.vue";
import RouterArt from "@/components/router-art/RouterArt.vue";
import SalesBarChart from "@/components/SalesBarChart.vue";
import BandwidthChart from "@/components/BandwidthChart.vue";
import DashboardCard from "@/components/DashboardCard.vue";

const store = useRoutersStore();

// ── State ─────────────────────────────────────────────────────
const resource = ref<Record<string, string>>({});
const loading = ref(false);
const error = ref("");
const routerTime = ref("");
const routerDate = ref("");

const allUsers = ref<Record<string, string>[]>([]);
const activeList = ref<Record<string, string>[]>([]);
const cleanupInstalled = ref<boolean | null>(null);
const cleanupInterval = ref("");
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
const expiredUsers = computed(() => {
  const n = Math.floor(Date.now() / 1000);
  return allUsers.value.filter((u) => {
    const c = u.comment ?? "";
    if (!c.startsWith("exp:")) return false;
    const exp = c.slice(4);
    const epoch = parseInt(exp);
    if (!isNaN(epoch) && epoch > 1_000_000_000) return epoch < n;
    if (exp.includes("-")) return new Date(exp.replace(" ", "T")) < new Date();
    return false;
  }).length;
});
const disabledUsers = computed(
  () => allUsers.value.filter((u) => u.disabled === "true").length,
);

// Cleanup interval is stored as a router-shorthand duration (e.g. "7d", "12h").
// Expand it to a human phrase for the dashboard card.
const cleanupIntervalLabel = computed(() => {
  const s = cleanupInterval.value;
  const m = s.match(/^(\d+)([wdhm])$/i);
  if (!m) return s || "a regular interval";
  const n = parseInt(m[1]);
  const unit = { w: "week", d: "day", h: "hour", m: "minute" }[
    m[2].toLowerCase()
  ]!;
  return `${n} ${unit}${n === 1 ? "" : "s"}`;
});

// ── Health computeds ──────────────────────────────────────────
const cpuLoad = computed(
  () => parseInt(resource.value["cpu-load"] ?? "0") || 0,
);
const freeMemory = computed(
  () => parseInt(resource.value["free-memory"] ?? "0") || 0,
);
const totalMemory = computed(
  () => parseInt(resource.value["total-memory"] ?? "1") || 1,
);

const healthScore = computed(() => {
  if (!store.activeId || error.value) return "—";
  const cpuScore = Math.round((1 - cpuLoad.value / 100) * 40);
  const ramScore = Math.round(
    Math.min(freeMemory.value / totalMemory.value / 0.5, 1) * 30,
  );
  return String(Math.min(100, cpuScore + ramScore + 30));
});
const healthColor = computed(() => {
  if (error.value) return "var(--color-red)";
  const s = parseInt(healthScore.value) || 0;
  if (s >= 75) return "var(--color-green)";
  if (s >= 50) return "var(--color-amber)";
  return "var(--color-red)";
});
const ringCirc = (2 * Math.PI * 50).toFixed(1);
const ringOffset = computed(() => {
  const s = parseInt(healthScore.value) || 0;
  return (2 * Math.PI * 50 * (1 - s / 100)).toFixed(1);
});

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
    // A trailing-30-day window can dip into December of the previous year —
    // only fetch it when we're actually close enough to the boundary to need it.
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

onMounted(() => {
  document.addEventListener("visibilitychange", () => {
    if (!document.hidden && store.activeId) poll();
  });
});

onUnmounted(stopTimers);

// ── Sales summary (current month) ──────────────────────────────
const monthSalesEntries = computed(() => {
  const now = new Date();
  return salesLedger.value.filter((e) => {
    const d = new Date(e.at);
    return d.getFullYear() === now.getFullYear() && d.getMonth() === now.getMonth();
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

// Trailing 30-day revenue run-rate — a rolling total, not a true subscription
// MRR (vouchers are one-off sales), but gives an at-a-glance performance signal.
const revenue30d = computed(() => {
  const cutoff = Date.now() - 30 * 86_400_000;
  let revenue = 0;
  for (const e of [...salesLedger.value, ...prevYearSalesLedger.value]) {
    if (new Date(e.at).getTime() < cutoff) continue;
    revenue += (parseFloat(e.price) || 0) * e.count;
  }
  return revenue;
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

// ── Formatters ────────────────────────────────────────────────
function formatBytes(n: number): string {
  if (!n || isNaN(n)) return "—";
  if (n >= 1_073_741_824) return (n / 1_073_741_824).toFixed(1) + " GB";
  if (n >= 1_048_576) return (n / 1_048_576).toFixed(1) + " MB";
  if (n >= 1_024) return (n / 1_024).toFixed(0) + " KB";
  return n + " B";
}
</script>
